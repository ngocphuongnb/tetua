package mockrepository

import (
	"context"
	"math"
	"strings"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/utils"
)

type PermissionRepository struct {
	*Repository[entities.Permission]
}

func (p *PermissionRepository) All(ctx context.Context) ([]*entities.Permission, error) {
	return p.entities, nil
}

func (m *PermissionRepository) Find(ctx context.Context, filters ...*entities.PermissionFilter) ([]*entities.Permission, error) {
	if len(filters) == 0 {
		return m.entities, nil
	}
	filter := *filters[0]
	result := make([]*entities.Permission, 0)
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	for index, perm := range m.entities {
		if index < offset {
			continue
		}
		if index >= offset+filter.Limit {
			break
		}

		if filter.Search != "" {
			if !strings.Contains(perm.Action, filter.Search) {
				continue
			}
		}

		if len(filter.ExcludeIDs) > 0 && utils.SliceContains(filter.ExcludeIDs, perm.ID) {
			continue
		}

		if len(filter.RoleIDs) > 0 && !utils.SliceContains(filter.RoleIDs, perm.RoleID) {
			continue
		}

		result = append(result, perm)
	}

	return result, nil
}

func (m *PermissionRepository) Count(ctx context.Context, filters ...*entities.PermissionFilter) (int, error) {
	if len(filters) == 0 {
		return len(m.entities), nil
	}
	filter := filters[0]
	count := 0
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	for index, file := range m.entities {
		if index < offset {
			continue
		}
		if index >= offset+filter.Limit {
			break
		}

		if filter.Search != "" {
			if !strings.Contains(file.Action, filter.Search) {
				continue
			}
		}

		if len(filter.ExcludeIDs) > 0 && utils.SliceContains(filter.ExcludeIDs, file.ID) {
			continue
		}

		if len(filter.RoleIDs) > 0 && !utils.SliceContains(filter.RoleIDs, file.RoleID) {
			continue
		}

		count++
	}

	return count, nil
}

func (m *PermissionRepository) Paginate(ctx context.Context, filters ...*entities.PermissionFilter) (*entities.Paginate[entities.Permission], error) {
	files, err := m.Find(ctx, filters...)
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
	return &entities.Paginate[entities.Permission]{
		Data:        files,
		PageSize:    filter.Limit,
		PageCurrent: filter.Page,
		Total:       int(math.Ceil(float64(count) / float64(filter.Limit))),
	}, nil
}
