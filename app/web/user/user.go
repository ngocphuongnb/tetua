package webuser

import (
	"net/http"

	"time"

	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/views"
)

type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Login(c server.Context) (err error) {
	if c.User() != nil && c.User().ID > 0 {
		return c.Redirect(utils.Url(""))
	}
	c.Meta().Title = "Login"
	return c.Render(views.Login())
}

func PostLogin(c server.Context) (err error) {
	loginData := &LoginData{}
	if err := c.BodyParser(loginData); err != nil {
		c.Logger().Error(err)
		c.Messages().AppendError("Something went wrong")
		return c.Render(views.Login())
	}

	foundUsers, err := repositories.User.ByUsernameOrEmail(c.Context(), loginData.Login, loginData.Login)

	if err != nil {
		c.Logger().Error(err)
		c.Messages().AppendError("Something went wrong")
		return c.Render(views.Login())
	}

	if len(foundUsers) == 0 {
		c.Messages().AppendError("Invalid login information")
		return c.Render(views.Login())
	}

	if err = utils.CheckHash(loginData.Password, foundUsers[0].Password); err != nil {
		c.Messages().AppendError("Invalid login information")
		return c.Render(views.Login())
	}

	if !foundUsers[0].IsRoot() && !foundUsers[0].Active {
		return c.Redirect(utils.Url("/inactive"))
	}

	if err = auth.SetLoginInfo(c, foundUsers[0]); err != nil {
		c.Logger().Error("Error setting login info", err)
		return c.Status(http.StatusBadGateway).SendString("Something went wrong")
	}

	return c.Redirect(utils.Url(""))
}

func Inactive(c server.Context) (err error) {
	return c.Render(views.Inactive())
}

func Logout(c server.Context) (err error) {
	c.Cookie(&server.Cookie{
		Name:    config.APP_TOKEN_KEY,
		Value:   "",
		Expires: time.Now().Add(time.Hour * 100 * 365 * 24),
	})

	return c.Redirect("/")
}
