package models

import (
	"time"
)

type Message struct {
	Id     int
	RoomId int
	Sender string
	Body   string
	Date   time.Time
}

func GetArchiveMessages(roomId int) []Message {
	var result []Message

	result = make([]Message, 1)
	msg := Message{Id: 0, RoomId: 1,
		Sender: "Chaat boot", Body: "hello",
		Date: time.Now()}

	result[0] = msg

	return result
}
