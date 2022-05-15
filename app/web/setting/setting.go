package websetting

import (
	"net/http"
	"strings"
	"time"

	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/services"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/views"
)

func Index(c server.Context) (err error) {
	var user *entities.User
	c.Meta().Title = "Settings"
	user, err = repositories.User.ByID(c.Context(), c.User().ID)

	if err != nil {
		c.WithError("Error while getting user", err)
	}

	return c.Render(views.UserSetting(user))
}

func Save(c server.Context) (err error) {
	user := c.User()
	data := getSettingSaveData(c)

	if c.Messages().HasError() {
		return c.Render(views.UserSetting(user))
	}

	if data.Username == "" || data.Email == "" {
		c.Messages().AppendError("Username and email are required")
		return c.Render(views.UserSetting(user))
	}

	existedUsers, err := repositories.User.ByUsernameOrEmail(c.Context(), data.Username, data.Email)

	if err == nil && len(existedUsers) > 0 {
		existedUsers = utils.SliceFilter(existedUsers, func(existedUser *entities.User) bool {
			return existedUser.ID != user.ID
		})

		if len(existedUsers) > 0 {
			c.Messages().AppendError("Username or email is already taken")
			return c.Render(views.UserSetting(user))
		}
	}

	user, err = repositories.User.Setting(c.Context(), user.ID, data)

	if err != nil {
		user = c.User()
		c.WithError("Error saving user", err)
		return c.Render(views.UserSetting(user))
	}

	user, err = repositories.User.ByID(c.Context(), user.ID)

	if err != nil {
		user = c.User()
		c.WithError("Error saving user", err)
		return c.Render(views.UserSetting(user))
	}

	exp := time.Now().Add(time.Hour * 100 * 365 * 24)
	jwtToken, err := user.JwtClaim(exp)

	if err != nil {
		c.Logger().Error("Error setting jwt", err)
		return c.Status(http.StatusBadRequest).Render(views.Error("Something went wrong"))
	}

	c.Cookie(&server.Cookie{
		Name:     config.APP_TOKEN_KEY,
		Value:    jwtToken,
		Expires:  exp,
		HTTPOnly: false,
		SameSite: "lax",
		Secure:   true,
	})

	return c.Redirect("/settings")
}

func getSettingSaveData(c server.Context) *entities.SettingMutation {
	var err error
	data := &entities.SettingMutation{}
	if err = c.BodyParser(data); err != nil {
		c.WithError("Error parsing body", err)
		return data
	}

	data.Username = strings.TrimSpace(data.Username)
	data.DisplayName = strings.TrimSpace(data.DisplayName)
	data.Email = strings.TrimSpace(data.Email)
	data.URL = strings.TrimSpace(data.URL)
	data.Bio = strings.TrimSpace(data.Bio)
	data.BioHTML, err = utils.MarkdownToHtml(data.Bio)

	if err != nil {
		c.WithError("Error convert markdown to html", err)
	}

	if avatarImage, err := services.SaveFile(c, "avatar_image"); err != nil {
		c.WithError("Error saving avatar image", err)
	} else if avatarImage != nil {
		data.AvatarImageID = avatarImage.ID
	}

	if data.Username == "" || len(data.Username) > 250 {
		c.Messages().AppendError("Username is required and can't be more than 250 characters")
	}

	if data.Password != "" {
		if data.Password, err = utils.GenerateHash(data.Password); err != nil {
			c.WithError("Error generating password hash", err)
			return data
		}
	}

	return data
}
