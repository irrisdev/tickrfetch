package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/irrisdev/go-coinlore/coinlore"
	"github.com/irrisdev/tickrfetch/logger"
	"github.com/irrisdev/tickrfetch/services"
	"github.com/joho/godotenv"
)

func main() {
	client := coinlore.NewClient("https://api.coinlore.net")
	s := services.NewCoinService(client, nil)

	state := make(chan services.ServiceState)
	go s.Start(state, false)

	if msg := <-state; msg != services.Running {
		logger.Error("Service failed to start")
		return
	}

	logger.Info("Service is running")

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		logger.Error("BOT_TOKEN environment variable is not set")
		return
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logger.Error("Failed to create new bot: ", err)
		return
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text != "" && (update.Message.Text[0] == '$' || update.Message.Text[0] == '/') {
			symbol := update.Message.Text[1:]
			coin := s.GetCoin(symbol)
			if coin == nil {
				logger.Error("Failed to get coin: ", symbol)
				continue
			}

			price, err := formatPrice(coin.PriceUSD)
			if err != nil {
				logger.Error("Failed to format price: ", err)
				continue
			}

			mcap, err := formatPrice(coin.MarketCapUSD)
			if err != nil {
				logger.Error("Failed to format price: ", err)
				continue
			}
            change1h := formatChange(coin.PercentChange1h)
            change24h := formatChange(coin.PercentChange24h)
            change7d := formatChange(coin.PercentChange7d)
			response := formatRespMessage(
				coin.Rank, coin.Name, coin.Symbol, price, mcap, change1h, change24h, change7d,
			)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			if _, err := bot.Send(msg); err != nil {
				logger.Error("Failed to send message: ", err)
			}
		}
	}
}

func formatChange(changeStr string) string {
    change, err := strconv.ParseFloat(changeStr, 64)
    if err != nil {
        return changeStr
    }

    if change > 0 {
        return fmt.Sprintf("+%.2f", change)
    }
    return fmt.Sprintf("%.2f", change)
}

func formatRespMessage(id int, name, symbol, price, marketCap, change1h, change24h, change7d string) string {
	const respMessage = "`#%d %s (%s)\n\nPrice        %-10s\nM-Cap        %-10s\n\n1h Change    %-5s%%\n24h Change   %-5s%%\n7d Change    %-5s%%`"

	return escapeMarkdown(fmt.Sprintf(respMessage, id, name, symbol, price, marketCap, change1h, change24h, change7d))
}

func formatPrice(priceStr string) (string, error) {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return "", err
	}

	switch {
	case price >= 1000:
		return fmt.Sprintf("$%s", humanize.Comma(int64(math.Round(price)))), nil
	case price >= 1:
		return fmt.Sprintf("$%.2f", price), nil
	default:
		return fmt.Sprintf("$%.4f", price), nil
	}
}

func escapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}
}
