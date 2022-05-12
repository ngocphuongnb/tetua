package entities

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/utils"
)

// Post is the model entity for the Post schema.
type Post struct {
	ID              int        `json:"id,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	Name            string     `json:"name,omitempty" validate:"max=255"`
	Slug            string     `json:"slug,omitempty" validate:"max=255"`
	Description     string     `json:"description,omitempty" validate:"max=255"`
	Content         string     `json:"content,omitempty" validate:"required"`
	ContentHTML     string     `json:"content_html,omitempty"`
	ViewCount       int64      `json:"view_count,omitempty"`
	CommentCount    int64      `json:"comment_count,omitempty"`
	RatingCount     int64      `json:"rating_count,omitempty"`
	RatingTotal     int64      `json:"rating_total,omitempty"`
	Draft           bool       `json:"draft,omitempty"`
	Approved        bool       `json:"approved,omitempty"`
	FeaturedImageID int        `json:"featured_image_id,omitempty"`
	UserID          int        `json:"user_id,omitempty"`
	User            *User      `json:"user,omitempty"`
	FeaturedImage   *File      `json:"featured_image,omitempty"`
	Topics          []*Topic   `json:"topics,omitempty"`
	TopicIDs        []int      `json:"topic_ids,omitempty"`
}

type PostMutation struct {
	Name            string `form:"name" json:"name"`
	Slug            string `form:"name" json:"slug"`
	Description     string `form:"description" json:"description"`
	Content         string `form:"content" json:"content"`
	ContentHTML     string `form:"content_html" json:"content_html"`
	TopicIDs        []int  `form:"topic_ids" json:"topic_ids"`
	Draft           bool   `form:"draft" json:"draft"`
	FeaturedImageID int    `form:"featured_image_id" json:"featured_image_id"`
}

type PostFilter struct {
	*Filter
	Approve  string `form:"approve" json:"approve"`           // approve = all, approved, pending
	Publish  string `form:"publish_type" json:"publish_type"` // publish_type = all, published, draft
	UserIDs  []int  `form:"user_ids" json:"user_ids"`
	TopicIDs []int  `form:"topic_ids" json:"topic_ids"`
}

func (p *Post) Url() string {
	return config.Url(fmt.Sprintf("%s-%d.html", p.Slug, p.ID))
}

func (p *PostFilter) Base() string {
	q := url.Values{}
	if !utils.SliceContains(p.IgnoreUrlParams, "search") && p.Search != "" {
		q.Add("q", p.Search)
	}
	if !utils.SliceContains(p.IgnoreUrlParams, "topic") && len(p.TopicIDs) > 0 {
		q.Add("topic", strconv.Itoa(p.TopicIDs[0]))
	}
	if !utils.SliceContains(p.IgnoreUrlParams, "user") && len(p.UserIDs) > 0 {
		q.Add("user", strconv.Itoa(p.UserIDs[0]))
	}
	if !utils.SliceContains(p.IgnoreUrlParams, "publish") && p.Publish != "" {
		q.Add("publish", p.Publish)
	}
	if !utils.SliceContains(p.IgnoreUrlParams, "approve") && p.Approve != "" {
		q.Add("approve", p.Approve)
	}

	if queryString := q.Encode(); queryString != "" {
		return p.FilterBaseUrl() + "?" + q.Encode()
	}

	return p.FilterBaseUrl()
}
