package repositories

import (
	"context"

	"github.com/ngocphuongnb/tetua/app/entities"
)

type PostRepository interface {
	Repository[entities.Post, entities.PostFilter]
	Approve(ctx context.Context, id int) error
	PublishedPostByID(ctx context.Context, id int) (*entities.Post, error)
	IncreaseViewCount(ctx context.Context, id int, views int64) error
}
