package entrepository

import (
	"context"

	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/logger"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent"
)

type SettingRepository struct {
	*Repository
}

func (ur *SettingRepository) All(ctx context.Context) []*config.SettingItem {
	result := make([]*config.SettingItem, 0)
	settings, err := ur.Client.Setting.Query().All(ctx)

	if err != nil {
		logger.Get().Error(err)
		return result
	}

	for _, setting := range settings {
		result = append(result, &config.SettingItem{
			Name:  setting.Name,
			Value: setting.Value,
			Type:  "",
		})
	}

	return result
}

func (ur *SettingRepository) Save(ctx context.Context, settings []*entities.Setting) error {
	var builders []*ent.SettingCreate

	for _, s := range settings {
		us := ur.Client.Setting.
			Create().
			SetName(s.Name).
			SetValue(s.Value).
			SetType(s.Type)
		builders = append(builders, us)
	}

	return ur.Client.Setting.
		CreateBulk(builders...).
		OnConflict().
		UpdateNewValues().
		Exec(ctx)
}
