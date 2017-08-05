package main

import (
  "fmt"
  "log"
  "crypto/sha1"
  "net/http"

  "github.com/fatih/color"
)

func hash(p string) (h string)  {
  d := []byte(p)
  h = fmt.Sprintf("%x", sha1.Sum(d))

  return
}

func request(h string) (c int) {
  url := "https://haveibeenpwned.com/api/v2/pwnedpassword/" + h + "?originalPasswordIsAHash=true"
  res, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()

  c = res.StatusCode

  return
}

func main()  {
  var password string

  // password query
  fmt.Println("Please enter your password, I will hash it locally")
  fmt.Printf("Password: ")
  fmt.Scanln(&password)

  // hash password with SHA-1
  hash := hash(password)
  fmt.Printf("Thanks, the SHA-1 hash of your password is ")
  color.New(color.Bold).Println(hash)

  // send hash to haveibeenpwned.com API
  fmt.Println("I will send the hash now to haveibeenpwned.com to check if you have been pwned ...")
  code := request(hash)

  // output result
  fmt.Println()
  if code == 200 {
    color.Red("  Oh no, you have been pwned! :-(\n\n")
    fmt.Println("  If you use the password somewhere, then change it immediately!")
  } else if code == 404 {
    color.Green("  Congrats, you does not seem to be pwned :-)")
  } else {
    color.Yellow("  Ooops, we have an unkown error here! :-/")
  }
  fmt.Println()

}
