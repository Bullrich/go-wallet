package wallet

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Tokens is a list with most ERC20 tokens and their addresses
type Tokens []TokenData

type TokenData struct {
	Address string `json:"address"`
	Symbol  string `json:"symbol"`
	Decimal int    `json:"decimal"`
	Type    string `json:"type"`
}

func GetTokens() *Tokens {
	bs, err := ioutil.ReadFile("erc20-addresses.json")
	if err != nil {
		log.Fatal(err)
	}

	tokens := &Tokens{}
	jsonError := json.Unmarshal(bs, tokens)
	if jsonError != nil {
		log.Fatal(jsonError)
	}

	return tokens
}
