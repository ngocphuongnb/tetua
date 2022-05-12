package entrepository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ngocphuongnb/tetua/app/entities"
	e "github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent/post"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent/topic"
)

type PostRepository struct {
	*BaseRepository[e.Post, ent.Post, *ent.PostQuery, *e.PostFilter]
}

func (p *PostRepository) IncreaseViewCount(ctx context.Context, id int, views int64) (err error) {
	_, err = p.Client.Post.UpdateOneID(id).AddViewCount(views).Save(ctx)
	return
}

func (p *PostRepository) Approve(ctx context.Context, id int) (err error) {
	_, err = p.Client.Post.UpdateOneID(id).SetApproved(true).Save(ctx)
	return
}

func (p *PostRepository) PublishedPostByID(ctx context.Context, id int) (*entities.Post, error) {
	post, err := p.Client.Post.
		Query().
		Where(post.DeletedAtIsNil()).
		Where(post.IDEQ(id)).
		Where(post.DraftEQ(false)).
		Where(post.Approved(true)).
		WithUser(func(uq *ent.UserQuery) {
			uq.WithAvatarImage()
		}).
		WithTopics().
		WithFeaturedImage().
		Only(ctx)

	if err != nil {
		return nil, EntError(err, fmt.Sprintf("post not found with id: %d", id))
	}

	return entPostToPost(post), nil
}

func CreatePostRepository(client *ent.Client) *PostRepository {
	return &PostRepository{
		BaseRepository: &BaseRepository[e.Post, ent.Post, *ent.PostQuery, *e.PostFilter]{
			Name:      "post",
			Client:    client,
			ConvertFn: entPostToPost,
			ByIDFn: func(ctx context.Context, client *ent.Client, id int) (*ent.Post, error) {
				return client.Post.Query().
					Where(post.DeletedAtIsNil()).
					Where(post.IDEQ(id)).
					WithUser(func(uq *ent.UserQuery) {
						uq.WithAvatarImage()
					}).
					WithTopics().
					WithFeaturedImage().
					Only(ctx)
			},
			DeleteByIDFn: func(ctx context.Context, client *ent.Client, id int) error {
				return client.Post.UpdateOneID(id).SetDeletedAt(time.Now()).Exec(ctx)
			},
			CreateFn: func(ctx context.Context, client *ent.Client, data *e.Post) (*ent.Post, error) {
				if data.UserID == 0 {
					return nil, errors.New("user_id is required")
				}
				cq := client.Post.Create().
					SetName(data.Name).
					SetDescription(data.Description).
					SetSlug(data.Slug).
					SetContentHTML(data.ContentHTML).
					SetContent(data.Content).
					SetDraft(data.Draft).
					SetApproved(data.Approved).
					SetUserID(data.UserID).
					AddTopicIDs(data.TopicIDs...)

				if data.FeaturedImageID != 0 {
					cq.SetFeaturedImageID(data.FeaturedImageID)
				}

				return cq.Save(ctx)
			},
			UpdateFn: func(ctx context.Context, client *ent.Client, data *e.Post) (*ent.Post, error) {
				if data.ID == 0 {
					return nil, errors.New("post id is required in order to update")
				}
				uq := client.Post.UpdateOneID(data.ID).
					SetName(data.Name).
					SetDescription(data.Description).
					SetSlug(data.Slug).
					SetContentHTML(data.ContentHTML).
					SetContent(data.Content).
					SetDraft(data.Draft).
					SetApproved(data.Approved)

				if len(data.TopicIDs) > 0 {
					oldPostEnt, err := client.Post.Query().
						WithTopics().
						Where(post.DeletedAtIsNil()).
						Where(post.IDEQ(data.ID)).Only(ctx)

					if err != nil {
						return nil, err
					}

					oldPost := entPostToPost(oldPostEnt)
					oldTopicIDs := utils.SliceMap(oldPost.Topics, func(t *entities.Topic) int {
						return t.ID
					})
					uq.RemoveTopicIDs(oldTopicIDs...).AddTopicIDs(data.TopicIDs...)
				}

				if data.FeaturedImageID != 0 {
					uq.SetFeaturedImageID(data.FeaturedImageID)
				}

				return uq.Save(ctx)
			},
			QueryFilterFn: func(client *ent.Client, filters ...*e.PostFilter) *ent.PostQuery {
				query := client.Post.Query()
				publish := "published"
				approved := "approved"

				if len(filters) > 0 {
					if filters[0].Publish != "" {
						publish = filters[0].Publish
					}
					if filters[0].Approve != "" {
						approved = filters[0].Approve
					}

					if len(filters[0].UserIDs) > 0 {
						query = query.Where(post.UserIDIn(filters[0].UserIDs...))
					}

					if len(filters[0].TopicIDs) > 0 {
						query = query.Where(post.HasTopicsWith(topic.IDIn(filters[0].TopicIDs...)))
					}

					if len(filters[0].ExcludeIDs) > 0 {
						query = query.Where(post.IDNotIn(filters[0].ExcludeIDs...))
					}
				}

				if publish != "all" {
					if publish == "published" {
						query = query.Where(post.DraftEQ(false))
					}

					if publish == "draft" {
						query = query.Where(post.DraftEQ(true))
					}
				}

				if approved != "all" {
					if approved == "approved" {
						query = query.Where(post.Approved(true))
					}

					if approved == "pending" {
						query = query.Where(post.Approved(false))
					}
				}

				return query
			},
			FindFn: func(ctx context.Context, query *ent.PostQuery, filters ...*e.PostFilter) ([]*ent.Post, error) {
				page, limit, sorts := getPaginateParams(filters[0])
				return query.
					WithUser(func(uq *ent.UserQuery) {
						uq.WithAvatarImage()
					}).
					WithTopics().
					WithFeaturedImage().
					Limit(limit).
					Offset((page - 1) * limit).
					Order(sorts...).
					All(ctx)
			},
		},
	}
}
