package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net/http"
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

func start(f bool) {
	var password string

	// password input
	fmt.Println()
	if f {
		fmt.Println("Please enter your password, I will hash it locally")
	}
	fmt.Printf("Password: ")
	fmt.Scanln(&password)

	// hash password with SHA-1
	h := hash(password)
	fmt.Println()
	fmt.Printf("Thanks, the SHA-1 hash of your password is ")
	color.New(color.Bold).Printf("%x\n", h)

	// send hash to haveibeenpwned.com API
	fmt.Println("I will send the hash now to haveibeenpwned.com to check if the password have been pwned ...")
	code := request(h)

	// output result
	fmt.Println()
	if code == 200 {
		color.Red("  Oh no, it have been pwned! :-(\n")
		fmt.Println("  If you use the password somewhere, then change it immediately!")
	} else if code == 404 {
		color.Green("  Congrats, it does not seem to be pwned :-)")
	} else {
		color.Yellow("  Ooops, we have an unkown error here! :-/")
	}
	fmt.Println()

}

func main() {

	// first start
	start(true)

	// repeat?
	for {
		var r string

		fmt.Println("Want to check another password?")
		fmt.Printf("Yes or no (y/n): ")
		fmt.Scanln(&r)

		r = strings.ToLower(r)

		if r == "yes" || r == "y" {
			start(false)
		} else if r == "no" || r == "n" {
			break
		} else {
			fmt.Println("Sorry, I don't know what you mean. Please try again")
		}
	}

}
