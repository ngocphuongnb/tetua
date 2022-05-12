package entities

import (
	"net/url"
	"strconv"
	"time"

	"github.com/ngocphuongnb/tetua/app/utils"
)

// Permission is the model entity for the Permission schema.
type Permission struct {
	ID        int        `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"Deleted_at,omitempty"`
	RoleID    int        `json:"role_id,omitempty"`
	Action    string     `json:"action,omitempty" validate:"max=255"`
	Value     string     `json:"value,omitempty"`
	Role      *Role      `json:"role,omitempty"`
}

type PermissionFilter struct {
	*Filter
	RoleIDs []int `form:"role_ids" json:"role_ids"`
}

func (p *PermissionFilter) Base() string {
	q := url.Values{}
	if !utils.SliceContains(p.IgnoreUrlParams, "search") && p.Search != "" {
		q.Add("q", p.Search)
	}
	if !utils.SliceContains(p.IgnoreUrlParams, "role") && len(p.RoleIDs) > 0 {
		q.Add("role", strconv.Itoa(p.RoleIDs[0]))
	}

	if queryString := q.Encode(); queryString != "" {
		return p.FilterBaseUrl() + "?" + q.Encode()
	}

	return p.FilterBaseUrl()
}
