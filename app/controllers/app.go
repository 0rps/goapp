package controllers

import (
	"errors"
	"github.com/0rps/goapp/app/chatroom"
	"github.com/0rps/goapp/app/models"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
	"strconv"
)

type App struct {
	*revel.Controller
}

func (c *App) getSession() (models.Session, error) {
	var session models.Session
	var err error

	id, ok1 := c.Session["sessionid"]
	secret, ok2 := c.Session["sessionsecret"]

	if ok1 && ok2 {
		sessionId, _ := strconv.ParseInt(id, 10, 32)
		session, modelsErr := models.GetSession(int(sessionId), secret)

		if err == nil {
			return session, nil
		}

		err = modelsErr

	} else {
		err = errors.New("couldnot get session from cookie")
	}

	return session, err
}

func (c *App) loginWrapper() {
	session, err := c.getSession()

	if err == nil {
		revel.INFO.Println("user is logged in")
		c.RenderArgs["loggedIn"] = true
		c.RenderArgs["userId"] = session.UserId
	} else {
		revel.INFO.Println("user is not logged in, err", err)

		c.RenderArgs["loggedIn"] = false
		c.RenderArgs["userId"] = -1
	}
}

func (c *App) userId() int {
	return c.RenderArgs["userId"].(int)
}

func (c *App) isLoggedIn() bool {
	return c.RenderArgs["loggedIn"].(bool)
}

func (c App) Login(login, password string) revel.Result {
	revel.INFO.Printf("login is '%s', pass is '%s'", login, password)

	c.loginWrapper()

	if c.isLoggedIn() {
		return c.Redirect("/")
	}

	if login != "" || password != "" {
		session, noerror := models.Authorize(login, password)
		if noerror {
			c.Session["sessionid"] = strconv.FormatInt(int64(session.Id), 10)
			c.Session["sessionsecret"] = session.Secret

			revel.INFO.Printf("session is installed,%s,%s", c.Session["sessionid"], c.Session["sessionsecret"])

			c.Flash.Success("Пользователь авторизован")
			return c.Redirect("/")
		} else {
			c.Flash.Error("Не существует пользователя или неправильный пароль")
			return c.Redirect("/login")
		}
	}

	return c.Render("App/Login.html")
}

func (c App) Register(login, password, repassword string) revel.Result {
	revel.INFO.Printf("login is '%s', pass is '%s', repass is '%s'", login, password, repassword)

	c.loginWrapper()

	if c.isLoggedIn() {
		return c.Redirect("/")
	}

	if login != "" || password != "" || repassword != "" {
		c.Validation.Required(login).Message("Введите логин")
		c.Validation.Required(password).Message("Введите пароль")
		c.Validation.MinSize(password, 4).Message("Минимальная длина пароля - 4 символа")

		if c.Validation.HasErrors() {
			c.Validation.Keep()
			c.FlashParams()
			return c.Redirect("/register")
		}

		if password != repassword {
			c.Flash.Error("Пароли не совпадают")
			return c.Redirect("/register")
		}

		if _, err := models.FindUserByLogin(login); err == nil {
			c.Flash.Error("Пользователь с таким именем существует")
			return c.Redirect("/register")
		}

		user, err := models.CreateUser(login, password)
		if err == nil {
			c.Flash.Success("Пользователь создан, id = ", user.Id)
			return c.Redirect("/")
		} else {
			revel.INFO.Println(err)
		}
	}

	return c.RenderTemplate("App/Register.html")
}

func (c App) Index() revel.Result {
	c.loginWrapper()
	return c.RenderTemplate("App/Index.html")
}

func (c App) Logout() revel.Result {
	c.loginWrapper()

	if c.isLoggedIn() {
		session, _ := c.getSession()
		err := session.Remove()
		if err == nil {
			delete(c.Session, "sessionid")
			delete(c.Session, "sessionsecret")
		} else {
			revel.INFO.Println("couldnot logout, ", err)
		}
	}

	return c.Redirect("/")
}

func (c App) Rooms() revel.Result {
	return c.Redirect("/")
}

func (c App) Config() revel.Result {
	return c.Redirect("/")
}

func (c App) Room() revel.Result {
	id := 1
	revel.INFO.Printf("room request, id = ", id)

	c.loginWrapper()

	if !c.isLoggedIn() {
		return c.Redirect("/")
	}

	c.RenderArgs["moreScripts"] = []string{"js/react.js", "js/react-dom.js"}

	c.RenderArgs["moreBabelScripts"] = []string{"js/chat.js"}
	return c.RenderTemplate("App/Room.html")
}

func (c App) RoomSocket(ws *websocket.Conn) revel.Result {
	roomId := 1
	c.loginWrapper()
	if !c.isLoggedIn() {
		return nil
	}

	subscription := chatroom.Subscribe(roomId)
	defer subscription.Cancel()

	user, _ := models.FindUserById(c.userId())

	chatroom.Join(roomId, user.Login)
	defer chatroom.Leave(roomId, user.Login)

	// for _, msg := range models.GetArchiveMessages(roomId) {
	// 	if websocket.JSON.Send(ws, &msg) != nil {
	// 		return nil
	// 	}
	// }

	newMessages := make(chan string)
	go func() {
		var msg string
		for {
			revel.INFO.Println("try to reveive socket message")
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				close(newMessages)
				return
			}
			newMessages <- msg
		}
	}()

	// Now listen for new events from either the websocket or the chatroom.
	for {
		select {
		case event := <-subscription.New:
			revel.INFO.Println("main loop: new message: ", event.Text)

			if websocket.JSON.Send(ws, &event) != nil {
				// They disconnected.
				return nil
			}
		case msg, ok := <-newMessages:
			revel.INFO.Println("main loop: try to say to all: ", msg)

			// If the channel is closed, they disconnected.
			if !ok {
				return nil
			}

			// Otherwise, say something.
			chatroom.Say(roomId, user.Login, msg)
		}
	}
	return nil

}
