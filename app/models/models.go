package models

type User struct {
	id       int
	login    string
	password string
}

type Session struct {
	id     int
	secret string
	userId int
}

// func createUser(login, pass string) User {
// 	user := User{login: login, password: pass, id: 1}
// 	return user
// }

func FindUser(id int) (User, bool) {
	var user User
	noerror := false
	if id == 1 {
		user = User{login: "god", password: "1234", id: 1}
		noerror = true
	}
	return user, noerror
}

func Authorize(login, pass string) (Session, bool) {
	if login == "god" && pass == "1234" {
		session := Session{id: 1, secret: "12345678", userId: 1}
		return session, true
	}

	return Session{}, false
}

func GetSession(id int, secret string) (Session, bool) {
	var session Session
	if id == 1 && secret == "12345678" {
		session = Session{id: 1, secret: "12345678", userId: 1}
		return session, true
	}

	return session, false
}

func (s *Session) GetUser() (User, bool) {
	return FindUser(s.userId)
}

func (s *Session) Remove() {

}
