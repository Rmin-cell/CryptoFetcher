package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var lastPrices map[string]float64

func main() {
	app := &cli.App{
		Name:  "Crypto Tracker",
		Usage: "Track crypto prices in the terminal",
		Action: func(c *cli.Context) error {
			return run(c)
		},
	}

	lastPrices = make(map[string]float64)

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	// Fetch top 10 most important cryptocurrencies
	top10, err := getTop10()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "COIN\tPRICE\n")

	// Get initial prices and print
	prices := getPrices(top10)
	printPrices(w, prices)
	w.Flush()

	// Update prices every 5 minutes
	for range time.Tick(5 * time.Minute) {
		clearTerminal()
		clearLastPrices()
		prices = getPrices(top10)
		printPrices(w, prices)
		w.Flush()
	}

	return nil
}

// Fetch top 10 most important cryptocurrencies
func getTop10() ([]string, error) {
	url := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=10&page=1"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var coins []struct {
		ID string `json:"id"`
	}

	json.NewDecoder(resp.Body).Decode(&coins)

	var top10 []string
	for _, coin := range coins {
		top10 = append(top10, coin.ID)
	}

	return top10, nil
}

// Get latest prices
func getPrices(coins []string) map[string]float64 {
	url := "https://api.coingecko.com/api/v3/simple/price?ids=" +
		strings.Join(coins, ",") + "&vs_currencies=usd"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	var data map[string]map[string]float64
	json.NewDecoder(resp.Body).Decode(&data)

	prices := make(map[string]float64)
	for _, coin := range coins {
		prices[coin] = data[coin]["usd"]
	}

	return prices
}

// Print formatted prices table
func printPrices(w *tabwriter.Writer, prices map[string]float64) {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	for coin, p := range prices {
		arrow := red("↓")
		if p >= lastPrices[coin] {
			arrow = green("↑")
		}
		fmt.Fprintf(w, "%s\t$%0.2f %s\n", coin, p, arrow)
		lastPrices[coin] = p // Update last prices
	}
}

// Clear terminal
func clearTerminal() {
	cmd := exec.Command("clear") // For Unix-like systems
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// Clear last prices
func clearLastPrices() {
	for coin := range lastPrices {
		lastPrices[coin] = 0
	}
}
