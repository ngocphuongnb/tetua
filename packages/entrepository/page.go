package entrepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/ngocphuongnb/tetua/app/entities"
	e "github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent/page"
)

type PageRepository struct {
	*BaseRepository[e.Page, ent.Page, *ent.PageQuery, *e.PageFilter]
}

func (p *PageRepository) PublishedPageBySlug(ctx context.Context, slug string) (*entities.Page, error) {
	page, err := p.Client.Page.
		Query().
		Where(page.DeletedAtIsNil()).
		Where(page.SlugEQ(slug)).
		Where(page.DraftEQ(false)).
		WithFeaturedImage().
		Only(ctx)

	if err != nil {
		return nil, EntError(err, fmt.Sprintf("page not found with slug: %s", slug))
	}

	return entPageToPage(page), nil
}

func CreatePageRepository(client *ent.Client) *PageRepository {
	return &PageRepository{
		BaseRepository: &BaseRepository[e.Page, ent.Page, *ent.PageQuery, *e.PageFilter]{
			Name:      "page",
			Client:    client,
			ConvertFn: entPageToPage,
			ByIDFn: func(ctx context.Context, client *ent.Client, id int) (*ent.Page, error) {
				return client.Page.Query().
					Where(page.DeletedAtIsNil()).
					Where(page.IDEQ(id)).
					WithFeaturedImage().
					Only(ctx)
			},
			DeleteByIDFn: func(ctx context.Context, client *ent.Client, id int) error {
				return client.Page.DeleteOneID(id).Exec(ctx)
			},
			CreateFn: func(ctx context.Context, client *ent.Client, data *e.Page) (*ent.Page, error) {
				cq := client.Page.Create().
					SetName(data.Name).
					SetSlug(data.Slug).
					SetContentHTML(data.ContentHTML).
					SetContent(data.Content).
					SetDraft(data.Draft)

				if data.FeaturedImageID != 0 {
					cq.SetFeaturedImageID(data.FeaturedImageID)
				}

				return cq.Save(ctx)
			},
			UpdateFn: func(ctx context.Context, client *ent.Client, data *e.Page) (*ent.Page, error) {
				if data.ID == 0 {
					return nil, errors.New("page id is required in order to update")
				}
				uq := client.Page.UpdateOneID(data.ID).
					SetName(data.Name).
					SetSlug(data.Slug).
					SetContentHTML(data.ContentHTML).
					SetContent(data.Content).
					SetDraft(data.Draft)

				if data.FeaturedImageID != 0 {
					uq.SetFeaturedImageID(data.FeaturedImageID)
				}

				return uq.Save(ctx)
			},
			QueryFilterFn: func(client *ent.Client, filters ...*e.PageFilter) *ent.PageQuery {
				query := client.Page.Query().Where(page.DeletedAtIsNil())
				publish := "published"

				if len(filters) > 0 {
					if filters[0].Publish != "" {
						publish = filters[0].Publish
					}

					if len(filters[0].ExcludeIDs) > 0 {
						query = query.Where(page.IDNotIn(filters[0].ExcludeIDs...))
					}
				}

				if publish != "all" {
					if publish == "published" {
						query = query.Where(page.DraftEQ(false))
					}

					if publish == "draft" {
						query = query.Where(page.DraftEQ(true))
					}
				}

				return query
			},
			FindFn: func(ctx context.Context, query *ent.PageQuery, filters ...*e.PageFilter) ([]*ent.Page, error) {
				page, limit, sorts := getPaginateParams(filters[0])
				return query.
					WithFeaturedImage().
					Limit(limit).
					Offset((page - 1) * limit).
					Order(sorts...).
					All(ctx)
			},
		},
	}
}
