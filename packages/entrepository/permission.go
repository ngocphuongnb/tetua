package entrepository

import (
	"context"

	e "github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent/permission"
)

type PermissionRepository = BaseRepository[e.Permission, ent.Permission, *ent.PermissionQuery, *e.PermissionFilter]

func CreatePermissionRepository(client *ent.Client) *PermissionRepository {
	return &PermissionRepository{
		Name:      "permission",
		Client:    client,
		ConvertFn: entPermissionToPermission,
		ByIDFn: func(ctx context.Context, client *ent.Client, id int) (*ent.Permission, error) {
			return client.Permission.Query().Where(permission.IDEQ(id)).Only(ctx)
		},
		DeleteByIDFn: func(ctx context.Context, client *ent.Client, id int) error {
			return client.Permission.DeleteOneID(id).Exec(ctx)
		},
		CreateFn: func(ctx context.Context, client *ent.Client, data *e.Permission) (*ent.Permission, error) {
			return client.Permission.Create().
				SetAction(data.Action).
				SetValue(data.Value).
				SetRoleID(data.RoleID).
				Save(ctx)
		},
		UpdateFn: func(ctx context.Context, client *ent.Client, data *e.Permission) (*ent.Permission, error) {
			return client.Permission.UpdateOneID(data.ID).
				SetAction(data.Action).
				SetValue(data.Value).
				SetRoleID(data.RoleID).
				Save(ctx)
		},
		QueryFilterFn: func(client *ent.Client, filters ...*e.PermissionFilter) *ent.PermissionQuery {
			query := client.Permission.Query()
			if len(filters) > 0 {
				if len(filters[0].RoleIDs) > 0 {
					query = query.Where(permission.RoleIDIn(filters[0].RoleIDs...))
				}

				if len(filters[0].ExcludeIDs) > 0 {
					query = query.Where(permission.IDNotIn(filters[0].ExcludeIDs...))
				}
			}

			if filters[0].Search != "" {
				query = query.Where(permission.ActionContainsFold(filters[0].Search))
			}

			return query
		},
		FindFn: func(ctx context.Context, query *ent.PermissionQuery, filters ...*e.PermissionFilter) ([]*ent.Permission, error) {
			page, limit, sorts := getPaginateParams(filters[0])
			return query.
				WithRole().
				Limit(limit).
				Offset((page - 1) * limit).
				Order(sorts...).All(ctx)
		},
	}
}
