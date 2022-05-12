package webpost

import (
	"net/http"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/views"
)

func Compose(c server.Context) (err error) {
	featuredImage := &entities.File{}
	postData := entities.PostMutation{}
	c.Meta().Title = "Create Post"

	if post := c.Post(); post != nil {
		c.Meta().Title = "Edit Post: " + post.Name
		featuredImage = post.FeaturedImage
		postData.Name = post.Name
		postData.Draft = post.Draft
		postData.Content = post.Content
		postData.FeaturedImageID = post.FeaturedImageID
		postTopics := post.Topics

		for _, topic := range postTopics {
			postData.TopicIDs = append(postData.TopicIDs, topic.ID)
		}
	}

	return getComposeView(c, &postData, featuredImage)
}

func getComposeView(c server.Context, data *entities.PostMutation, featuredImage *entities.File) (err error) {
	status := http.StatusOK
	topics, err := repositories.Topic.All(c.Context())

	if err != nil {
		c.Logger().Error("Error getting topics", err)
		c.Messages().AppendError("Error getting topics")
	}

	if c.Messages().HasError() {
		status = http.StatusBadRequest
	}

	return c.Status(status).Render(views.PostCompose(entities.PrintTopicsTree(topics, []int{}), data, featuredImage))
}
