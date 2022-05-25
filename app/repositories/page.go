package repositories

import (
	"context"

	"github.com/ngocphuongnb/tetua/app/entities"
)

type PageRepository interface {
	Repository[entities.Page, entities.PageFilter]
	PublishedPageBySlug(ctx context.Context, slug string) (*entities.Page, error)
}
