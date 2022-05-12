package entrepository

import (
	"context"
	"fmt"

	e "github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent/topic"
)

type TopicRepository struct {
	*BaseRepository[e.Topic, ent.Topic, *ent.TopicQuery, *e.TopicFilter]
}

func (c *TopicRepository) ByName(ctx context.Context, name string) (*e.Topic, error) {
	t, err := c.Client.Topic.Query().Where(topic.NameEQ(name)).Only(ctx)
	if err != nil {
		return nil, EntError(err, fmt.Sprintf("topic not found with name: %s", name))
	}
	return entTopicToTopic(t), err
}

func CreateTopicRepository(client *ent.Client) *TopicRepository {
	return &TopicRepository{
		BaseRepository: &BaseRepository[e.Topic, ent.Topic, *ent.TopicQuery, *e.TopicFilter]{
			Name:      "topic",
			Client:    client,
			ConvertFn: entTopicToTopic,
			ByIDFn: func(ctx context.Context, client *ent.Client, id int) (*ent.Topic, error) {
				return client.Topic.Query().Where(topic.IDEQ(id)).Only(ctx)
			},
			DeleteByIDFn: func(ctx context.Context, client *ent.Client, id int) error {
				return client.Topic.DeleteOneID(id).Exec(ctx)
			},
			CreateFn: func(ctx context.Context, client *ent.Client, data *e.Topic) (topic *ent.Topic, err error) {
				tc := client.Topic.Create().
					SetName(data.Name).
					SetSlug(data.Slug).
					SetDescription(data.Description).
					SetContent(data.Content).
					SetContentHTML(data.ContentHTML)
				if data.ParentID != 0 {
					tc = tc.SetParentID(data.ParentID)
				}
				return tc.Save(ctx)
			},
			UpdateFn: func(ctx context.Context, client *ent.Client, data *e.Topic) (topic *ent.Topic, err error) {
				tc := client.Topic.UpdateOneID(data.ID).
					SetName(data.Name).
					SetSlug(data.Slug).
					SetDescription(data.Description).
					SetContent(data.Content).
					SetContentHTML(data.ContentHTML)
				if data.ParentID != 0 {
					tc = tc.SetParentID(data.ParentID)
				}
				return tc.Save(ctx)
			},
			QueryFilterFn: func(client *ent.Client, filters ...*e.TopicFilter) *ent.TopicQuery {
				query := client.Topic.Query().Where(topic.DeletedAtIsNil())
				if len(filters) > 0 && filters[0].Search != "" {
					query = query.Where(topic.NameContainsFold(filters[0].Search))
				}
				return query
			},
			FindFn: func(ctx context.Context, query *ent.TopicQuery, filters ...*e.TopicFilter) ([]*ent.Topic, error) {
				page, limit, sorts := getPaginateParams(filters...)
				return query.
					Limit(limit).
					Offset((page - 1) * limit).
					Order(sorts...).
					All(ctx)
			},
		},
	}
}
