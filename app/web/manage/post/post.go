package managepost

import (
	"fmt"
	"net/http"

	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	e "github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/views"
)

func Index(c server.Context) error {
	page := c.QueryInt("page", 1)
	search := c.Query("q")
	publish := c.Query("publish", "all")
	approve := c.Query("approve", "all")
	topicID := c.QueryInt("topic", 0)
	userID := c.QueryInt("user", 0)
	topicIDs := []int{}
	userIDs := []int{}
	status := http.StatusOK
	topics, err := repositories.Topic.All(c.Context())

	fmt.Println(approve)

	if topicID > 0 {
		topicIDs = append(topicIDs, topicID)
	}
	if userID > 0 {
		userIDs = append(userIDs, userID)
	}

	if err != nil {
		c.WithError("Error getting topics", err)
	}

	data, err := repositories.Post.Paginate(c.Context(), &e.PostFilter{
		Publish:  publish,
		Approve:  approve,
		TopicIDs: topicIDs,
		UserIDs:  userIDs,
		Filter: &entities.Filter{
			BaseUrl: config.Url("/manage/posts"),
			Page:    page,
			Search:  search,
		}})
	c.Meta().Title = "Manage posts"

	if err != nil {
		status = http.StatusBadRequest
		c.WithError("Error getting posts", err)
	}

	return c.Status(status).Render(views.ManagePostIndex(
		data,
		entities.PrintTopicsTree(topics, []int{}),
		topicIDs,
		search,
		publish,
		approve,
	))
}

func Approve(c server.Context) error {
	if err := repositories.Post.Approve(c.Context(), c.ParamInt("id")); err != nil {
		c.Logger().Error("Error aprrove post", err)
		return c.Status(http.StatusBadRequest).Json(&entities.Message{
			Type:    "error",
			Message: "Error aprrove post",
		})
	}

	return c.Status(http.StatusOK).Json(&entities.Message{
		Type:    "success",
		Message: "Post aprroved",
	})
}
