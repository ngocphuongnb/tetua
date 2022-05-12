package cache

import (
	"context"
	"sync"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/utils"
)

var RolesPermissions = []*entities.RolePermissions{}
var Topics = []*entities.Topic{}
var Roles []*entities.Role

func All() error {
	var err1 error
	var err2 error
	var wg sync.WaitGroup

	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err1 = CacheTopics()
	}(&wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err2 = CachePermissions()
	}(&wg)
	wg.Wait()

	return utils.FirstError(err1, err2)
}

func CacheTopics(ctxs ...context.Context) (err error) {
	ctxs = append(ctxs, context.Background())
	Topics, err = repositories.Topic.All(ctxs[0])
	return
}

func CachePermissions(ctxs ...context.Context) (err error) {
	ctxs = append(ctxs, context.Background())
	Roles, err = repositories.Role.All(ctxs[0])

	if err != nil {
		return err
	}

	for _, role := range Roles {
		permissions := []*entities.PermissionValue{}

		for _, permission := range role.Permissions {
			permissions = append(permissions, &entities.PermissionValue{
				Action: permission.Action,
				Value:  entities.GetPermTypeValue(permission.Value),
			})
		}

		RolesPermissions = append(RolesPermissions, &entities.RolePermissions{
			RoleID:      role.ID,
			Permissions: permissions,
		})
	}

	return nil
}
