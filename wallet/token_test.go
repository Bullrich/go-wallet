package wallet

import (
	"log"
	"testing"
)

func getTokenInTestDirectory() *Tokens{
	return getTokens("../erc20-addresses.json")
}

func TestToken(t *testing.T) {
	tokens := getTokenInTestDirectory()

	if len(*tokens) == 0 {
		t.Fatal("No tokens found")
	}
}

func TestNoDuplicatedAddress(t *testing.T) {
	tokens := getTokenInTestDirectory()

	coins := make(map[string]bool)

	for _, coin := range *tokens {
		coins[coin.Address] = true
	}

	if len(*tokens) != len(coins) {
		log.Fatal("Length are not equal")
	}
}

func TestNoDuplicatedSymbols(t *testing.T) {
	tokens := getTokenInTestDirectory()

	coins := make(map[string]bool)

	for _, coin := range *tokens {
		coins[coin.Symbol] = true
	}

	if len(*tokens) != len(coins) {
		log.Fatal("Length are not equal")
	}
}
