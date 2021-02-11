package wallet

import (
	"log"
	"os"
	"testing"
)

func getClient() *User {
	apiKey := os.Getenv("INFURA_API_KEY")
	testAddress := os.Getenv("TEST_ADDRESS")

	return NewUser(apiKey, testAddress)
}

func TestUser(t *testing.T) {
	client := getClient()
	if client == nil {
		t.Fatal("Fatal is nil")
	}
}

func TestGetWeiBalance(t *testing.T) {
	client := getClient()

	balance, _ := client.GetWeiBalance().Balance.Float64()
	if balance == 0 {
		t.Fatal("No wei for address")
	}
}

func TestAddsWeiBalanceInAllBalances(t *testing.T) {
	client := getClient()
	tokens := getTokenInTestDirectory()

	eb := client.GetWeiBalance()

	balances := client.GetAllBalances(*tokens)

	for _, b := range balances {
		if b.Coin == eb.Coin {
			if b.Balance.Cmp(eb.Balance) != 0 {
				t.Fatal("ETH balance is not equal individually than in a group")
			}
			break
		}
	}
}

func TestReturnValueOfOtherTokens(t *testing.T) {
	client := getClient()
	tokens := getTokenInTestDirectory()

	eb := client.GetWeiBalance()

	balances := client.GetAllBalances(*tokens)

	tt := .0

	for _, b := range balances {
		// if it is not wei
		if b.Coin != eb.Coin {
			cb, _ := b.Balance.Float64()
			tt += cb
		}
	}

	if tt == 0 {
		log.Fatal("Account couldn't retrieve any tokens.")
	}
}
