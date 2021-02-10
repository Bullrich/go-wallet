package wallet

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sort"
	"strconv"
	"strings"

	token "github.com/Bullrich/go-wallet/token"
	"github.com/Bullrich/go-wallet/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Coin is the float value of a token (total value divided by 10 ^ 18)
type Coin *big.Float

// User is a container of eht client and the address of the wallet
type User struct {
	client  *ethclient.Client
	address common.Address
}

// NewUser constructs a User object to interact with the ethereum network
func NewUser(infuraAPIKey string, address string) *User {
	if !(utils.IsAddressValid(address)) {
		return nil
	}

	infuraAddress := fmt.Sprintf("https://mainnet.infura.io/v3/%s", infuraAPIKey)

	client, err := ethclient.Dial(infuraAddress)
	utils.CheckError(err)

	account := common.HexToAddress(address)

	return &User{client: client, address: account}
}

// GetWeiBalance returns balance of Wei (ether) in the account
func (u User) GetWeiBalance() *big.Float {
	balance, err := u.client.BalanceAt(context.Background(), u.address, nil)
	utils.CheckError(err)

	return divideTokens(balance, 18)
}

func (u User) getTokenBalance(t *TokenData) Coin {
	tokenAddress := common.HexToAddress(t.Address)
	instance, err := token.NewToken(tokenAddress, u.client)
	utils.CheckError(err)

	bal, err := instance.BalanceOf(&bind.CallOpts{}, u.address)
	utils.CheckError(err)

	dividedBalance := divideTokens(bal, t.Decimal)

	return dividedBalance
}

// divideTokens get the total amount of tokens and divided it by 10 ^ 18
func divideTokens(tokenAmount *big.Int, decimals int) Coin {
	tokenFloat := new(big.Float).SetInt(tokenAmount)
	d := big.NewFloat(math.Pow10(decimals))

	return new(big.Float).Quo(tokenFloat, d)
}

type CoinValue struct {
	Coin    string
	Balance *big.Float
}

type CoinValueString struct {
	Coin    string
	Balance string
}

func (u User) GetAllBalances(t Tokens) []CoinValue {
	c := make(chan CoinValue)

	for _, td := range t {
		go u.getTokenDataBalance(td, &c)
	}

	balances := make([]CoinValue, 0)

	defaultValue := big.NewFloat(0)

	eth := u.GetWeiBalance()

	balances = append(balances, CoinValue{Coin: "ETH", Balance: eth})

	for range t {
		balance := <-c
		if balance.Balance.Cmp(defaultValue) != 0 {
			balances = append(balances, balance)
		}
	}

	sort.Slice(balances, func(i, j int) bool {
		return balances[i].Balance.Cmp(balances[j].Balance) > 0
	})

	return balances
}

func (u User) getTokenDataBalance(t TokenData, c *chan CoinValue) {
	balance := u.getTokenBalance(&t)
	*c <- CoinValue{Coin: t.Symbol, Balance: balance}
}

func LimitDecimals(cv []CoinValue, decimals int) []CoinValueString {
	cvs := make([]CoinValueString, 0)
	format := "%." + strconv.Itoa(decimals) + "f"

	for _, coin := range cv {
		balance, _ := coin.Balance.Float64()
		stringBalance := strings.TrimRight(strings.TrimRight(fmt.Sprintf(format, balance), "0"), ".")
		if stringBalance != "0" {
			cvs = append(cvs, CoinValueString{Coin: coin.Coin, Balance: stringBalance})
		}
	}

	return cvs
}
