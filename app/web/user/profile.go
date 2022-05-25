package webuser

import (
	"net/http"
	"sync"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/views"
)

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
