package repositories

import (
	"context"

	"github.com/ngocphuongnb/tetua/app/entities"
)

type RoleRepository interface {
	Repository[entities.Role, entities.RoleFilter]
	All(ctx context.Context) ([]*entities.Role, error)
	ByName(ctx context.Context, name string) (*entities.Role, error)
	SetPermissions(ctx context.Context, id int, permissions []*entities.PermissionValue) error
}
