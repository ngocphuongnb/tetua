package managetopic

import (
	"net/http"
	"strconv"

	"github.com/gosimple/slug"
	"github.com/ngocphuongnb/tetua/app/cache"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/views"
)

func Index(c server.Context) (err error) {
	status := http.StatusOK
	topics, err := repositories.Topic.All(c.Context())
	c.Meta().Title = "Manage topics"

	if err != nil {
		status = http.StatusBadRequest
		c.WithError("Load all topics error", err)
	}

	return c.Status(status).Render(views.ManageTopicIndex(entities.PrintTopicsTree(topics, []int{})))
}

func Compose(c server.Context) (err error) {
	return getTopicComposeView(c, &entities.TopicMutation{}, false)
}

func Save(c server.Context) (err error) {
	var topic *entities.Topic
	composeData := getTopicSaveData(c)

	if topic, err = getProcessingTopic(c); err != nil {
		c.WithError("Query editting topic error", err)
	}

	if c.Messages().Length() > 0 {
		return getTopicComposeView(c, composeData, true)
	}

	topic.Name = composeData.Name
	topic.Content = composeData.Content
	topic.ParentID = composeData.ParentID
	topic.ContentHTML = composeData.ContentHTML

	if topic.ID > 0 {
		topic, err = repositories.Topic.Update(c.Context(), topic)
	} else {
		topic.Slug = slug.Make(composeData.Name)
		topic, err = repositories.Topic.Create(c.Context(), topic)
	}

	if err != nil {
		c.WithError("Error saving topic", err)
		return getTopicComposeView(c, composeData, true)
	}

	if err := cache.CacheTopics(c.Context()); err != nil {
		c.WithError("Error caching topics", err)
		return getTopicComposeView(c, composeData, true)
	}

	return c.Redirect("/manage/topics/" + strconv.Itoa(topic.ID))
}

func Delete(c server.Context) error {
	topic, err := getProcessingTopic(c)

	if err != nil {
		c.Logger().Error("Error deleting topic", err)
		return c.Status(http.StatusBadRequest).SendString("Error deleting topic")
	}

	if err := repositories.Topic.DeleteByID(c.Context(), topic.ID); err != nil {
		c.Logger().Error("Error deleting topic", err)
		return c.Status(http.StatusBadRequest).SendString("Error deleting topic")
	}

	return c.Status(http.StatusOK).SendString("Topic deleted")
}

func getProcessingTopic(c server.Context) (topic *entities.Topic, err error) {
	if c.Param("id") == "new" {
		return &entities.Topic{}, nil
	}
	return repositories.Topic.ByID(c.Context(), c.ParamInt("id"))
}

func getTopicComposeView(c server.Context, data *entities.TopicMutation, isSave bool) (err error) {
	var topics []*entities.Topic
	ignore := make([]int, 0)
	topic, err := getProcessingTopic(c)
	c.Meta().Title = "Create Topic"

	if err != nil {
		c.WithError("Query editting topic error", err)
	} else if !isSave {
		data.ID = topic.ID
		data.Name = topic.Name
		data.Content = topic.Content
		data.ParentID = topic.ParentID
	}

	if topics, err = repositories.Topic.All(c.Context()); err != nil {
		c.WithError("Load all topics error", err)
	}

	if topic.ID > 0 {
		c.Meta().Title = "Edit Topic: " + topic.Name
		ignore = append(ignore, topic.ID)
	}

	return c.Render(views.ManageTopicCompose(entities.PrintTopicsTree(topics, ignore), data))
}

func getTopicSaveData(c server.Context) *entities.TopicMutation {
	var err error
	data := &entities.TopicMutation{}
	if err = c.BodyParser(data); err != nil {
		c.WithError("Error parsing body", err)

		return data
	}

	data.Name = utils.SanitizePlainText(data.Name)
	data.Content = utils.SanitizeMarkdown(data.Content)

	if data.Name == "" || len(data.Name) > 250 {
		c.Messages().AppendError("Name is required and can't be more than 250 characters")
	}

	if data.ContentHTML, err = utils.MarkdownToHtml(data.Content); err != nil {
		c.WithError("Error convert markdown to html", err)
	}

	return data
}
