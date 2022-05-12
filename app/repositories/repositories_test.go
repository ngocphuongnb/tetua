package repositories_test

import (
	"testing"

	"github.com/ngocphuongnb/tetua/app/mock"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/stretchr/testify/assert"
)

func TestComment(t *testing.T) {
	repos := mock.Repositories()
	repositories.New(repos)

	assert.Equal(t, repos.File, repositories.File)
	assert.Equal(t, repos.Post, repositories.Post)
	assert.Equal(t, repos.Comment, repositories.Comment)
	assert.Equal(t, repos.Role, repositories.Role)
	assert.Equal(t, repos.Topic, repositories.Topic)
	assert.Equal(t, repos.User, repositories.User)
	assert.Equal(t, repos.Permission, repositories.Permission)
}
