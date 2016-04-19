package controllers

import "github.com/revel/revel"
import "github.com/0rps/goapp/app/models"

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

func (c App) Register(login, password, repassword string) revel.Result {
	revel.INFO.Printf("login is '%s', pass is '%s', repass is '%s'", login, password, repassword)

	if login != "" || password != "" || repassword != "" {
		c.Validation.Required(login).Message("Введите логин")
		c.Validation.Required(password).Message("Введите пароль")
		c.Validation.MinSize(password, 4).Message("Минимальная длина пароля - 8 символов")

		if c.Validation.HasErrors() {
			c.Validation.Keep()
			c.FlashParams()
			return c.Redirect("/register")
		}

		if password != repassword {
			c.Flash.Error("Пароли не совпадают")
			return c.Redirect("/register")
		}

		if models.IsUserExists(login) {
			c.Flash.Error("Пользователь с таким именем существует")
			return c.Redirect("/register")
		}

		user := models.CreateUser(login, password)
		if user.Id > 0 {
			c.Flash.Success("Пользователь создан, id = ", user.Id)
			return c.Redirect("/")
		}
	}

	return c.RenderTemplate("App/Register.html")
}

func (c App) Index() revel.Result {
	return c.RenderTemplate("App/Index.html")
}
