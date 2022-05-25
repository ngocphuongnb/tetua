package webcomment

import (
	"fmt"
	"net/http"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/views"
)

func List(c server.Context) error {
	paginate, err := repositories.Comment.PaginateWithPost(c.Context(), &entities.CommentFilter{
		Filter: &entities.Filter{
			BaseUrl:         utils.Url("/comments"),
			Page:            c.QueryInt("page"),
			IgnoreUrlParams: []string{"user"},
		},
		UserIDs: []int{c.User().ID},
	})

	if err != nil {
		c.WithError("Something went wrong", err)
	}

	return c.Render(views.CommentList(paginate))
}

func Save(c server.Context) (err error) {
	data := getCommentSaveData(c)
	postID := data.PostID

	if c.Messages().HasError() {
		return c.Status(http.StatusBadRequest).Json(c.Messages())
	}

	if data.ID == 0 {
		data.UserID = c.User().ID
		data, err = repositories.Comment.Create(c.Context(), data)
	} else {
		data, err = repositories.Comment.Update(c.Context(), data)
	}

	if err != nil {
		c.WithError("Error saving comment", err)
		return c.Status(http.StatusInternalServerError).SendString("Error saving comment")
	}

	return c.Redirect(fmt.Sprintf("/post-%d.html#comment-%d", postID, data.ID))
}

func Delete(c server.Context) error {
	if err := repositories.Comment.DeleteByID(c.Context(), c.ParamInt("id")); err != nil {
		c.Logger().Error("Error deleting comment", err)
		return c.Status(http.StatusBadRequest).Json(&entities.Message{
			Type:    "error",
			Message: "Error deleting comment",
		})
	}

	return c.Status(http.StatusOK).Json(&entities.Message{
		Type:    "success",
		Message: "Post deleted",
	})
}

func getCommentSaveData(c server.Context) *entities.Comment {
	var err error
	data := &entities.CommentMutation{}
	if err = c.BodyParser(data); err != nil {
		c.WithError("Error parsing body", err)

		return &entities.Comment{}
	}

	post, err := repositories.Post.ByID(c.Context(), data.PostID)

	if post == nil || err != nil {
		c.WithError("Invalid post ID", err)
	}

	if data.ParentID > 0 {
		parent, err := repositories.Comment.ByID(c.Context(), data.ParentID)

		if parent == nil || err != nil {
			c.WithError("Invalid parent ID", err)
		}
	}

	if data.Content == "" {
		c.Messages().AppendError("Content is required")
	}

	comment := &entities.Comment{
		ID:       c.ParamInt("id"),
		Content:  utils.SanitizeMarkdown(data.Content),
		PostID:   data.PostID,
		ParentID: data.ParentID,
	}

	if comment.ContentHTML, err = utils.MarkdownToHtml(data.Content); err != nil {
		c.WithError("Error convert markdown to html", err)
	}

	return comment
}
