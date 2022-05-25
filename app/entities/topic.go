package entities

import (
	"net/url"
	"strings"
	"time"

	"github.com/ngocphuongnb/tetua/app/utils"
)

// Topic is the model entity for the Topic schema.
type Topic struct {
	ID          int        `json:"id,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Name        string     `json:"name,omitempty" validate:"max=255"`
	Slug        string     `json:"slug,omitempty" validate:"max=255"`
	Description string     `json:"description,omitempty" validate:"max=255"`
	Content     string     `json:"content,omitempty" validate:"required"`
	ContentHTML string     `json:"content_html,omitempty" validate:"required"`
	ParentID    int        `json:"parent_id,omitempty"`
	Parent      *Topic     `json:"parent,omitempty"`
	Children    []*Topic   `json:"children,omitempty"`
	Posts       []*Post    `json:"posts,omitempty"`
}

type TopicFilter struct {
	*Filter
}

type TopicMutation struct {
	ID          int    `form:"id" json:"id"`
	Name        string `form:"name" json:"name"`
	Content     string `form:"content" json:"content"`
	ContentHTML string `json:"content_html,omitempty" validate:"required"`
	ParentID    int    `form:"parent_id" json:"parent_id"`
}

func (t *Topic) Url() string {
	return utils.Url(t.Slug)
}

func (t *Topic) FeedUrl() string {
	return utils.Url(t.Slug + "/feed")
}

func GetTopicsTree(topics []*Topic, rootTopic, level int, ignore []int) []*Topic {
	var result []*Topic

	for _, topic := range topics {
		if utils.SliceContains(ignore, topic.ID) {
			continue
		}
		if topic.ParentID == rootTopic {
			topic.Children = GetTopicsTree(topics, topic.ID, level+1, ignore)
			result = append(result, topic)
		}
	}

	return result
}

func topicTree(topics []*Topic, level int, ignore []int) []*Topic {
	var result []*Topic
	for _, topic := range topics {
		if utils.SliceContains(ignore, topic.ID) {
			continue
		}

		topic.Name = strings.Repeat("--", level) + topic.Name
		result = append(result, topic)

		if len(topic.Children) > 0 {
			result = append(result, topicTree(topic.Children, level+1, ignore)...)
		}
	}

	return result
}

func PrintTopicsTree(topics []*Topic, ignore []int) []*Topic {
	ts := GetTopicsTree(topics, 0, 0, []int{})
	topics = topicTree(ts, 0, ignore)
	return topics
}

func (p *TopicFilter) Base() string {
	q := url.Values{}
	if !utils.SliceContains(p.IgnoreUrlParams, "search") && p.Search != "" {
		q.Add("q", p.Search)
	}

	if queryString := q.Encode(); queryString != "" {
		return p.FilterBaseUrl() + "?" + q.Encode()
	}

	return p.FilterBaseUrl()
}
