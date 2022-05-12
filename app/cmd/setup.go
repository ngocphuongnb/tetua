package cmd

import (
	"context"

	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/logger"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/utils"
)

func createRole(role *entities.Role, ctxs ...context.Context) *entities.Role {
	ctxs = append(ctxs, context.Background())
	existedRole, err := repositories.Role.ByName(ctxs[0], role.Name)

	if err != nil && !entities.IsNotFound(err) {
		logger.Error(err)
		panic(err)
	}

	if existedRole != nil {
		return existedRole
	}

	role, err = repositories.Role.Create(ctxs[0], role)

	if err != nil {
		logger.Error(err)
		panic(err)
	}

	return role
}

func createPermission(role *entities.Role, action string, value entities.PermType, ctxs ...context.Context) *entities.Permission {
	ctxs = append(ctxs, context.Background())
	permission, err := repositories.Permission.Create(ctxs[0], &entities.Permission{
		Action: action,
		Value:  string(value),
		RoleID: role.ID,
	})

	if err != nil {
		logger.Error(err)
		panic(err)
	}

	return permission
}

func Setup(username, password string, ctxs ...context.Context) error {
	ctxs = append(ctxs, context.Background())
	adminRole := createRole(auth.ROLE_ADMIN)
	userRole := createRole(auth.ROLE_USER)
	guestRole := createRole(auth.ROLE_GUEST)

	for _, authConfig := range auth.ActionConfigs {
		createPermission(userRole, authConfig.Action, authConfig.DefaultValue)
	}
	for _, authConfig := range auth.ActionConfigs {
		createPermission(guestRole, authConfig.Action, authConfig.DefaultValue)
	}

	rootUser, err := repositories.User.ByUsername(ctxs[0], username)

	if err != nil && !entities.IsNotFound(err) {
		logger.Error(err)
		panic(err)
	}

	if rootUser != nil {
		logger.Info("Root user existed, skip setup")
		return nil
	}

	rootPassword, err := utils.GenerateHash(password)

	if err != nil {
		logger.Error(err)
		panic(err)
	}

	if _, err = repositories.User.Create(ctxs[0], &entities.User{
		Username: username,
		Password: rootPassword,
		Provider: "local",
		RoleIDs:  []int{adminRole.ID},
		Active:   true,
	}); err != nil {
		logger.Error(err)
		panic(err)
	}

	logger.Info("Setup root user successfully")

	return nil
}
