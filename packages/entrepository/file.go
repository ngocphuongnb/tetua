package entrepository

import (
	"context"
	"errors"

	e "github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent/file"
)

type FileRepository = BaseRepository[e.File, ent.File, *ent.FileQuery, *e.FileFilter]

func CreateFileRepository(client *ent.Client) *FileRepository {
	return &FileRepository{
		Name:      "file",
		Client:    client,
		ConvertFn: entFileToFile,
		ByIDFn: func(ctx context.Context, client *ent.Client, id int) (*ent.File, error) {
			return client.File.Query().Where(file.IDEQ(id)).WithUser().Only(ctx)
		},
		DeleteByIDFn: func(ctx context.Context, client *ent.Client, id int) error {
			return client.File.DeleteOneID(id).Exec(ctx)
		},
		CreateFn: func(ctx context.Context, client *ent.Client, data *e.File) (*ent.File, error) {
			if data.UserID == 0 {
				return nil, errors.New("user_id is required")
			}
			return client.File.Create().
				SetDisk(data.Disk).
				SetPath(data.Path).
				SetSize(data.Size).
				SetType(data.Type).
				SetUserID(data.UserID).
				Save(ctx)
		},
		UpdateFn: func(ctx context.Context, client *ent.Client, data *e.File) (*ent.File, error) {
			if data.ID == 0 {
				return nil, errors.New("ID is required")
			}
			return client.File.UpdateOneID(data.ID).
				SetDisk(data.Disk).
				SetPath(data.Path).
				SetSize(data.Size).
				SetType(data.Type).
				Save(ctx)
		},
		QueryFilterFn: func(client *ent.Client, filters ...*e.FileFilter) *ent.FileQuery {
			query := client.File.Query()
			if len(filters) > 0 {
				if filters[0].Search != "" {
					query = query.Where(file.PathContainsFold(filters[0].Search))
				}
				if len(filters[0].UserIDs) > 0 {
					query = query.Where(file.UserIDIn(filters[0].UserIDs...))
				}

				if len(filters[0].ExcludeIDs) > 0 {
					query = query.Where(file.IDNotIn(filters[0].ExcludeIDs...))
				}
			}

			return query
		},
		FindFn: func(ctx context.Context, query *ent.FileQuery, filters ...*e.FileFilter) ([]*ent.File, error) {
			page, limit, sorts := getPaginateParams(filters[0])
			return query.
				Limit(limit).
				Offset((page - 1) * limit).
				Order(sorts...).All(ctx)
		},
	}
}
