package webuser

import (
	"strconv"
	"strings"
	"time"

	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/views"
)

func Active(c server.Context) (err error) {
	code := c.Query("code")
	if code == "" {
		return c.Redirect(utils.Url(""))
	}

	activation, err := utils.Decrypt(code)
	if err != nil {
		return c.Render(views.Message("Something went wrong", "Invalid activation code.", "", 0))
	}

	parts := strings.Split(activation, "_")
	if len(parts) != 2 {
		return c.Render(views.Message("Something went wrong", "Invalid activation code.", "", 0))
	}

	userID, err := strconv.Atoi(parts[0])
	if err != nil {
		return c.Render(views.Message("Something went wrong", "Invalid activation code.", "", 0))
	}

	exp, err := strconv.ParseInt(parts[1], 10, 64)

	if err != nil {
		return c.Render(views.Message("Something went wrong", "Invalid activation code.", "", 0))
	}

	if time.Now().UnixMicro() > exp {
		return c.Render(views.Message("Something went wrong", "Activation code has expired.", "", 0))
	}

	user, err := repositories.User.ByID(c.Context(), userID)
	if err != nil {
		return c.Render(views.Message("Something went wrong", "Can't activate your account, please contact us for more information.", "", 0))
	}
	if user.Active {
		return c.Redirect(utils.Url(""))
	}
	user.Active = true
	if _, err := repositories.User.Update(c.Context(), user); err != nil {
		c.Logger().Error(err)
		return c.Render(views.Message("Something went wrong", "Can't activate your account, please contact us for more information.", "", 0))
	}

	return c.Render(views.Message("Success", "Your account has been activated, please login.", utils.Url(""), 5))
}
