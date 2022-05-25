package entrepository

import (
	"context"
	"fmt"
	"sync"

	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent"
)

type EntityType interface {
	ent.Comment | ent.File | ent.Permission | ent.Post | ent.Page | ent.Role | ent.Setting | ent.Topic | ent.User
}

type QueryFilter interface {
	GetSearch() string
	GetPage() int
	GetLimit() int
	GetSorts() []*entities.Sort
	GetIgnoreUrlParams() []string
	GetExcludeIDs() []int
	Base() string
}

type EntityQuery[EE EntityType] interface {
	*ent.CommentQuery | *ent.FileQuery | *ent.PermissionQuery | *ent.PostQuery | *ent.PageQuery | *ent.RoleQuery | *ent.SettingQuery | *ent.TopicQuery | *ent.UserQuery
	Count(context.Context) (int, error)
	All(context.Context) ([]*EE, error)
}

type BaseRepository[E entities.Entity, EE EntityType, EQ EntityQuery[EE], QF QueryFilter] struct {
	Name          string
	Client        *ent.Client
	ConvertFn     func(entEntity *EE) *E
	ByIDFn        func(ctx context.Context, client *ent.Client, id int) (*EE, error)
	DeleteByIDFn  func(ctx context.Context, client *ent.Client, id int) error
	CreateFn      func(ctx context.Context, client *ent.Client, data *E) (*EE, error)
	UpdateFn      func(ctx context.Context, client *ent.Client, data *E) (*EE, error)
	FindFn        func(ctx context.Context, query EQ, filters ...QF) ([]*EE, error)
	QueryFilterFn func(client *ent.Client, filters ...QF) EQ
}

func (b *BaseRepository[E, EE, EQ, QF]) ByID(ctx context.Context, id int) (*E, error) {
	entity, err := b.ByIDFn(ctx, b.Client, id)

	if err != nil {
		return nil, EntError(err, fmt.Sprintf("%s not found with id: %d", b.Name, id))
	}

	return b.ConvertFn(entity), nil
}

func (b *BaseRepository[E, EE, EQ, QF]) DeleteByID(ctx context.Context, id int) error {
	return b.DeleteByIDFn(ctx, b.Client, id)
}

func (b *BaseRepository[E, EE, EQ, QF]) Create(ctx context.Context, data *E) (*E, error) {
	entity, err := b.CreateFn(ctx, b.Client, data)
	if err != nil {
		return nil, err
	}

	return b.ConvertFn(entity), nil
}

func (b *BaseRepository[E, EE, EQ, QF]) Update(ctx context.Context, data *E) (*E, error) {
	entity, err := b.UpdateFn(ctx, b.Client, data)
	if err != nil {
		return nil, err
	}

	return b.ConvertFn(entity), nil
}

func (b *BaseRepository[E, EE, EQ, QF]) Count(ctx context.Context, filters ...QF) (int, error) {
	return b.QueryFilterFn(b.Client, filters...).Count(ctx)
}

func getPaginateParams[F QueryFilter](filters ...F) (int, int, []ent.OrderFunc) {
	page := 1
	limit := 10
	sorts := []ent.OrderFunc{ent.Desc("id")}

	if len(filters) > 0 {
		if filters[0].GetLimit() > 0 {
			limit = filters[0].GetLimit()
		}

		if filters[0].GetPage() > 0 {
			page = filters[0].GetPage()
		}

		if len(filters[0].GetSorts()) > 0 {
			sorts = getSortFNs(filters[0].GetSorts())
		}
	}

	return page, limit, sorts
}

func (b *BaseRepository[E, EE, EQ, QF]) All(ctx context.Context) ([]*E, error) {
	query := b.QueryFilterFn(b.Client)
	if items, err := query.All(ctx); err != nil {
		return nil, err
	} else {
		return utils.SliceMap(items, b.ConvertFn), nil
	}
}

func (b *BaseRepository[E, EE, EQ, QF]) Find(ctx context.Context, filters ...QF) ([]*E, error) {
	query := b.QueryFilterFn(b.Client, filters...)
	if items, err := b.FindFn(ctx, query, filters...); err != nil {
		return nil, err
	} else {
		return utils.SliceMap(items, b.ConvertFn), nil
	}
}

func (b *BaseRepository[E, EE, EQ, QF]) Paginate(ctx context.Context, filters ...QF) (*entities.Paginate[E], error) {
	var err1 error
	var err2 error
	var wg sync.WaitGroup
	total := 0
	base := ""
	items := make([]*EE, 0)
	page, limit, _ := getPaginateParams(filters[0])

	wg.Add(2)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		total, err1 = b.QueryFilterFn(b.Client, filters...).Count(ctx)
	}(&wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		items, err2 = b.FindFn(ctx, b.QueryFilterFn(b.Client, filters...), filters...)
	}(&wg)
	wg.Wait()

	if err := utils.FirstError(err1, err2); err != nil {
		return nil, err
	}

	if len(filters) > 0 {
		base = filters[0].Base()
	}

	return &entities.Paginate[E]{
		Data:        utils.SliceMap(items, b.ConvertFn),
		BaseUrl:     base,
		Total:       total,
		PageSize:    limit,
		PageCurrent: page,
	}, nil
}

func EntError(err error, msg string) error {
	if ent.IsNotFound(err) {
		return &entities.NotFoundError{Message: msg}
	}

	return err
}
