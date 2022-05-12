package repositories

import (
	"context"

	"github.com/ngocphuongnb/tetua/app/entities"
)

type CommentRepository interface {
	Repository[entities.Comment, entities.CommentFilter]
	FindWithPost(ctx context.Context, filters ...*entities.CommentFilter) ([]*entities.Comment, error)
	PaginateWithPost(ctx context.Context, filters ...*entities.CommentFilter) (*entities.Paginate[entities.Comment], error)
}
