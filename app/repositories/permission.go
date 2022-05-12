package repositories

import (
	"context"

	"github.com/ngocphuongnb/tetua/app/entities"
)

type PermissionRepository interface {
	Repository[entities.Permission, entities.PermissionFilter]
	All(ctx context.Context) ([]*entities.Permission, error)
}
