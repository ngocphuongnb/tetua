package entrepository

import (
	"context"
	"fmt"

	e "github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent/role"
)

type RoleRepository struct {
	*BaseRepository[e.Role, ent.Role, *ent.RoleQuery, *e.RoleFilter]
}

func (p *RoleRepository) ByName(ctx context.Context, name string) (*e.Role, error) {
	r, err := p.Client.Role.Query().Where(role.NameEQ(name)).Only(ctx)
	if err != nil {
		return nil, EntError(err, fmt.Sprintf("role not found with name: %s", name))
	}
	return entRoleToRole(r), err
}

func (p *RoleRepository) SetPermissions(ctx context.Context, id int, permissions []*e.PermissionValue) error {
	var builders []*ent.PermissionCreate

	for _, permission := range permissions {
		us := p.Client.Permission.
			Create().
			SetRoleID(id).
			SetAction(permission.Action).
			SetValue(string(permission.Value))
		builders = append(builders, us)
	}

	return p.Client.Permission.
		CreateBulk(builders...).
		OnConflict().
		UpdateNewValues().
		Exec(ctx)
}

func CreateRoleRepository(client *ent.Client) *RoleRepository {
	return &RoleRepository{
		BaseRepository: &BaseRepository[e.Role, ent.Role, *ent.RoleQuery, *e.RoleFilter]{
			Name:      "role",
			Client:    client,
			ConvertFn: entRoleToRole,
			ByIDFn: func(ctx context.Context, client *ent.Client, id int) (*ent.Role, error) {
				return client.Role.Query().WithPermissions().Where(role.IDEQ(id)).Only(ctx)
			},
			DeleteByIDFn: func(ctx context.Context, client *ent.Client, id int) error {
				return client.Role.DeleteOneID(id).Exec(ctx)
			},
			CreateFn: func(ctx context.Context, client *ent.Client, data *e.Role) (*ent.Role, error) {
				return client.Role.Create().
					SetName(data.Name).
					SetDescription(data.Description).
					SetRoot(data.Root).
					Save(ctx)
			},
			UpdateFn: func(ctx context.Context, client *ent.Client, data *e.Role) (*ent.Role, error) {
				return client.Role.UpdateOneID(data.ID).
					SetName(data.Name).
					SetDescription(data.Description).
					SetRoot(data.Root).
					Save(ctx)
			},
			QueryFilterFn: func(client *ent.Client, filters ...*e.RoleFilter) *ent.RoleQuery {
				query := client.Role.Query().WithPermissions()

				if len(filters) > 0 && filters[0].Search != "" {
					query = query.Where(role.NameContainsFold(filters[0].Search))
				}

				return query
			},
			FindFn: func(ctx context.Context, query *ent.RoleQuery, filters ...*e.RoleFilter) ([]*ent.Role, error) {
				page, limit, sorts := getPaginateParams(filters...)
				return query.
					WithPermissions().
					Limit(limit).
					Offset((page - 1) * limit).
					Order(sorts...).All(ctx)
			},
		},
	}
}
