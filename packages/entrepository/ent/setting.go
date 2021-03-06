// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent/setting"
)

// Setting is the model entity for the Setting schema.
type Setting struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"omitempty"`
	// DeletedAt holds the value of the "deleted_at" field.
	DeletedAt time.Time `json:"omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Value holds the value of the "value" field.
	Value string `json:"value,omitempty"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Setting) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case setting.FieldID:
			values[i] = new(sql.NullInt64)
		case setting.FieldName, setting.FieldValue, setting.FieldType:
			values[i] = new(sql.NullString)
		case setting.FieldCreatedAt, setting.FieldUpdatedAt, setting.FieldDeletedAt:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Setting", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Setting fields.
func (s *Setting) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case setting.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			s.ID = int(value.Int64)
		case setting.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				s.CreatedAt = value.Time
			}
		case setting.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				s.UpdatedAt = value.Time
			}
		case setting.FieldDeletedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field deleted_at", values[i])
			} else if value.Valid {
				s.DeletedAt = value.Time
			}
		case setting.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				s.Name = value.String
			}
		case setting.FieldValue:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field value", values[i])
			} else if value.Valid {
				s.Value = value.String
			}
		case setting.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				s.Type = value.String
			}
		}
	}
	return nil
}

// Update returns a builder for updating this Setting.
// Note that you need to call Setting.Unwrap() before calling this method if this Setting
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *Setting) Update() *SettingUpdateOne {
	return (&SettingClient{config: s.config}).UpdateOne(s)
}

// Unwrap unwraps the Setting entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (s *Setting) Unwrap() *Setting {
	tx, ok := s.config.driver.(*txDriver)
	if !ok {
		panic("ent: Setting is not a transactional entity")
	}
	s.config.driver = tx.drv
	return s
}

// String implements the fmt.Stringer.
func (s *Setting) String() string {
	var builder strings.Builder
	builder.WriteString("Setting(")
	builder.WriteString(fmt.Sprintf("id=%v", s.ID))
	builder.WriteString(", created_at=")
	builder.WriteString(s.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(s.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", deleted_at=")
	builder.WriteString(s.DeletedAt.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(s.Name)
	builder.WriteString(", value=")
	builder.WriteString(s.Value)
	builder.WriteString(", type=")
	builder.WriteString(s.Type)
	builder.WriteByte(')')
	return builder.String()
}

// Settings is a parsable slice of Setting.
type Settings []*Setting

func (s Settings) config(cfg config) {
	for _i := range s {
		s[_i].config = cfg
	}
}
