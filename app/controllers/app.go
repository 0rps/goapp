package controllers

import "github.com/revel/revel"
import "myapp/app/models"

type App struct {
	*revel.Controller
}

func (c App) Login(login, password string) revel.Result {
	revel.INFO.Printf("login is '%s', pass is '%s'", login, password)

	if login != "" || password != "" {
		_, noerror := models.Authorize(login, password)
		if noerror {
			c.Flash.Success("Пользователь авторизован")
			return c.Redirect("/")
		} else {
			c.Flash.Error("Не существует пользователя или неправильный пароль")
			return c.Redirect("/login")
		}
	}

	return c.Render("App/Login.html")
}

func (c App) Register() revel.Result {
	return c.RenderTemplate("App/Index.html")
}

func (c App) Index() revel.Result {
	return c.RenderTemplate("App/Index.html")
}
