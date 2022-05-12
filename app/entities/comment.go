package entities

import (
	"net/url"
	"strconv"
	"time"

	"github.com/ngocphuongnb/tetua/app/utils"
)

type Comment struct {
	ID          int        `json:"id,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Content     string     `json:"content,omitempty" validate:"required"`
	ContentHTML string     `json:"content_html,omitempty" validate:"required"`
	Votes       int64      `json:"votes,omitempty"`
	PostID      int        `json:"post_id,omitempty"`
	UserID      int        `json:"user_id,omitempty"`
	ParentID    int        `json:"parent_id,omitempty"`
	Parent      *Comment
	Post        *Post
	User        *User
}

type CommentMutation struct {
	PostID   int    `json:"post_id" form:"post_id" validate:"required"`
	ParentID int    `json:"parent_id" form:"parent_id"`
	Content  string `json:"content" form:"content" validate:"required"`
}

type CommentFilter struct {
	*Filter
	PostIDs   []int `form:"post_ids" json:"post_ids"`
	UserIDs   []int `form:"user_ids" json:"user_ids"`
	ParentIDs []int `form:"parent_ids" json:"parent_ids"`
}

func (p *CommentFilter) Base() string {
	q := url.Values{}
	if !utils.SliceContains(p.IgnoreUrlParams, "search") && p.Search != "" {
		q.Add("q", p.Search)
	}
	if !utils.SliceContains(p.IgnoreUrlParams, "post") && len(p.PostIDs) > 0 {
		q.Add("post", strconv.Itoa(p.PostIDs[0]))
	}
	if !utils.SliceContains(p.IgnoreUrlParams, "user") && len(p.UserIDs) > 0 {
		q.Add("user", strconv.Itoa(p.UserIDs[0]))
	}
	if !utils.SliceContains(p.IgnoreUrlParams, "parent") && len(p.ParentIDs) > 0 {
		q.Add("parent", strconv.Itoa(p.ParentIDs[0]))
	}

	if queryString := q.Encode(); queryString != "" {
		return p.FilterBaseUrl() + "?" + q.Encode()
	}

	return p.FilterBaseUrl()
}
