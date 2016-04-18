package models

import (
	"fmt"
	"github.com/revel/revel"
	"goapp/app"
)

type User struct {
	Id       int
	Login    string
	Password string
}

type Session struct {
	Id     int
	Secret string
	UserId int
}

func CreateUser(login, pass string) User {
	var user User

	queryRes, err := app.DB.Exec(`
		INSERT INTO Users(login, password)
		VALUES(?,?)`, login, pass)

	if err != nil {
		revel.INFO.Print("Error in createuser: ", err)
	} else {
		var id int64
		id, err = queryRes.LastInsertId()
		user = User{Login: login, Password: pass, Id: int(id)}
	}
	return user
}

func FindUser(id int) (User, bool) {
	var user User
	noerror := false
	if id == 1 {
		user = User{Login: "god", Password: "1234", Id: 1}
		noerror = true
	}
	return user, noerror
}

func FindUser(login string) bool {
	// queryStr := fmt.Sprintf("SELECT COUNT(*) AS cnt FROM Users WHERE login like '%%s%'", login)
	// /// TODO: fix
	// queryRes, err := app.DB.Query(queryStr)

	return false
}

func FindUser(login, password string) (User, bool) {
	return User{}, false
}

func Authorize(login, pass string) (Session, bool) {
	var session Session
	var err bool
	// user, isExists :=FindUser(login, pass)
	// if isExists {
	// 	session, err = createSession(user.Id)
	// } else {
	// 	return session, false
	// }
	return session, false
}

func GetSession(id int, secret string) (Session, bool) {
	var session Session

	return session, false
}

func generateSecret() string {
	return "12345667"
}

func createSession(userId int) (Session, bool) {
	var session Session
	// secret := generateSecret()
	// queryRes, err := app.DB.Exec(`
	// 	INSERT INTO Sessions(secret, user_id)
	// 	VALUES(?,?)`, secret, userId)

	// if err != nil {
	// 	revel.INFO.Print("Error in createSession: ", err)
	// } else {
	// 	var id int64
	// 	id, err = queryRes.LastInsertId()
	// 	user = Session{Secret: secret, UserId: userId, Id: int(id)}
	// 	return session, true
	// }

	return session, false
}

func (s *Session) GetUser() (User, bool) {
	return FindUser(s.userId)
}

func (s *Session) Remove() {

}
