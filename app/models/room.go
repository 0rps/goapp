package models

type Room struct {
	id   int
	name string
}

type UserInRoom struct {
	id     int
	userId int
	roomId int
}
