package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Bullrich/go-wallet/wallet"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	startWebServer()
}

func startWebServer() {
	tokens := wallet.GetTokens()
	engine := html.New("./views", ".html")
	engine.Delims("{{", "}}")

	apiKey := os.Getenv("INFURA_API_KEY")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/balance/", func(c *fiber.Ctx) error {
		return c.Render("balances", nil)
	})

	app.Get("/balance/:address", func(c *fiber.Ctx) error {
		address := c.Params("address")
		balanceMap := obtainFormattedBalance(apiKey, address)

		return c.Render("walletBalance", balanceMap)
	})

	app.Get("/balances/:address", func(c *fiber.Ctx) error {
		address := c.Params("address")
		user := wallet.NewUser(apiKey, address)
		balances := user.GetAllBalances(*tokens)

		fmt.Println(fmt.Sprintf("%+v", balances))

		data := &fiber.Map{
			"balances": balances,
		}

		return c.Render("walletBalance", data)
	})

	app.Get("/api/balance/:address", func(c *fiber.Ctx) error {
		user := wallet.NewUser(apiKey, c.Params("address"))
		balance := user.GetBalances()
		balanceJSON, err := json.Marshal(balance)
		if err != nil {
			return c.SendStatus(400)
		}

		return c.SendString(string(balanceJSON))
	})

	app.Get("/api/balances/:address", func(c *fiber.Ctx) error {
		user := wallet.NewUser(apiKey, c.Params("address"))
		balance := user.GetAllBalances(*tokens)
		balanceJSON, err := json.Marshal(balance)
		if err != nil {
			return c.SendStatus(400)
		}

		return c.SendString(string(balanceJSON))
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	err := app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

func obtainFormattedBalance(apiKey string, address string) *fiber.Map {
	user := wallet.NewUser(apiKey, address)
	if user == nil {
		return &fiber.Map{
			"validAddress": false,
		}
	}

	balance := user.GetBalances()

	return &fiber.Map{
		"address":      address,
		"wei":          balance["wei"],
		"sai":          balance["sai"],
		"mkr":          balance["mkr"],
		"dai":          balance["dai"],
		"validAddress": true,
	}
}
