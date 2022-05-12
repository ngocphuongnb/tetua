package mockrepository

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/utils"
)

type RoleRepository struct {
	*Repository[entities.Role]
}

func (m *RoleRepository) ByName(ctx context.Context, name string) (*entities.Role, error) {
	if ctx.Value("query_error") != nil {
		return nil, errors.New("ByName error")
	}

	return getEntityByField(name, m.entities, "Name", name)
}

func (m *RoleRepository) All(ctx context.Context) ([]*entities.Role, error) {
	if ctx.Value("query_error") != nil {
		return nil, errors.New("Get all roles error")
	}
	return m.entities, nil
}
func (m *RoleRepository) Find(ctx context.Context, filters ...*entities.RoleFilter) ([]*entities.Role, error) {
	if len(filters) == 0 {
		return m.entities, nil
	}
	filter := *filters[0]
	result := make([]*entities.Role, 0)
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	for index, role := range m.entities {
		if index < offset {
			continue
		}
		if index >= offset+filter.Limit {
			break
		}

		if filter.Search != "" {
			if !strings.Contains(role.Name, filter.Search) {
				continue
			}
		}

		if len(filter.ExcludeIDs) > 0 && utils.SliceContains(filter.ExcludeIDs, role.ID) {
			continue
		}

		result = append(result, role)
	}

	return result, nil
}

func (m *RoleRepository) Count(ctx context.Context, filters ...*entities.RoleFilter) (int, error) {
	if len(filters) == 0 {
		return len(m.entities), nil
	}
	filter := *filters[0]
	count := 0
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	for index, role := range m.entities {
		if index < offset {
			continue
		}
		if index >= offset+filter.Limit {
			break
		}

		if filter.Search != "" {
			if !strings.Contains(role.Name, filter.Search) {
				continue
			}
		}

		if len(filter.ExcludeIDs) > 0 && utils.SliceContains(filter.ExcludeIDs, role.ID) {
			continue
		}

		count++
	}

	return count, nil
}

func (m *RoleRepository) Paginate(ctx context.Context, filters ...*entities.RoleFilter) (*entities.Paginate[entities.Role], error) {
	comments, err := m.Find(ctx, filters...)
	if err != nil {
		return nil, err
	}

	count, err := m.Count(ctx, filters...)
	if err != nil {
		return nil, err
	}

	filter := filters[0]
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	return &entities.Paginate[entities.Role]{
		Data:        comments,
		PageSize:    filter.Limit,
		PageCurrent: filter.Page,
		Total:       int(math.Ceil(float64(count) / float64(filter.Limit))),
	}, nil
}

func (m *RoleRepository) SetPermissions(ctx context.Context, id int, permissions []*entities.PermissionValue) error {
	found := false

	for index, role := range m.entities {
		if role.ID == id {
			found = true
			permissionID := 0
			m.entities[index].Permissions = utils.SliceMap(permissions, func(item *entities.PermissionValue) *entities.Permission {
				permissionID++
				return &entities.Permission{
					ID:     permissionID,
					Action: item.Action,
					Value:  item.Value.String(),
				}
			})
		}
	}

	if !found {
		return &entities.NotFoundError{Message: "Role not found with id: " + strconv.Itoa(id)}
	}

	return nil
}
