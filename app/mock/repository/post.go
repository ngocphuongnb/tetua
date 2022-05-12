package mockrepository

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/utils"
)

type PostRepository struct {
	*Repository[entities.Post]
}

func (m *PostRepository) Find(ctx context.Context, filters ...*entities.PostFilter) ([]*entities.Post, error) {
	if err, ok := FakeRepoErrors[m.Name+"_find"]; ok && err != nil {
		return nil, err
	}
	if len(filters) == 0 {
		return m.entities, nil
	}
	filter := *filters[0]
	result := make([]*entities.Post, 0)
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	for index, post := range m.entities {
		if index < offset {
			continue
		}
		if index >= offset+filter.Limit {
			break
		}

		if filter.Search != "" {
			if !strings.Contains(post.Name, filter.Search) && !strings.Contains(post.Content, filter.Search) {
				continue
			}
		}

		if len(filter.ExcludeIDs) > 0 && utils.SliceContains(filter.ExcludeIDs, post.ID) {
			continue
		}

		if len(filter.UserIDs) > 0 && !utils.SliceContains(filter.UserIDs, post.UserID) {
			continue
		}

		if len(filter.TopicIDs) > 0 && len(utils.SliceOverlap(filter.TopicIDs, post.TopicIDs)) == 0 {
			continue
		}

		if filter.Publish == "published" && post.Draft {
			continue
		}
		if filter.Publish == "draft" && !post.Draft {
			continue
		}

		user, err := repositories.User.ByID(context.Background(), post.UserID)

		if err != nil || user == nil {
			user = &entities.User{
				ID:       post.UserID,
				Username: fmt.Sprintf("testuser%d", post.UserID),
				Provider: "local",
			}
		}
		post.User = user

		result = append(result, post)

	}

	return result, nil
}

func (m *PostRepository) Count(ctx context.Context, filters ...*entities.PostFilter) (int, error) {
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

	for index, post := range m.entities {
		if index < offset {
			continue
		}
		if index >= offset+filter.Limit {
			break
		}

		if filter.Search != "" {
			if !strings.Contains(post.Name, filter.Search) && !strings.Contains(post.Content, filter.Search) {
				continue
			}
		}

		if len(filter.ExcludeIDs) > 0 && utils.SliceContains(filter.ExcludeIDs, post.ID) {
			continue
		}

		if len(filter.UserIDs) > 0 && !utils.SliceContains(filter.UserIDs, post.UserID) {
			continue
		}

		if len(filter.TopicIDs) > 0 && len(utils.SliceOverlap(filter.TopicIDs, post.TopicIDs)) == 0 {
			continue
		}

		if filter.Publish == "published" && post.Draft {
			continue
		}
		if filter.Publish == "draft" && !post.Draft {
			continue
		}

		count++
	}

	return count, nil
}

func (m *PostRepository) Paginate(ctx context.Context, filters ...*entities.PostFilter) (*entities.Paginate[entities.Post], error) {
	if err, ok := FakeRepoErrors[m.Name+"_paginate"]; ok && err != nil {
		return nil, err
	}

	posts, err := m.Find(ctx, filters...)
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
	return &entities.Paginate[entities.Post]{
		Data:        posts,
		PageSize:    filter.Limit,
		PageCurrent: filter.Page,
		Total:       int(math.Ceil(float64(count) / float64(filter.Limit))),
	}, nil
}

func (m *PostRepository) IncreaseViewCount(ctx context.Context, id int, views int64) error {
	foundPost := false
	m.mu.Lock()
	defer m.mu.Unlock()

	for index, post := range m.entities {
		if post.ID == id {
			foundPost = true
			post.ViewCount += views
			m.entities[index] = post
		}
	}

	if !foundPost {
		return &entities.NotFoundError{Message: "post not found with id: " + strconv.Itoa(id)}
	}

	return nil
}

func (m *PostRepository) PublishedPostByID(ctx context.Context, id int) (post *entities.Post, err error) {
	if post, err = m.ByID(ctx, id); err != nil {
		return nil, err
	}

	if post.Draft || !post.Approved {
		return nil, &entities.NotFoundError{Message: "post not found with id: " + strconv.Itoa(id)}
	}

	return post, nil
}

func (m *PostRepository) Approve(ctx context.Context, id int) error {
	post, err := m.ByID(ctx, id)
	if err != nil {
		return err
	}

	post.Approved = true
	_, err = m.Update(ctx, post)

	return err
}
