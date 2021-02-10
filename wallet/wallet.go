package wallet

import (
	"context"
	"fmt"
	"math"
	"math/big"

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

const daiContract = "0x6B175474E89094C44Da98b954EedeAC495271d0F"
const mkrContract = "0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2"
const saiContract = "0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359"

// NewUser constructs a User object to interact with the ethereum network
func NewUser(infuraAPIKey string, address string) *User {
	if !(utils.IsAddressValid(address)) {
		return nil
	}

	infuraAddress := fmt.Sprintf("https://mainnet.infura.io/v3/%s", infuraAPIKey)
	//infuraAddress := fmt.Sprintf("https://rinkeby.infura.io/v3/%s", infuraAPIKey)

	client, err := ethclient.Dial(infuraAddress)
	utils.CheckError(err)

	account := common.HexToAddress(address)

	return &User{client: client, address: account}
}

// GetWeiBalance returns balance of Wei (ether) in the account
func (u User) GetWeiBalance() *big.Float {
	balance, err := u.client.BalanceAt(context.Background(), u.address, nil)
	utils.CheckError(err)

	return divideTokens(balance)
}

// GetDaiBalance returns balance of dai stable coin in the account
func (u User) GetDaiBalance() Coin {
	daiAddress := common.HexToAddress(daiContract)
	return u.getTokenBalance(daiAddress)
}

// GetMkrBalance returns the balance of mkr in the account
func (u User) GetMkrBalance() Coin {
	mkrAddress := common.HexToAddress(mkrContract)
	return u.getTokenBalance(mkrAddress)
}

// GetSaiBalance returns the balance of sai in the account
func (u User) GetSaiBalance() Coin {
	sntAddress := common.HexToAddress(saiContract)
	return u.getTokenBalance(sntAddress)
}

func (u User) getTokenBalance(tokenAddress common.Address) Coin {
	instance, err := token.NewToken(tokenAddress, u.client)
	utils.CheckError(err)

	bal, err := instance.BalanceOf(&bind.CallOpts{}, u.address)
	utils.CheckError(err)

	dividedBalance := divideTokens(bal)

	return dividedBalance
}

// divideTokens get the total amount of tokens and divided it by 10 ^ 18
func divideTokens(tokenAmount *big.Int) Coin {
	tokenFloat := new(big.Float).SetInt(tokenAmount)
	decimals := big.NewFloat(math.Pow10(18))

	return new(big.Float).Quo(tokenFloat, decimals)
}

type valueFunc func(user User) Coin

type CoinValue struct {
	Coin    string
	Balance *big.Float
}

func (u User) GetAllBalances(t Tokens) []CoinValue {
	c := make(chan CoinValue)

	for _, token := range t {
		go u.getTokenDataBalance(token, &c)
	}

	balances := make([]CoinValue, 0)

	defaultValue := big.NewFloat(0)

	for range t {
		balance := <-c
		if balance.Balance.Cmp(defaultValue) != 0 {
			balances = append(balances, balance)
			// balances[balance.coin] = balance.balance
		}
	}

	return balances
}

func (u User) getTokenDataBalance(t TokenData, c *chan CoinValue) {
	tokenAddress := common.HexToAddress(t.Address)
	balance := u.getTokenBalance(tokenAddress)
	(*c) <- CoinValue{Coin: t.Symbol, Balance: balance}
}

// GetBalances returns a map with the balances of sai, mkr, wei and dai in the account
func (u User) GetBalances() map[string]Coin {
	tokens := map[string]valueFunc{
		"sai": func(user User) Coin { return user.GetSaiBalance() },
		"mkr": func(user User) Coin { return user.GetMkrBalance() },
		"wei": func(user User) Coin { return user.GetWeiBalance() },
		"dai": func(user User) Coin { return user.GetDaiBalance() },
	}

	c := make(chan CoinValue)

	for coin, invocation := range tokens {
		go u.fetchBalance(coin, invocation, c)
	}

	balances := make(map[string]Coin)

	for range tokens {
		balance := <-c
		balances[balance.Coin] = balance.Balance
	}

	return balances
}

func (u User) fetchBalance(coin string, invocation valueFunc, c chan CoinValue) {
	balance := invocation(u)
	c <- CoinValue{Coin: coin, Balance: balance}
}
