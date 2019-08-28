# pwned-check
Simple tool that checks if your password have been pwned via the [haveibeenpwned.com](https://haveibeenpwned.com) API.<br>
It hashes your password locally and sends the first 5 characters ([k-anonymity](https://en.wikipedia.org/wiki/K-anonymity)) to the server. Afterwards, the list obtained from the server is locally searched for matches.

## Usage

```
Usage of pwned-check:
  -hash
    	SHA1 hash as input
  -p string
    	Password to check
```

### Examples
```BASH
pwned-check -p foobar
```
```BASH
pwned-check < password.txt
```
```BASH
echo -n "foobar" | sha1sum | cut -d " " -f1 | pwned-check -hash
```
