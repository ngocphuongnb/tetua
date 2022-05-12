package config

type SettingItem struct {
	Name  string `json:"name" form:"name"`
	Value string `json:"value" form:"value"`
	Type  string `json:"type" form:"type"`
}

type SettingsMutation struct {
	Settings []*SettingItem `json:"settings,omitempty" form:"settings"`
}

var defaultSettings = []*SettingItem{
	{"app_name", "Tetua", "input"},
	{"app_desc", "Tetua", "input"},
	{"app_base_url", "", "input"},
	{"file_base_url", "", "input"},
	{"app_logo", "", "image"},
	{"contact_name", "", "input"},
	{"contact_email", "", "input"},
	{"twitter_site", "", "input"},
	{"inject_header", "", "textarea"},
	{"inject_footer", "", "textarea"},
	{"footer_content", "", "textarea"},
	{"auto_approve_user", "", "switch"},
	{"auto_approve_post", "", "switch"},
	{"auto_approve_comment", "", "switch"},
}
var settings = defaultSettings

func Settings(values []*SettingItem, overrideValues ...[]*SettingItem) {
	for _, s := range values {
		updateOrCreateSetting(s.Name, s.Value, s.Type)
	}

	for _, s := range overrideValues {
		for _, v := range s {
			updateOrCreateSetting(v.Name, v.Value, v.Type)
		}
	}
}

func updateOrCreateSetting(key, value, stype string) {
	for i, s := range settings {
		if s.Name == key {
			settings[i].Value = value
			return
		}
	}
	settings = append(settings, &SettingItem{key, value, stype})
}

func getSetting(key string) string {
	for _, s := range settings {
		if s.Name == key {
			return s.Value
		}
	}
	return ""
}

func Setting(key string, defaultValues ...string) string {
	value := getSetting(key)
	if value != "" {
		return value
	}

	if len(defaultValues) > 0 {
		return defaultValues[0]
	}

	return ""
}

func AllSettings() []*SettingItem {
	return settings
}
