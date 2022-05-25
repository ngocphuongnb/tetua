package cmd

import (
	"context"
	"testing"

	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/mock"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/test"
	"github.com/stretchr/testify/assert"
)

var mockLogger *mock.MockLogger

func init() {
	mockLogger = mock.CreateLogger()
	mock.CreateRepositories()
}

func TestFindRoleError(t *testing.T) {
	ctx := context.WithValue(context.Background(), "query_error", "true")
	defer test.RecoverPanic(t, "ByName error", "ByName error")

	createRole(auth.ROLE_ADMIN, ctx)
}

func TestCreateRoleError(t *testing.T) {
	ctx := context.WithValue(context.Background(), "create_error", "true")
	defer test.RecoverPanic(t, "Error create role", "Error create role")

	createRole(auth.ROLE_ADMIN, ctx)
}

func TestExistedRole(t *testing.T) {
	repositories.Role.Create(context.Background(), auth.ROLE_ADMIN)
	assert.Equal(t, auth.ROLE_ADMIN, createRole(auth.ROLE_ADMIN))
}

func TestCreatePermission(t *testing.T) {
	createdPerm := createPermission(auth.ROLE_ADMIN, "test.action", entities.PERM_ALL)
	assert.Equal(t, 1, createdPerm.ID)
	assert.Equal(t, 1, createdPerm.RoleID)
	assert.Equal(t, "test.action", createdPerm.Action)
	assert.Equal(t, string(entities.PERM_ALL), createdPerm.Value)

	defer test.RecoverPanic(t, "Error create permission", "Error create permission")
	ctx := context.WithValue(context.Background(), "create_error", "true")
	createPermission(auth.ROLE_ADMIN, "test.action", entities.PERM_ALL, ctx)
}

func TestSetupCreateUserFailed(t *testing.T) {
	defer test.RecoverPanic(t, "Error create user", "Error create user")
	ctx := context.WithValue(context.Background(), "create_error", "true")
	Setup("test_setup_user", "test_setup_password", ctx)
}

func TestSetupCreateUserEmptyPassword(t *testing.T) {
	defer test.RecoverPanic(t, "hash: input cannot be empty", "hash: input cannot be empty")
	Setup("test_setup_user", "")
}

func TestSetup(t *testing.T) {
	auth.Config(&server.AuthConfig{
		Action:       "test.action",
		DefaultValue: entities.PERM_NONE,
	})
	Setup("test_setup_user", "test_setup_password")

	adminRole, err := repositories.Role.ByID(context.Background(), 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, 1, adminRole.ID)

	userRole, err := repositories.Role.ByID(context.Background(), 2)
	assert.Equal(t, err, nil)
	assert.Equal(t, 2, userRole.ID)

	guestRole, err := repositories.Role.ByID(context.Background(), 3)
	assert.Equal(t, err, nil)
	assert.Equal(t, 3, guestRole.ID)

	perms, err := repositories.Permission.All(context.Background())

	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(perms))
	Setup("test_setup_user", "test_setup_password")

	lastLoggerIndex := len(mockLogger.Messages) - 1
	assert.Equal(t, mockLogger.Messages[lastLoggerIndex].Params[0], "Root user existed, skip setup")

	defer test.RecoverPanic(t, "ByUsername error", "ByUsername error")
	ctx := context.WithValue(context.Background(), "query_error", "true")
	Setup("test_setup_user", "test_setup_password", ctx)

}
