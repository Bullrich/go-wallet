# Ethereum Wallet Balance

![Testing](https://github.com/Bullrich/go-wallet/workflows/Continuous%20Testing/badge.svg)


A small app built in Go using Fiber.

Shows the balance of eth and the most populars ERC20 tokens in your wallet.

## How to

You need to get an api key from [infura.io](https://infura.io/)

### With Go installed in your system

Clone the repo and execute the following:
```shell
go get .
INFURA_API_KEY="<infura-api-key>" go run *.go
```

### Testing

To run the test execute the following command:
```shell
INFURA_API_KEY="<infura-api-key>" TEST_ADDRESS="<address>" go test ./...
```

The address must have wei and one of the tokens that is evaluated in the system as it verifies that it can obtain that information.

## Hosting

Check the project online here: [https://go-wallet.herokuapp.com/](https://go-wallet.herokuapp.com/)