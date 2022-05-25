package config

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/ngocphuongnb/tetua/app/test"
	"github.com/stretchr/testify/assert"
)

const sampleConfigContent = `{
	"app_env": "testing",
  "app_key": "app_key_test",
  "app_token_key": "app_token_key_test",
  "app_port": "3001",
  "app_theme": "test_theme",
  "app_base_url": "http://localhost:3001",
  "cookie_uuid": "test_uuid",
  "db_dsn": "root:123@tcp(127.0.0.1:3306)/tetua",
  "github_client_id": "github_client_id_test",
  "github_client_secret": "github_client_secret_test",
  "show_tetua_block": true,
	"db_query_logging": true,
  "storage": {
    "default_disk": "local_public_test",
    "disks": [
      {
        "name": "local_public_test",
        "driver": "local",
        "root": "./public/storage",
        "base_url": "http://localhost:3001"
      },
      {
        "name": "local_public_test",
        "driver": "local",
        "root": "./public/storage",
        "base_url": "http://localhost:3001"
      },
      {
        "name": "local_private_test",
        "driver": "local",
        "root": "./private/storage"
      },
      {
        "name": "s3_public_test",
        "driver": "s3",
        "root": "/files",
        "provider": "DigitalOcean",
        "endpoint": "ams3.digitaloceanspaces.com",
        "region": "ams3",
        "bucket": "",
        "access_key_id": "",
        "secret_access_key": ""
      }
    ]
  }
}`

func createCustomWorkingDir(configContent string) string {
	workingDir := test.CreateDir("app-")
	configFile := path.Join(workingDir, "config.json")

	if err := ioutil.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		panic(err)
	}

	return workingDir
}

func TestConfigInitEmptyAppKey(t *testing.T) {
	workingDir := createCustomWorkingDir(strings.ReplaceAll(sampleConfigContent, "app_key_test", ""))
	APP_KEY = ""
	os.Setenv("APP_KEY ", "")
	defer test.RecoverPanic(t, "APP_KEY is not set, please set it in config.json or setting the environment variable.\n")
	defer os.RemoveAll(workingDir)

	Init(workingDir)
}

func TestConfigInitEmptyDBDsn(t *testing.T) {
	workingDir := createCustomWorkingDir(strings.ReplaceAll(sampleConfigContent, "root:123@tcp(127.0.0.1:3306)/tetua", ""))
	DB_DSN = ""
	os.Setenv("DB_DSN ", "")
	defer test.RecoverPanic(t, "DB_DSN is not set, please set it in config.json or setting the environment variable.\n")
	defer os.RemoveAll(workingDir)

	Init(workingDir)
}

func TestConfigInit(t *testing.T) {
	workingDir := createCustomWorkingDir(sampleConfigContent)
	defer os.RemoveAll(workingDir)

	os.Setenv("APP_ENV", "development")
	Init(workingDir)
	assert.Equal(t, settings, defaultSettings)
	assert.Equal(t, workingDir, WD)
	assert.Equal(t, "development", APP_ENV)
	assert.Equal(t, true, DEVELOPMENT)
}

func TestParseConfigFileFailed(t *testing.T) {
	defer test.RecoverPanic(t, "unexpected end of JSON input")
	WD = createCustomWorkingDir("")
	defer os.RemoveAll(WD)
	parseConfigFile()
}

func TestParseConfigFile(t *testing.T) {
	WD = createCustomWorkingDir(sampleConfigContent)
	defer os.RemoveAll(WD)
	parseConfigFile()

	assert.Equal(t, "testing", APP_ENV)
	assert.Equal(t, "3001", APP_PORT)
	assert.Equal(t, "root:123@tcp(127.0.0.1:3306)/tetua", DB_DSN)
	assert.Equal(t, "app_key_test", APP_KEY)
	assert.Equal(t, "app_token_key_test", APP_TOKEN_KEY)
	assert.Equal(t, "test_theme", APP_THEME)
	assert.Equal(t, "test_uuid", COOKIE_UUID)
	assert.Equal(t, true, SHOW_TETUA_BLOCK)
	assert.Equal(t, true, DB_QUERY_LOGGING)
	assert.Equal(t, 5, len(STORAGES.DiskConfigs))
	assert.Equal(t, "local_public_test", STORAGES.DiskConfigs[2].Name)
	assert.Equal(t, "local_private_test", STORAGES.DiskConfigs[3].Name)
	assert.Equal(t, "s3_public_test", STORAGES.DiskConfigs[4].Name)
	assert.Equal(t, "/files", STORAGES.DiskConfigs[0].BaseUrlFn())
}

func TestParseEnv(t *testing.T) {
	os.Setenv("APP_ENV", "env_development")
	os.Setenv("APP_PORT", "3002")
	os.Setenv("DB_DSN", "env_root:123@tcp")
	os.Setenv("APP_KEY", "env_app_key_test")
	os.Setenv("APP_TOKEN_KEY", "env_app_token_key_test")
	os.Setenv("APP_THEME", "env_test_theme")
	os.Setenv("COOKIE_UUID", "env_test_uuid")
	os.Setenv("SHOW_TETUA_BLOCK", "true")
	os.Setenv("DB_QUERY_LOGGING", "true")

	parseENV()

	assert.Equal(t, "env_development", APP_ENV)
	assert.Equal(t, "3002", APP_PORT)
	assert.Equal(t, "env_root:123@tcp", DB_DSN)
	assert.Equal(t, "env_app_key_test", APP_KEY)
	assert.Equal(t, "env_app_token_key_test", APP_TOKEN_KEY)
	assert.Equal(t, "env_app_token_key_test", APP_TOKEN_KEY)
	assert.Equal(t, "env_test_theme", APP_THEME)
	assert.Equal(t, "env_test_uuid", COOKIE_UUID)
	assert.Equal(t, true, SHOW_TETUA_BLOCK)
	assert.Equal(t, true, DB_QUERY_LOGGING)

	os.Setenv("SHOW_TETUA_BLOCK", "false")
	os.Setenv("DB_QUERY_LOGGING", "false")

	parseENV()

	assert.Equal(t, false, SHOW_TETUA_BLOCK)
	assert.Equal(t, false, DB_QUERY_LOGGING)
}

func TestUrl(t *testing.T) {
	Settings([]*SettingItem{{Name: "app_base_url", Value: "http://localhost:3002"}})
	assert.Equal(t, "http://localhost:3002", Setting("app_base_url"))
}

func TestCreateConfigFile(t *testing.T) {
	WD = test.CreateDir("app-")
	defer os.RemoveAll(WD)
	assert.Equal(t, nil, CreateConfigFile(WD))
	assert.Equal(t, nil, CreateConfigFile(WD))
	parseConfigFile()

	assert.Equal(t, "production", APP_ENV)
	assert.Equal(t, "3002", APP_PORT)
	assert.Equal(t, true, DB_DSN != "")
}
