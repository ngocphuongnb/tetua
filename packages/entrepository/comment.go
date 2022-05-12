package entrepository

import (
	"context"
	"errors"
	"sync"

	e "github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent/comment"
)

type CommentRepository struct {
	*BaseRepository[e.Comment, ent.Comment, *ent.CommentQuery, *e.CommentFilter]
}

func (c *CommentRepository) FindWithPost(ctx context.Context, filters ...*e.CommentFilter) ([]*e.Comment, error) {
	query := c.QueryFilterFn(c.Client, filters...).WithPost()
	if items, err := c.FindFn(ctx, query, filters...); err != nil {
		return nil, err
	} else {
		return utils.SliceMap(items, c.ConvertFn), nil
	}
}

func (c *CommentRepository) PaginateWithPost(ctx context.Context, filters ...*e.CommentFilter) (*e.Paginate[e.Comment], error) {
	var err1 error
	var err2 error
	var wg sync.WaitGroup
	total := 0
	base := ""
	items := make([]*ent.Comment, 0)
	page, limit, _ := getPaginateParams(filters[0])

	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		total, err1 = c.QueryFilterFn(c.Client, filters...).WithPost().Count(ctx)
	}(&wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		items, err2 = c.FindFn(ctx, c.QueryFilterFn(c.Client, filters...).WithPost(), filters...)
	}(&wg)
	wg.Wait()

	if err := utils.FirstError(err1, err2); err != nil {
		return nil, err
	}

	if len(filters) > 0 {
		base = filters[0].Base()
	}

	return &e.Paginate[e.Comment]{
		Data:        utils.SliceMap(items, c.ConvertFn),
		BaseUrl:     base,
		Total:       total,
		PageSize:    limit,
		PageCurrent: page,
	}, nil
}

func CreateCommentRepository(client *ent.Client) *CommentRepository {
	return &CommentRepository{
		BaseRepository: &BaseRepository[e.Comment, ent.Comment, *ent.CommentQuery, *e.CommentFilter]{
			Name:      "comment",
			Client:    client,
			ConvertFn: entCommentToComment,
			ByIDFn: func(ctx context.Context, client *ent.Client, id int) (*ent.Comment, error) {
				return client.Comment.Query().Where(comment.IDEQ(id)).WithUser().Only(ctx)
			},
			DeleteByIDFn: func(ctx context.Context, client *ent.Client, id int) error {
				return client.Comment.DeleteOneID(id).Exec(ctx)
			},
			CreateFn: func(ctx context.Context, client *ent.Client, data *e.Comment) (comment *ent.Comment, err error) {
				if data.UserID == 0 {
					return nil, errors.New("user_id is required")
				}
				cc := client.Comment.Create().
					SetContent(data.Content).
					SetContentHTML(data.ContentHTML).
					SetPostID(data.PostID).
					SetUserID(data.UserID)

				if data.ParentID > 0 {
					cc.SetParentID(data.ParentID)
				}

				if comment, err = cc.Save(ctx); err != nil {
					return nil, err
				}

				if err = client.Post.UpdateOneID(data.PostID).AddCommentCount(1).Exec(ctx); err != nil {
					return nil, err
				}

				return
			},
			UpdateFn: func(ctx context.Context, client *ent.Client, data *e.Comment) (comment *ent.Comment, err error) {
				return client.Comment.
					UpdateOneID(data.ID).
					SetContent(data.Content).
					SetContentHTML(data.ContentHTML).
					Save(ctx)
			},
			QueryFilterFn: func(client *ent.Client, filters ...*e.CommentFilter) *ent.CommentQuery {
				query := client.Comment.Query()
				if len(filters) > 0 {
					if len(filters[0].UserIDs) > 0 {
						query = query.Where(comment.UserIDIn(filters[0].UserIDs...))
					}

					if len(filters[0].PostIDs) > 0 {
						query = query.Where(comment.PostIDIn(filters[0].PostIDs...))
					}

					if filters[0].Search != "" {
						query = query.Where(comment.ContentContainsFold(filters[0].Search))
					}
				}
				return query
			},
			FindFn: func(ctx context.Context, query *ent.CommentQuery, filters ...*e.CommentFilter) ([]*ent.Comment, error) {
				page, limit, sorts := getPaginateParams(filters[0])
				return query.
					WithUser(func(uq *ent.UserQuery) {
						uq.WithAvatarImage()
					}).
					Limit(limit).
					Offset((page - 1) * limit).
					Order(sorts...).
					All(ctx)
			},
		},
	}
}
