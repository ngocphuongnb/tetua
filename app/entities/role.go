package entities

import (
	"net/url"
	"time"

	"github.com/ngocphuongnb/tetua/app/utils"
)

// Role is the model entity for the Role schema.
type Role struct {
	ID          int           `json:"id,omitempty"`
	CreatedAt   *time.Time    `json:"created_at,omitempty"`
	UpdatedAt   *time.Time    `json:"updated_at,omitempty"`
	DeletedAt   *time.Time    `json:"deleted_at,omitempty"`
	Name        string        `json:"name,omitempty" validate:"max=255"`
	Description string        `json:"description,omitempty" validate:"max=255"`
	Root        bool          `json:"root,omitempty"`
	Users       []*User       `json:"users,omitempty"`
	Permissions []*Permission `json:"permissions,omitempty"`
}

type PermType string

const (
	PERM_ALL  PermType = "all"
	PERM_OWN  PermType = "own"
	PERM_NONE PermType = "none"
)

func (s PermType) String() string {
	switch s {
	case PERM_ALL:
		return "all"
	case PERM_OWN:
		return "own"
	case PERM_NONE:
		return "none"
	}
	return "none"
}

type RoleMutation struct {
	Name        string             `form:"name" json:"name"`
	Description string             `form:"description" json:"description"`
	Root        bool               `form:"root" json:"root"`
	Permissions []*PermissionValue `form:"permissions" json:"permissions"`
}

type PermissionValue struct {
	Action string   `json:"action,omitempty" validate:"max=255"`
	Value  PermType `json:"value,omitempty"`
}

type RolePermissions struct {
	RoleID      int                `json:"role_id,omitempty"`
	Permissions []*PermissionValue `json:"permissions,omitempty"`
}

func GetPermTypeValue(value string) PermType {
	switch value {
	case "all":
		return PERM_ALL
	case "own":
		return PERM_OWN
	case "none":
		return PERM_NONE
	default:
		return PERM_NONE
	}
}

type RoleFilter struct {
	*Filter
}

func (p *RoleFilter) Base() string {
	q := url.Values{}
	if !utils.SliceContains(p.IgnoreUrlParams, "search") && p.Search != "" {
		q.Add("q", p.Search)
	}

	if queryString := q.Encode(); queryString != "" {
		return p.FilterBaseUrl() + "?" + q.Encode()
	}

	return p.FilterBaseUrl()
}
