package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Bullrich/go-wallet/wallet"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	startWebServer()
}

func getPort() int {
	envPort := os.Getenv("PORT")
	if len(envPort) == 0 {
		fmt.Println("No port env detected. Using default")
		return 3000
	}

	port, err := strconv.Atoi(envPort)
	if err != nil {
		log.Fatal(err)
	}

	return port
}

func startWebServer() {
	tokens := wallet.GetTokens()
	engine := html.New("./views", ".html")
	engine.Delims("{{", "}}")

	apiKey := os.Getenv("INFURA_API_KEY")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("balances", nil)
	})

	app.Get("/balances", func(c *fiber.Ctx) error {
		return c.Render("balances", nil)
	})

	app.Get("/balances/:address", func(c *fiber.Ctx) error {
		address := c.Params("address")
		user := wallet.NewUser(apiKey, address)
		if user == nil {
			return c.Render("walletBalance", &fiber.Map{"address": address})
		}

		balances := user.GetAllBalances(*tokens)

		formattedBalance := wallet.LimitDecimals(balances, 4)

		data := &fiber.Map{
			"balances": formattedBalance,
			"address":  address,
		}

		return c.Render("walletBalance", data)
	})

	app.Get("/api/balances/:address", func(c *fiber.Ctx) error {
		user := wallet.NewUser(apiKey, c.Params("address"))
		if user == nil {
			return c.SendStatus(404)
		}

		balance := user.GetAllBalances(*tokens)
		balanceJSON, err := json.Marshal(balance)
		if err != nil {
			return c.SendStatus(500)
		}

		return c.SendString(string(balanceJSON))
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Ok 👋!")
	})

	err := app.Listen(fmt.Sprintf(":%v", getPort()))
	if err != nil {
		log.Fatal(err)
	}
}
