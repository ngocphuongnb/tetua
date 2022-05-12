package entities

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/utils"
)

// Entities are used in all other parts. This will store properties of business objects and associated methods. Example: Article, User

type Entity interface {
	Comment | File | Permission | Post | Role | Setting | Topic | User
}

type EntityFilter interface {
	PostFilter | FileFilter | CommentFilter | UserFilter | PermissionFilter | RoleFilter | TopicFilter
}

type NotFoundError struct {
	Message string
}

// Error implements the error interface.
func (e *NotFoundError) Error() string {
	return e.Message
}

// IsNotFound returns a boolean indicating whether the error is a not found error.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	var e *NotFoundError
	return errors.As(err, &e)
}

type Paginate[E Entity] struct {
	BaseUrl     string `json:"base_url"`
	QueryString string `json:"query_string"`
	Total       int    `json:"total"`
	PageSize    int    `json:"page_size"`
	PageCurrent int    `json:"page_current"`
	Data        []*E   `json:"data"`
}

type Message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Messages []*Message

type Meta struct {
	Title       string
	Description string
	Query       string
	Type        string
	Image       string
	Canonical   string
	User        *User
	Messages    *Messages
}

type Map map[string]interface{}

func (ms *Messages) Append(m *Message) {
	*ms = append(*ms, m)
}

func (ms *Messages) AppendError(m string) {
	*ms = append(*ms, &Message{
		Type:    "error",
		Message: m,
	})
}

func (ms *Messages) Length() int {
	// if ms == nil {
	// 	return 0
	// }
	return len(*ms)
}

func (ms *Messages) HasError() bool {
	// if ms == nil {
	// 	return false
	// }

	errorCount := 0

	for _, m := range *ms {
		if m.Type == "error" {
			errorCount++
		}
	}

	return errorCount > 0
}

func (ms *Messages) Get() []*Message {
	return *ms
}

type GetPostFn func(ctx context.Context, id int) (*Post, error)
type GetPostsFn func(ctx context.Context, limit, offset int) ([]*Post, error)
type PostPaginateFn func(ctx context.Context, filters ...*PostFilter) (*Paginate[Post], error)

type Filter struct {
	BaseUrl         string   `json:"base_url"`
	Search          string   `form:"search" json:"search"`
	Page            int      `form:"page" json:"page"`
	Limit           int      `form:"limit" json:"limit"`
	Sorts           []*Sort  `form:"orders" json:"orders"`
	IgnoreUrlParams []string `form:"ignore_url_params" json:"ignore_url_params"`
	ExcludeIDs      []int    `form:"exclude_ids" json:"exclude_ids"`
}

func (p *Filter) FilterBaseUrl() string {
	if p.BaseUrl == "" {
		return config.Url("")
	}

	return p.BaseUrl
}

func (p *Filter) Base() string {
	q := url.Values{}
	if !utils.SliceContains(p.IgnoreUrlParams, "search") && p.Search != "" {
		q.Add("q", p.Search)
	}

	if queryString := q.Encode(); queryString != "" {
		return p.FilterBaseUrl() + "?" + q.Encode()
	}

	return p.FilterBaseUrl()
}

func (f *Filter) GetSearch() string {
	return f.Search
}

func (f *Filter) GetPage() int {
	return f.Page
}

func (f *Filter) GetLimit() int {
	return f.Limit
}

func (f *Filter) GetSorts() []*Sort {
	return f.Sorts
}

func (f *Filter) GetIgnoreUrlParams() []string {
	return f.IgnoreUrlParams
}

func (f *Filter) GetExcludeIDs() []int {
	return f.ExcludeIDs
}

type Sort struct {
	Field string
	Order string
}

type PaginateLink struct {
	Link  string
	Label string
	Class string
}

func (pp *Paginate[E]) Links() []*PaginateLink {
	var url string
	var links = []*PaginateLink{}
	totalPages := int(math.Ceil(float64(pp.Total) / float64(pp.PageSize)))

	for i := 1; i <= totalPages; i++ {
		page := strconv.Itoa(i)
		class := []string{}

		if i == 1 {
			class = []string{"first"}
		}

		if i == totalPages {
			class = []string{"last"}
		}

		if i == pp.PageCurrent {
			class = append(class, "active")
		}

		if strings.Contains(pp.BaseUrl, "?") {
			url = pp.BaseUrl + "&page=" + page
		} else {
			url = pp.BaseUrl + "?page=" + page
		}

		links = append(links, &PaginateLink{
			Link:  url,
			Label: page,
			Class: strings.Join(class, " "),
		})
	}

	return links
}

func (m *Meta) GetTitle() string {
	title := config.Setting("app_name")

	if m.Title != "" {
		title = fmt.Sprintf("%s - %s", m.Title, title)
	}

	return title
}
