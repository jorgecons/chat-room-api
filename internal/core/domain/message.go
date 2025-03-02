package domain

import (
	"fmt"
	"strings"
	"time"
)

const (
	BotMessage       = "%s quote is $%f per share."
	StockMessageType = "stock"
	BotMessageType   = "bot"
	BotUsername      = "bot"
	StockPrefix      = "/stock="
)

type Message struct {
	Room     string    `json:"room"`
	Username string    `json:"username"`
	Text     string    `json:"text"`
	Date     time.Time `json:"date"`
}

func NewMessage(room, username, text string, date time.Time) Message {
	return Message{
		Room:     room,
		Username: username,
		Text:     text,
		Date:     date,
	}
}

func CreateBotMessage(stockName string, price float64) string {
	return fmt.Sprintf(BotMessage, stockName, price)
}

func GetStockName(text string) string {
	return strings.TrimPrefix(text, StockPrefix)
}
