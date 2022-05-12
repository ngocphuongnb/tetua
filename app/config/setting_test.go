package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettings(t *testing.T) {
	var newSettings = []*SettingItem{
		{"app_name", "App Testing", "input"},
		{"test_key", "Test value", "input"},
	}
	var overrideSettings = []*SettingItem{
		{"app_name", "App Testing", "input"},
		{"test_key", "Test value 2", "input"},
	}

	assert.Equal(t, "Tetua", Setting("app_desc"))

	Settings(defaultSettings)
	assert.Equal(t, "Tetua", Setting("app_name"))

	Settings(newSettings)
	assert.Equal(t, "App Testing", Setting("app_name"))
	assert.Equal(t, "Test value", Setting("test_key"))

	Settings(newSettings, overrideSettings)
	assert.Equal(t, "Test value 2", Setting("test_key"))
	assert.Equal(t, "Tetua", Setting("app_desc"))
	assert.Equal(t, "", Setting("unknown"))
	assert.Equal(t, "default", Setting("unknown", "default"))
	assert.Equal(t, settings, AllSettings())
}
