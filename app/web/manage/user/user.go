package manageuser

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/entities"
	e "github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/services"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/views"
)

func Index(c server.Context) error {
	c.Meta().Title = "Manage users"
	status := http.StatusOK
	page := c.QueryInt("page", 1)
	search := c.Query("q")
	data, err := repositories.User.Paginate(c.Context(), &e.UserFilter{Filter: &entities.Filter{Page: page, Search: search}})

	if err != nil {
		status = http.StatusBadRequest
		c.WithError("Error getting users", err)
	}

	return c.Status(status).Render(views.ManageUserIndex(data, search))
}

func Compose(c server.Context) (err error) {
	return composeView(c, &entities.User{}, false)
}

func Save(c server.Context) (err error) {
	var user *entities.User
	userID := c.ParamInt("id", 0)
	data := getUserSaveData(c)

	if c.Messages().HasError() {
		return composeView(c, data, true)
	}

	if userID > 0 {
		data.ID = userID
		user, err = repositories.User.Update(c.Context(), data)
	} else {
		user, err = repositories.User.Create(c.Context(), data)
	}

	if err != nil {
		c.WithError("Error saving user", err)
		return composeView(c, data, true)
	}

	return c.Redirect("/manage/users/" + strconv.Itoa(user.ID))
}

func Delete(c server.Context) error {
	user, err := getProcessingUser(c)

	if user.ID == 1 {
		return c.Status(http.StatusBadRequest).SendString("Error deleting user")
	}

	if err != nil {
		c.Logger().Error("Error deleting user", err)
		return c.Status(http.StatusBadRequest).SendString("Error deleting user")
	}

	if err := repositories.User.DeleteByID(c.Context(), user.ID); err != nil {
		c.Logger().Error("Error deleting user", err)
		return c.Status(http.StatusBadRequest).SendString("Error deleting user")
	}

	return c.Status(http.StatusOK).SendString("Success")
}

func getProcessingUser(c server.Context) (user *entities.User, err error) {
	if c.Param("id") == "new" {
		return &entities.User{}, nil
	}

	return repositories.User.ByID(c.Context(), c.ParamInt("id"))
}

func composeView(c server.Context, composeData *entities.User, isSave bool) (err error) {
	var roles []*entities.Role
	user, err := getProcessingUser(c)
	c.Meta().Title = "Create User"

	if err != nil {
		c.WithError("Query editting user error", err)
	} else {
		if !isSave {
			composeData = user
		}
	}

	if roles, err = repositories.Role.All(c.Context()); err != nil {
		c.WithError("Load roles error", err)
	}

	if user.ID > 0 {
		c.Meta().Title = "Edit User: " + user.Username
		user.RoleIDs = []int{}
		for _, role := range user.Roles {
			user.RoleIDs = append(user.RoleIDs, role.ID)
		}
	}

	return c.Render(views.ManageUserCompose(user.ID, composeData, roles, auth.Providers()))
}

func getUserSaveData(c server.Context) *entities.User {
	var err error
	user := &entities.User{}
	data := &entities.UserMutation{}
	if err = c.BodyParser(data); err != nil {
		c.WithError("Error parsing body", err)
		return &entities.User{}
	}

	user.Username = utils.SanitizePlainText(strings.TrimSpace(data.Username))
	user.DisplayName = utils.SanitizePlainText(strings.TrimSpace(data.DisplayName))
	user.Email = utils.SanitizePlainText(strings.TrimSpace(data.Email))
	user.URL = utils.SanitizePlainText(strings.TrimSpace(data.URL))
	user.Bio = utils.SanitizeMarkdown(strings.TrimSpace(data.Bio))
	user.Provider = utils.SanitizePlainText(strings.TrimSpace(data.Provider))
	user.ProviderID = utils.SanitizePlainText(strings.TrimSpace(data.ProviderID))
	user.ProviderUsername = utils.SanitizePlainText(strings.TrimSpace(data.ProviderUsername))
	user.ProviderAvatar = utils.SanitizePlainText(strings.TrimSpace(data.ProviderAvatar))
	user.Password = utils.SanitizePlainText(strings.TrimSpace(data.Password))
	user.RoleIDs = data.RoleIDs
	user.Active = data.Active

	if avatarImage, err := services.SaveFile(c, "avatar_image"); err != nil {
		c.WithError("Error saving avatar image", err)
	} else if avatarImage != nil {
		user.AvatarImageID = avatarImage.ID
	}

	if data.Username == "" || len(data.Username) > 250 {
		c.Messages().AppendError("Username is required and can't be more than 250 characters")
	}

	if data.Password != "" {
		if user.Password, err = utils.GenerateHash(data.Password); err != nil {
			c.WithError("Error generating password hash", err)
			return user
		}
	}

	return user
}
