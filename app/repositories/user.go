package repositories

import (
	"context"

	"github.com/ngocphuongnb/tetua/app/entities"
)

type UserRepository interface {
	Repository[entities.User, entities.UserFilter]
	ByUsername(ctx context.Context, username string) (*entities.User, error)
	ByProvider(ctx context.Context, providerName, providerId string) (*entities.User, error)
	ByUsernameOrEmail(ctx context.Context, username, email string) ([]*entities.User, error)
	CreateIfNotExistsByProvider(ctx context.Context, userData *entities.User) (*entities.User, error)
	Setting(ctx context.Context, id int, userData *entities.SettingMutation) (*entities.User, error)
}
