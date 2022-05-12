package webuser

import (
	"net/http"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
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
		return c.Redirect(config.Url(""))
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
		spew.Dump(err)
		c.Messages().AppendError("Invalid login information")
		return c.Render(views.Login())
	}

	if !foundUsers[0].IsRoot() && !foundUsers[0].Active {
		return c.Redirect(config.Url("/inactive"))
	}

	if err = auth.SetLoginInfo(c, foundUsers[0]); err != nil {
		c.Logger().Error("Error setting login info", err)
		return c.Status(http.StatusBadGateway).SendString("Something went wrong")
	}

	return c.Redirect(config.Url(""))
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

func Profile(c server.Context) (err1 error) {
	username := c.Param("username")
	user, err1 := repositories.User.ByUsername(c.Context(), username)
	var err2 error
	var paginate *entities.Paginate[entities.Post]
	var comments []*entities.Comment
	var wg sync.WaitGroup

	if err1 != nil {
		return c.Render(views.Error("User not found"))
	}

	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		paginate, err1 = repositories.Post.Paginate(c.Context(), &entities.PostFilter{
			UserIDs: []int{user.ID},
			Filter: &entities.Filter{
				Page:            c.ParamInt("page", 1),
				IgnoreUrlParams: []string{"user"},
			},
		})

		if err1 != nil {
			c.Logger().Error(err1)
		}
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		if c.ParamInt("page", 1) > 1 {
			return
		}
		comments, err2 = repositories.Comment.FindWithPost(c.Context(), &entities.CommentFilter{
			UserIDs: []int{user.ID},
			Filter: &entities.Filter{
				Limit: 5,
			},
		})
		if err1 != nil {
			c.Logger().Error(err1)
		}
	}(&wg)

	wg.Wait()

	if err := utils.FirstError(err1, err2); err != nil {
		c.Logger().Error("Error loading profile", err)
		return c.Status(http.StatusBadRequest).Render(views.Error("Error loading profile"))
	}

	c.Meta().Title = user.Name()
	return c.Render(views.Profile(user, paginate, comments))
}
