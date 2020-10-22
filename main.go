package main

import (
	"flag"
	"github.com/hashicorp/vault/api"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}
var vaultAddr = os.Getenv("VAULT_ADDR")

func main() {
	flag.Usage = func() {
		println("Utility to copy with replace whole branches in vault")
		println("Usage: vault-copy [Options]")
		println("Options:")
		flag.PrintDefaults()
	}
	var tokenFile = flag.String("t", "./token", "Path to file with token [mandatory]")
	var regExp = flag.String("r", "", "Sed regular expression to replace old variables (see https://github.com/rwtodd/Go.Sed) [optional]")
	var passwordLength = flag.Int("p", 15, "Password length [optional]")
	var input = flag.String("i", "", "Path to copy [mandatory")
	var output = flag.String("o", "", "Path where to copy [mandatory]")
	flag.Parse()
	if *input == "" {
		panic("Input branch can't be empty!")
	}
	if *output == "" {
		panic("Output branch can't be empty!")
	}
	btoken, err := ioutil.ReadFile(*tokenFile)
	if err != nil {
		panic(err)
	}
	token := string(btoken)
	client, err := api.NewClient(&api.Config{Address: vaultAddr, HttpClient: httpClient})
	if err != nil {
		panic(err)
	}
	token = strings.Replace(token, "\n", "", -1)
	client.SetToken(token)
	vaultCopy(client, *input, *output, *regExp, *passwordLength)
	println("Ok! Branch " + *input + " has been successfully copied to " + *output)
	println("Ok! All fields 'password' and 'secretKey' has been changed to random strings")
	if *regExp != "" {
		println("Ok! All values has been processed by regexp '" + *regExp + "'")
	}
}
