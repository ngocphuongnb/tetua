package mockrepository

import (
	"context"
	"math"
	"strings"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/utils"
)

type FileRepository struct {
	*Repository[entities.File]
}

func (m *FileRepository) Find(ctx context.Context, filters ...*entities.FileFilter) ([]*entities.File, error) {
	if err, ok := FakeRepoErrors["file_find"]; ok && err != nil {
		return nil, err
	}
	if len(filters) == 0 {
		return m.entities, nil
	}
	filter := *filters[0]
	result := make([]*entities.File, 0)
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
			if !strings.Contains(file.Path, filter.Search) {
				continue
			}
		}

		if len(filter.ExcludeIDs) > 0 && utils.SliceContains(filter.ExcludeIDs, file.ID) {
			continue
		}

		if len(filter.UserIDs) > 0 && !utils.SliceContains(filter.UserIDs, file.UserID) {
			continue
		}

		result = append(result, file)
	}

	return result, nil
}

func (m *FileRepository) Count(ctx context.Context, filters ...*entities.FileFilter) (int, error) {
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
			if !strings.Contains(file.Path, filter.Search) {
				continue
			}
		}

		if len(filter.ExcludeIDs) > 0 && utils.SliceContains(filter.ExcludeIDs, file.ID) {
			continue
		}

		if len(filter.UserIDs) > 0 && !utils.SliceContains(filter.UserIDs, file.UserID) {
			continue
		}

		count++
	}

	return count, nil
}

func (m *FileRepository) Paginate(ctx context.Context, filters ...*entities.FileFilter) (*entities.Paginate[entities.File], error) {
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
	return &entities.Paginate[entities.File]{
		Data:        files,
		PageSize:    filter.Limit,
		PageCurrent: filter.Page,
		Total:       int(math.Ceil(float64(count) / float64(filter.Limit))),
	}, nil
}
