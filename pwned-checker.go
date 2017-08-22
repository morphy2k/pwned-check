package main

import (
	"bufio"
	"crypto/sha1"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

func hash(p string) [20]byte {
	return sha1.Sum([]byte(p))
}

func request(h [20]byte) (c int) {
	url := fmt.Sprintf("https://haveibeenpwned.com/api/v2/pwnedpassword/%x?originalPasswordIsAHash=true", h)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	c = res.StatusCode

	return
}

func stdin() {

	r := bufio.NewReader(os.Stdin)
	p, _ := r.ReadString('\n')
	p = strings.Replace(p, "\n", "", -1)

	if len(p) > 0 {
		h := hash(p)
		c := request(h)

		if c == 200 {
			fmt.Println("\nPassword have been pwned!")
			os.Exit(0)
		} else if c == 404 {
			os.Exit(0)
		} else {
			os.Exit(1)
		}

	} else {
		os.Exit(1)
	}
}

// interactive mode
func interactive(f bool) {
	var p string

	// password input
	fmt.Println()
	if f {
		fmt.Println("Please enter your password, I will hash it locally")
	}
	fmt.Printf("Password: ")
	fmt.Scanln(&p)

	// hash password with SHA-1
	h := hash(p)
	fmt.Println()
	fmt.Printf("Thanks, the SHA-1 hash of your password is ")
	color.New(color.Bold).Printf("%x\n", h)

	// send hash to haveibeenpwned.com API
	fmt.Println("I will send the hash now to haveibeenpwned.com to check if the password have been pwned ...")
	c := request(h)

	// output result
	fmt.Println()
	if c == 200 {
		color.Red("  Oh no, it have been pwned! :-(\n")
		fmt.Println("  If you use the password somewhere, then change it immediately!")
	} else if c == 404 {
		color.Green("  Congrats, it does not seem to be pwned :-)")
	} else {
		color.Yellow("  Ooops, we have an unkown error here! :-/")
	}
	fmt.Println()

}

func main() {

	var m int

	flag.IntVar(&m, "mode", 0, "0 interactive mode, 1 stdin mode")
	flag.Parse()

	if m == 1 {
		stdin()
	}

	// interactive mode (default)
	interactive(true)

	// repeat?
	for {
		var r string

		fmt.Println("Want to check another password?")
		fmt.Printf("Yes or no (y/n): ")
		fmt.Scanln(&r)

		r = strings.ToLower(r)

		if r == "yes" || r == "y" {
			interactive(false)
		} else if r == "no" || r == "n" {
			break
		} else {
			fmt.Println("Sorry, I don't know what you mean. Please try again")
		}
	}

}
