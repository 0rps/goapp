package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/0rps/goapp/app"
	"github.com/revel/revel"
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

func CreateUser(login, pass string) (User, error) {
	var user User
	var err error
	queryRes, queryErr := app.DB.Exec(`
		INSERT INTO Users(login, password)
		VALUES(?,?)`, login, pass)

	if queryErr != nil {
		revel.INFO.Print("Error in createuser: ", queryErr)
		err = errors.New("DB error")
	} else {
		var id int64
		id, _ = queryRes.LastInsertId()
		user = User{Login: login, Password: pass, Id: int(id)}
	}

	return user, err
}

func FindUserById(id int) (User, error) {
	var (
		user     User
		err      error
		login    string
		password string
	)

	queryErr := app.DB.QueryRow(`
		SELECT login, password FROM Users 
		WHERE id=?`, id).Scan(&login, &password)

	switch {
	case queryErr == sql.ErrNoRows:
		err = errors.New("No such user")
	case err != nil:
		err = queryErr
	default:
		user = User{Id: id, Login: login, Password: password}
	}

	return user, err
}

func FindUserByLogin(login string) (User, error) {
	var (
		user     User
		err      error
		id       int
		password string
	)

	queryErr := app.DB.QueryRow(`
		SELECT id, password FROM Users 
		WHERE login LIKE ?`, login).Scan(&id, &password)

	switch {
	case queryErr == sql.ErrNoRows:
		err = errors.New("No such user")
	case err != nil:
		err = queryErr
	default:
		user = User{Id: id, Login: login, Password: password}
	}

	return user, err
}

func Authorize(login, pass string) (Session, bool) {
	var session Session

	user, err := FindUserByLogin(login)
	if err == nil {

		if user.Password == pass {
			session, err = createSession(user.Id)
			if err == nil {
				return session, true
			} else {
				revel.INFO.Println("cannot authorize, cannot create session, err: ", err)
			}
		} else {
			revel.INFO.Printf("cannot authorize, diff passwords: %s and %s", pass, user.Password)
		}
	} else {
		revel.INFO.Printf("cannot authorize, no such user with login: %s, err = %s", login, err.Error())
	}

	return session, false
}

func GetSession(id int, secret string) (Session, error) {
	var (
		session Session
		err     error
	)

	queryErr := app.DB.QueryRow(`
		SELECT secret, user_id FROM Sessions 
		WHERE id=?`, id).Scan(&session.Secret, &session.UserId)

	switch {
	case queryErr == sql.ErrNoRows:
		err = errors.New(fmt.Sprintf("No such session with id: %d", id))
	case err != nil:
		err = queryErr
	default:
		//	revel.INFO.Printf("sql session id=%d, sec=%s, uid = %d", id, ssecret, sid)
		if session.Secret == secret {
			session.Id = id
		} else {
			session = Session{}
			err = errors.New(fmt.Sprintf("Wrong secret for session with id: %d", id))
		}
	}

	return session, err
}

func generateSecret() string {
	return "12345667"
}

func createSession(userId int) (Session, error) {
	var session Session
	var err error

	secret := generateSecret()
	queryRes, queryErr := app.DB.Exec(`
		INSERT INTO Sessions(secret, user_id)
		VALUES(?,?)`, secret, userId)

	if queryErr != nil {
		revel.INFO.Print("Error in createSesson: ", queryErr)
		err = errors.New("DB error")
	} else {
		var id int64
		id, _ = queryRes.LastInsertId()
		session = Session{Id: int(id), Secret: secret, UserId: userId}
	}

	return session, err
}

func (s *Session) GetUser() User {
	user, _ := FindUserById(s.UserId)
	return user
}

func (s *Session) Remove() error {

	var err error
	_, queryErr := app.DB.Exec(`
		DELETE FROM Sessions WHERE id=?`, s.Id)

	if queryErr != nil {
		revel.INFO.Print("Error in removeSession: ", queryErr)
		err = errors.New("DB error")
	}

	return err
}
