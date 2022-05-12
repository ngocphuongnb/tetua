package mockrepository

import (
	"context"
	"errors"
	"math"
	"strings"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/utils"
)

type TopicRepository struct {
	*Repository[entities.Topic]
}

func (m *TopicRepository) ByName(ctx context.Context, name string) (*entities.Topic, error) {
	return getEntityByField(name, m.entities, "Name", name)
}

func (m *TopicRepository) All(ctx context.Context) ([]*entities.Topic, error) {
	if ctx.Value("query_error") != nil {
		return nil, errors.New("Get all topics error")
	}
	return m.entities, nil
}
func (m *TopicRepository) Find(ctx context.Context, filters ...*entities.TopicFilter) ([]*entities.Topic, error) {
	if len(filters) == 0 {
		return m.entities, nil
	}
	filter := *filters[0]
	result := make([]*entities.Topic, 0)
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	for index, topic := range m.entities {
		if index < offset {
			continue
		}
		if index >= offset+filter.Limit {
			break
		}

		if filter.Search != "" {
			if !strings.Contains(topic.Name, filter.Search) {
				continue
			}
		}

		if len(filter.ExcludeIDs) > 0 && utils.SliceContains(filter.ExcludeIDs, topic.ID) {
			continue
		}

		result = append(result, topic)
	}

	return result, nil
}

func (m *TopicRepository) Count(ctx context.Context, filters ...*entities.TopicFilter) (int, error) {
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

	for index, topic := range m.entities {
		if index < offset {
			continue
		}
		if index >= offset+filter.Limit {
			break
		}

		if filter.Search != "" {
			if !strings.Contains(topic.Name, filter.Search) {
				continue
			}
		}

		if len(filter.ExcludeIDs) > 0 && utils.SliceContains(filter.ExcludeIDs, topic.ID) {
			continue
		}

		count++
	}

	return count, nil
}

func (m *TopicRepository) Paginate(ctx context.Context, filters ...*entities.TopicFilter) (*entities.Paginate[entities.Topic], error) {
	topics, err := m.Find(ctx, filters...)
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
	return &entities.Paginate[entities.Topic]{
		Data:        topics,
		PageSize:    filter.Limit,
		PageCurrent: filter.Page,
		Total:       int(math.Ceil(float64(count) / float64(filter.Limit))),
	}, nil
}
