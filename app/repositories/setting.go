package repositories

import (
	"context"

	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
)

type SettingRepository interface {
	All(ctx context.Context) []*config.SettingItem
	Save(ctx context.Context, data []*entities.Setting) error
}
