package cache_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/cache"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/mock"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/stretchr/testify/assert"
)

func init() {
	mock.CreateRepositories()
}

func TestCache(t *testing.T) {
	repositories.Topic.Create(context.Background(), &entities.Topic{
		ID:   1,
		Name: "Topic 1",
	})
	repositories.Role.Create(context.Background(), auth.ROLE_ADMIN)
	repositories.Role.Create(context.Background(), auth.ROLE_USER)
	repositories.Role.Create(context.Background(), auth.ROLE_GUEST)

	err := repositories.Role.SetPermissions(context.Background(), auth.ROLE_USER.ID, []*entities.PermissionValue{{
		Action: "post.view",
		Value:  entities.PERM_ALL,
	}})

	assert.Equal(t, nil, err)
	assert.Equal(t, nil, cache.All())
	assert.Equal(t, 1, len(cache.Topics))
	assert.Equal(t, 3, len(cache.Roles))
	assert.Equal(t, 2, cache.RolesPermissions[1].RoleID)
	assert.Equal(t, []*entities.PermissionValue{{
		Action: "post.view",
		Value:  entities.PERM_ALL,
	}}, cache.RolesPermissions[1].Permissions)
}

func TestCacheError(t *testing.T) {
	ctx := context.WithValue(context.Background(), "query_error", true)

	assert.Equal(t, errors.New("Get all topics error"), cache.CacheTopics(ctx))
	assert.Equal(t, errors.New("Get all roles error"), cache.CachePermissions(ctx))
}
