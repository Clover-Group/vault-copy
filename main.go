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
	var tokenFile = flag.String("t", "./token", "Path to file with token")
	var regExp = flag.String("r", "", "Regular expression to replace old variables")
	var passwordLength = flag.Int("p", 15, "Password length")
	var input = flag.String("i", "", "Path to copy")
	var output = flag.String("o", "", "Path where to copy")
	flag.Parse()
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
}
