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

type UserRepository struct {
	*Repository[entities.User]
}

var ErrorCreateIfNotExistsByProvider = false

func (m *UserRepository) CreateIfNotExistsByProvider(ctx context.Context, userData *entities.User) (*entities.User, error) {
	if ErrorCreateIfNotExistsByProvider || ctx.Value("CreateIfNotExistsByProvider") != nil {
		return nil, errors.New("CreateIfNotExistsByProvider error")
	}

	for _, user := range m.entities {
		if user.Provider == userData.Provider && user.ProviderID == userData.ProviderID {
			return user, nil
		}
	}

	return m.Create(ctx, userData)
}

func (m *UserRepository) Setting(ctx context.Context, id int, userData *entities.SettingMutation) (*entities.User, error) {
	for index, user := range m.entities {
		if user.ID == id {
			user.Username = userData.Username
			user.DisplayName = userData.DisplayName
			user.URL = userData.URL
			user.Email = userData.Email
			user.Password = userData.Password
			user.Bio = userData.Bio
			user.BioHTML = userData.BioHTML
			if user.AvatarImageID > 0 {
				user.AvatarImageID = userData.AvatarImageID
			}
			m.entities[index] = user
			return user, nil
		}
	}

	return nil, &entities.NotFoundError{Message: "User not found with id " + strconv.Itoa(id)}
}

func (m *UserRepository) ByUsername(ctx context.Context, name string) (*entities.User, error) {
	if ctx.Value("query_error") != nil {
		return nil, errors.New("ByUsername error")
	}
	return getEntityByField(name, m.entities, "Username", name)
}

func (m *UserRepository) ByProvider(ctx context.Context, name, id string) (*entities.User, error) {
	for _, user := range m.entities {
		if user.Provider == name && user.ProviderID == id {
			return user, nil
		}
	}

	return nil, &entities.NotFoundError{Message: "User not found with provider " + name + " and id " + id}
}

func (m *UserRepository) ByUsernameOrEmail(ctx context.Context, username, email string) ([]*entities.User, error) {
	result := make([]*entities.User, 0)
	for _, user := range m.entities {
		if user.Username == username || user.Email == email {
			result = append(result, user)
		}
	}

	if len(result) == 0 {
		return nil, &entities.NotFoundError{Message: "User not found with username " + username + " or email " + email}
	}

	return result, nil
}

func (m *UserRepository) All(ctx context.Context) ([]*entities.User, error) {
	if ctx.Value("query_error") != nil {
		return nil, errors.New("Get all users error")
	}
	return m.entities, nil
}
func (m *UserRepository) Find(ctx context.Context, filters ...*entities.UserFilter) ([]*entities.User, error) {
	if len(filters) == 0 {
		return m.entities, nil
	}
	filter := *filters[0]
	result := make([]*entities.User, 0)
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	for index, user := range m.entities {
		if index < offset {
			continue
		}
		if index >= offset+filter.Limit {
			break
		}

		if filter.Search != "" {
			if !strings.Contains(user.Username, filter.Search) && !strings.Contains(user.Email, filter.Search) {
				continue
			}
		}

		if len(filter.ExcludeIDs) > 0 && utils.SliceContains(filter.ExcludeIDs, user.ID) {
			continue
		}

		result = append(result, user)
	}

	return result, nil
}

func (m *UserRepository) Count(ctx context.Context, filters ...*entities.UserFilter) (int, error) {
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

	for index, user := range m.entities {
		if index < offset {
			continue
		}
		if index >= offset+filter.Limit {
			break
		}

		if filter.Search != "" {
			if !strings.Contains(user.Username, filter.Search) && !strings.Contains(user.Email, filter.Search) {
				continue
			}
		}

		if len(filter.ExcludeIDs) > 0 && utils.SliceContains(filter.ExcludeIDs, user.ID) {
			continue
		}

		count++
	}

	return count, nil
}

func (m *UserRepository) Paginate(ctx context.Context, filters ...*entities.UserFilter) (*entities.Paginate[entities.User], error) {
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
	return &entities.Paginate[entities.User]{
		Data:        topics,
		PageSize:    filter.Limit,
		PageCurrent: filter.Page,
		Total:       int(math.Ceil(float64(count) / float64(filter.Limit))),
	}, nil
}
