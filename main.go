package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	apiURL    = "https://api.pwnedpasswords.com/range"
	userAgent = "pwned-check/1.0 (+https://github.com/morphy2k/pwned-check)"
)

var (
	client *http.Client
)

func init() {
	os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1")
	client = &http.Client{}
}

func toHash(s string) string {
	h := sha1.New()
	io.WriteString(h, s)

	return hex.EncodeToString(h.Sum(nil))
}

func getHashes(hash string) (string, error) {
	url := fmt.Sprintf("%s/%s", apiURL, hash[:5])

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", userAgent)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return "", nil
		}
		return "", fmt.Errorf("unexpected response code: %v", res.StatusCode)
	}

	if ct := res.Header.Get("Content-Type"); ct != "text/plain" {
		return "", fmt.Errorf("wrong content type: %v", ct)
	}

	var b bytes.Buffer

	if _, err := b.ReadFrom(res.Body); err != nil {
		return "", err
	}

	return b.String(), nil
}

func compareHashes(hash, hashes string) (int64, error) {
	var match string

	hashes = strings.ToLower(hashes)

	for _, k := range strings.Split(hashes, "\n") {
		if k[:35] == hash[5:] {
			match = k[36:]
			break
		}
	}

	if match == "" {
		return 0, nil
	}

	return strconv.ParseInt(match[:len(match)-1], 10, 64)
}

func main() {
	var input string
	var inputIsHash bool

	flag.BoolVar(&inputIsHash, "hash", false, "SHA1 hash as input")

	flag.StringVar(&input, "p", "", "Password to check")
	flag.Parse()

	if input == "" {
		r := bufio.NewReader(os.Stdin)

		p, err := r.ReadString('\n')
		if err != nil {
			fmt.Println("I/O error:", err)
			os.Exit(2)
		}

		fmt.Println()

		input = strings.Replace(p, "\n", "", -1)
	}

	if inputIsHash && len(input) != 40 {
		fmt.Println("Invalid SHA1 hash")
		os.Exit(1)
	}

	if !inputIsHash {
		input = toHash(input)
	} else {
		input = strings.ToLower(input)
	}

	hashes, err := getHashes(input)
	if err != nil {
		fmt.Println("API request error:", err.Error())
		os.Exit(1)
	}

	if hashes != "" {
		count, err := compareHashes(input, hashes)
		if err != nil {
			fmt.Println("Invalid API data")
			os.Exit(1)
		}

		if count > 0 {
			fmt.Printf("Oh no - pwned! :-(\nPassword found %v times in the data set\n", count)
		}
	}
}
