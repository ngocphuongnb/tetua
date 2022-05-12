package repositories

import (
	"context"

	"github.com/ngocphuongnb/tetua/app/entities"
)

type TopicRepository interface {
	Repository[entities.Topic, entities.TopicFilter]
	All(ctx context.Context) ([]*entities.Topic, error)
	ByName(ctx context.Context, name string) (*entities.Topic, error)
}
