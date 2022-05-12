package mockrepository

import (
	"context"
	"math"
	"strings"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/utils"
)

type CommentRepository struct {
	*Repository[entities.Comment]
}

func (m *CommentRepository) Find(ctx context.Context, filters ...*entities.CommentFilter) ([]*entities.Comment, error) {
	if len(filters) == 0 {
		return m.entities, nil
	}
	filter := *filters[0]
	result := make([]*entities.Comment, 0)
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	for index, comment := range m.entities {
		if index < offset {
			continue
		}
		if index >= offset+filter.Limit {
			break
		}

		if filter.Search != "" {
			if !strings.Contains(comment.Content, filter.Search) && !strings.Contains(comment.Content, filter.Search) {
				continue
			}
		}

		if len(filter.ExcludeIDs) > 0 && utils.SliceContains(filter.ExcludeIDs, comment.ID) {
			continue
		}

		if len(filter.UserIDs) > 0 && !utils.SliceContains(filter.UserIDs, comment.UserID) {
			continue
		}

		if len(filter.PostIDs) > 0 && !utils.SliceContains(filter.PostIDs, comment.UserID) {
			continue
		}

		if len(filter.ParentIDs) > 0 && !utils.SliceContains(filter.ParentIDs, comment.UserID) {
			continue
		}

		result = append(result, comment)
	}

	return result, nil
}

func (m *CommentRepository) FindWithPost(ctx context.Context, filters ...*entities.CommentFilter) ([]*entities.Comment, error) {
	var err error
	posts, err := m.Find(ctx, filters...)

	if err != nil {
		return nil, err
	}

	return utils.SliceMap(posts, func(comment *entities.Comment) *entities.Comment {
		comment.Post, err = repositories.Post.ByID(ctx, comment.PostID)
		return comment
	}), err
}

func (m *CommentRepository) Count(ctx context.Context, filters ...*entities.CommentFilter) (int, error) {
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

	for index, comment := range m.entities {
		if index < offset {
			continue
		}
		if index >= offset+filter.Limit {
			break
		}

		if filter.Search != "" {
			if !strings.Contains(comment.Content, filter.Search) && !strings.Contains(comment.Content, filter.Search) {
				continue
			}
		}

		if len(filter.ExcludeIDs) > 0 && utils.SliceContains(filter.ExcludeIDs, comment.ID) {
			continue
		}

		if len(filter.UserIDs) > 0 && !utils.SliceContains(filter.UserIDs, comment.UserID) {
			continue
		}

		if len(filter.PostIDs) > 0 && !utils.SliceContains(filter.PostIDs, comment.UserID) {
			continue
		}

		if len(filter.ParentIDs) > 0 && !utils.SliceContains(filter.ParentIDs, comment.UserID) {
			continue
		}

		count++
	}

	return count, nil
}

func (m *CommentRepository) Paginate(ctx context.Context, filters ...*entities.CommentFilter) (*entities.Paginate[entities.Comment], error) {
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
	return &entities.Paginate[entities.Comment]{
		Data:        comments,
		PageSize:    filter.Limit,
		PageCurrent: filter.Page,
		Total:       int(math.Ceil(float64(count) / float64(filter.Limit))),
	}, nil
}

func (m *CommentRepository) PaginateWithPost(ctx context.Context, filters ...*entities.CommentFilter) (*entities.Paginate[entities.Comment], error) {
	comments, err := m.FindWithPost(ctx, filters...)
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
	return &entities.Paginate[entities.Comment]{
		Data:        comments,
		PageSize:    filter.Limit,
		PageCurrent: filter.Page,
		Total:       int(math.Ceil(float64(count) / float64(filter.Limit))),
	}, nil
}
