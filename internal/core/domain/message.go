package domain

import (
	"fmt"
	"strings"
	"time"
)

const (
	BotMessage       = "%s quote is $%f per share"
	BotErrorMessage  = "%s: %s"
	StockMessageType = "stock"
	BotMessageType   = "bot"
	BotUsername      = "bot"
	StockPrefix      = "/stock="
)

var Now = time.Now // test purposes

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

func NewBotMessage(room string) Message {
	return Message{
		Room:     room,
		Username: BotUsername,
		Date:     Now().UTC(),
	}
}

func CreateBotText(stockName string, price float64) string {
	return fmt.Sprintf(BotMessage, stockName, price)
}

func CreateBotErrorText(stockName string, err error) string {
	return fmt.Sprintf(BotErrorMessage, stockName, err.Error())
}

func GetStockName(text string) string {
	return strings.TrimPrefix(text, StockPrefix)
}
