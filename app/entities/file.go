package entities

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/ngocphuongnb/tetua/app/fs"
	"github.com/ngocphuongnb/tetua/app/utils"
)

type File struct {
	ID        int        `json:"id,omitempty"`
	Disk      string     `json:"disk,omitempty"`
	Path      string     `json:"path,omitempty"`
	Type      string     `json:"type,omitempty"`
	Size      int        `json:"size,omitempty"`
	UserID    int        `json:"user_id,omitempty"`
	User      *User      `json:"user,omitempty"`
	Posts     []*Post    `json:"post,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (f *File) Url() string {
	if f.Disk == "" || f.Path == "" {
		return ""
	}

	fileDisk := fs.Disk(f.Disk)
	if fileDisk == nil {
		return ""
	}

	return fileDisk.Url(f.Path)
}

func (f *File) Delete(ctx context.Context) error {
	if f.Disk == "" || f.Path == "" {
		return errors.New("disk or path is empty")
	}

	fileDisk := fs.Disk(f.Disk)
	if fileDisk == nil {
		return errors.New("disk not found")
	}

	return fileDisk.Delete(ctx, f.Path)
}

type FileFilter struct {
	*Filter
	UserIDs []int `form:"user_ids" json:"user_ids"`
}

func (p *FileFilter) Base() string {
	q := url.Values{}
	if !utils.SliceContains(p.IgnoreUrlParams, "search") && p.Search != "" {
		q.Add("q", p.Search)
	}
	if !utils.SliceContains(p.IgnoreUrlParams, "user") && len(p.UserIDs) > 0 {
		q.Add("user", strconv.Itoa(p.UserIDs[0]))
	}

	if queryString := q.Encode(); queryString != "" {
		return p.FilterBaseUrl() + "?" + q.Encode()
	}

	return p.FilterBaseUrl()
}
