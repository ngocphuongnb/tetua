package entities

import "time"

type Setting struct {
	ID        int        `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Name      string     `json:"name,omitempty"`
	Value     string     `json:"value,omitempty"`
	Type      string     `json:"type,omitempty"`
}
