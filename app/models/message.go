package models

type Message struct {
	id int
	roomId int
	receiverId int
	sender string
	body string
	date string
}