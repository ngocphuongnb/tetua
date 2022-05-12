package managesetting

import (
	"fmt"

	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/services"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/views"
)

func Settings(c server.Context) (err error) {
	c.Meta().Title = "Settings"
	overrideSettings := []*config.SettingItem{}
	if config.Setting("app_base_url") == "" {
		overrideSettings = append(overrideSettings, &config.SettingItem{
			Type:  "input",
			Name:  "app_base_url",
			Value: c.BaseUrl(),
		})
	}
	if config.Setting("file_base_url") == "" {
		overrideSettings = append(overrideSettings, &config.SettingItem{
			Type:  "input",
			Name:  "file_base_url",
			Value: c.BaseUrl(),
		})
	}
	if config.Setting("footer_content") == "" {
		overrideSettings = append(overrideSettings, &config.SettingItem{
			Type: "textarea",
			Name: "footer_content",
			Value: fmt.Sprintf(
				`%s - %s.<br/>Built on Tetua - An open source CMS platform for Blogging.<br/>%s © 2021 - 2022.`,
				config.Setting("app_name"),
				config.Setting("app_desc"),
				config.Setting("app_name"),
			),
		})
	}

	if len(overrideSettings) > 0 {
		config.Settings(overrideSettings)
	}

	return c.Render(views.ManageSettings(config.AllSettings()))
}

func Save(c server.Context) (err error) {
	data := &config.SettingsMutation{}
	if err = c.BodyParser(data); err != nil {
		c.WithError("Error parsing body", err)
	}

	var settingItems = utils.SliceFilter(data.Settings, func(s *config.SettingItem) bool {
		return s != nil
	})

	for _, s := range config.AllSettings() {
		if s.Type != "image" {
			continue
		}

		if image, err := services.SaveFile(c, s.Name); err != nil {
			c.WithError("Error saving setting image: "+s.Name, err)
		} else if image != nil {
			settingItems = append(settingItems, &config.SettingItem{
				Type:  "image",
				Name:  s.Name,
				Value: image.Url(),
			})
		}
	}

	var settings = utils.SliceMap(settingItems, func(s *config.SettingItem) *entities.Setting {
		return &entities.Setting{
			Name:  s.Name,
			Value: s.Value,
			Type:  s.Type,
		}
	})

	if err = repositories.Setting.Save(c.Context(), settings); err != nil {
		c.WithError("Error saving settings", err)
	}

	config.Settings(settingItems)

	return c.Redirect("/manage/settings")
}
