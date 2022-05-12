package managecomment

import (
	"net/http"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/views"
)

func Index(c server.Context) (err error) {
	status := http.StatusOK
	search := c.Query("q")
	postID := c.QueryInt("post")
	userID := c.QueryInt("user")
	filter := &entities.CommentFilter{
		Filter: &entities.Filter{
			Page:   c.QueryInt("page", 1),
			Search: search,
		},
	}
	if postID > 0 {
		filter.PostIDs = append(filter.PostIDs, postID)
	}
	if userID > 0 {
		filter.UserIDs = append(filter.PostIDs, userID)
	}
	paginate, err := repositories.Comment.PaginateWithPost(c.Context(), filter)
	c.Meta().Title = "Manage comments"

	if err != nil {
		status = http.StatusBadRequest
		c.WithError("Load comments error", err)
	}

	return c.Status(status).Render(views.ManageCommentIndex(paginate, search, userID, postID))
}
