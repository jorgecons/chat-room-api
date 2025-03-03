package message

import "time"

type message struct {
	ID       uint
	Room     string
	Username string
	Text     string
	Date     time.Time
}
