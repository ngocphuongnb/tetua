package entities

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ngocphuongnb/tetua/app/utils"
)

// Page is the model entity for the Page schema.
type Page struct {
	ID              int        `json:"id,omitempty"`
	Name            string     `json:"name,omitempty" validate:"max=255"`
	Slug            string     `json:"slug,omitempty" validate:"max=255"`
	Content         string     `json:"content,omitempty" validate:"required"`
	ContentHTML     string     `json:"content_html,omitempty"`
	Draft           bool       `json:"draft,omitempty"`
	FeaturedImageID int        `json:"featured_image_id,omitempty"`
	FeaturedImage   *File      `json:"featured_image,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type PageFilter struct {
	*Filter
	Publish string `form:"publish_type" json:"publish_type"` // publish_type = all, published, draft
}

func (p *Page) Url() string {
	return utils.Url(fmt.Sprintf("/%s.html", p.Slug))
}

func (p *PageFilter) Base() string {
	q := url.Values{}
	if !utils.SliceContains(p.IgnoreUrlParams, "search") && p.Search != "" {
		q.Add("q", p.Search)
	}
	if !utils.SliceContains(p.IgnoreUrlParams, "publish") && p.Publish != "" {
		q.Add("publish", p.Publish)
	}

	if queryString := q.Encode(); queryString != "" {
		return p.FilterBaseUrl() + "?" + q.Encode()
	}

	return p.FilterBaseUrl()
}
