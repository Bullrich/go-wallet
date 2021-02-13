package wallet

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
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
	return getTokens("erc20-addresses.json")
}

func getTokens(file string) *Tokens {
	af, fe := filepath.Abs(file)
	if fe != nil {
		log.Fatal(fe)
	}
	bs, err := ioutil.ReadFile(af)
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

func (t *Tokens) GetSymbols() []string {
	names := make([]string, 0)
	for _, token := range *t {
		names = append(names, token.Symbol)
	}

	return names
}
