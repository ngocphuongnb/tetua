package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/ngocphuongnb/tetua/app/fs"
	"github.com/ngocphuongnb/tetua/app/utils"
)

type ConfigFile struct {
	APP_ENV              string            `json:"app_env"`
	APP_KEY              string            `json:"app_key"`
	APP_TOKEN_KEY        string            `json:"app_token_key,omitempty"`
	APP_PORT             string            `json:"app_port"`
	APP_THEME            string            `json:"app_theme,omitempty"`
	DB_DSN               string            `json:"db_dsn"`
	GITHUB_CLIENT_ID     string            `json:"github_client_id"`
	GITHUB_CLIENT_SECRET string            `json:"github_client_secret"`
	COOKIE_UUID          string            `json:"cookie_uuid,omitempty"`
	SHOW_TETUA_BLOCK     bool              `json:"show_tetua_block,omitempty"`
	DB_QUERY_LOGGING     bool              `json:"db_query_logging"`
	STORAGES             *fs.StorageConfig `json:"storage,omitempty"`
}

var (
	WD                   = "."
	DEVELOPMENT          = false
	APP_VERSION          = "0.0.1"
	STORAGES             = &fs.StorageConfig{}
	APP_ENV              = ""
	APP_KEY              = ""
	APP_TOKEN_KEY        = "token"
	APP_PORT             = "3000"
	APP_THEME            = "default"
	DB_DSN               = ""
	DB_QUERY_LOGGING     = false
	GITHUB_CLIENT_ID     = ""
	GITHUB_CLIENT_SECRET = ""
	ROOT_DIR             = ""
	PUBLIC_DIR           = "public"
	PRIVATE_DIR          = "private"
	COOKIE_UUID          = "uuid"
	SHOW_TETUA_BLOCK     = false
)

func ConfigError(name string) {
	panic(fmt.Sprintf(
		"%s is not set, please set it in config.json or setting the environment variable.\n",
		name,
	))
}

func Init(workingDir string) {
	WD = workingDir
	PUBLIC_DIR = path.Join(WD, "public")
	PRIVATE_DIR = path.Join(WD, "private")

	parseConfigFile()
	parseENV()
	DEVELOPMENT = APP_ENV == "development"

	if DEVELOPMENT {
		ROOT_DIR = utils.GetRootDir()
	}

	if APP_KEY == "" {
		ConfigError("APP_KEY")
	}

	if DB_DSN == "" {
		ConfigError("DB_DSN")
	}
}

func Url(path string) string {
	appBase := strings.TrimRight(Setting("app_base_url"), "/")

	if path == "" {
		return appBase + "/"
	}

	path = strings.TrimLeft(path, "/")
	path = fmt.Sprintf("%s/%s", appBase, path)
	return path
}

func CreateConfigFile(workingDir string) (err error) {
	configFile := path.Join(workingDir, "config.json")

	if _, err := os.Stat(configFile); err == nil {
		fmt.Println("config.json already exists", configFile)
		return nil
	}

	cfg := &ConfigFile{}
	cfg.APP_ENV = "production"
	cfg.SHOW_TETUA_BLOCK = false
	cfg.APP_PORT = APP_PORT
	cfg.APP_KEY = utils.RandomString(32)

	file, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(configFile, file, 0644); err != nil {
		return err
	}

	fmt.Println("Config file created at", configFile)

	return nil
}

func parseConfigFile() {
	configFile := path.Join(WD, "config.json")
	cfg := &ConfigFile{}

	if _, err := os.Stat(configFile); !errors.Is(err, os.ErrNotExist) {
		file, _ := ioutil.ReadFile(configFile)

		if err := json.Unmarshal([]byte(file), cfg); err != nil {
			panic(err)
		}

		if cfg.APP_ENV != "" {
			APP_ENV = cfg.APP_ENV
		}

		if cfg.APP_PORT != "" {
			APP_PORT = cfg.APP_PORT
		}

		if cfg.DB_DSN != "" {
			DB_DSN = cfg.DB_DSN
		}

		if cfg.APP_KEY != "" {
			APP_KEY = cfg.APP_KEY
		}

		if cfg.APP_TOKEN_KEY != "" {
			APP_TOKEN_KEY = cfg.APP_TOKEN_KEY
		}

		if cfg.APP_THEME != "" {
			APP_THEME = cfg.APP_THEME
		}

		if cfg.GITHUB_CLIENT_ID != "" {
			GITHUB_CLIENT_ID = cfg.GITHUB_CLIENT_ID
		}

		if cfg.GITHUB_CLIENT_SECRET != "" {
			GITHUB_CLIENT_SECRET = cfg.GITHUB_CLIENT_SECRET
		}

		if cfg.COOKIE_UUID != "" {
			COOKIE_UUID = cfg.COOKIE_UUID
		}

		SHOW_TETUA_BLOCK = cfg.SHOW_TETUA_BLOCK
		DB_QUERY_LOGGING = cfg.DB_QUERY_LOGGING
		STORAGES = &fs.StorageConfig{
			DefaultDisk: "local_public",
			DiskConfigs: []*fs.DiskConfig{{
				Name:   "local_public",
				Driver: "local",
				Root:   path.Join(PUBLIC_DIR, "files"),
				BaseUrlFn: func() string {
					return Setting("file_base_url") + "/files"
				},
			}, {
				Name:   "local_private",
				Driver: "local",
				Root:   path.Join(PRIVATE_DIR, "storage"),
			}},
		}

		if cfg.STORAGES != nil {
			if cfg.STORAGES.DefaultDisk != "" {
				STORAGES.DefaultDisk = cfg.STORAGES.DefaultDisk
			}

			if cfg.STORAGES.DiskConfigs != nil && len(cfg.STORAGES.DiskConfigs) > 0 {
				for _, disk := range cfg.STORAGES.DiskConfigs {
					STORAGES.DiskConfigs = utils.SliceAppendIfNotExists(STORAGES.DiskConfigs, disk, func(t *fs.DiskConfig) bool {
						return t.Name == disk.Name
					})
				}
			}
		}
	}
}

func parseENV() {
	if os.Getenv("APP_ENV") != "" {
		APP_ENV = os.Getenv("APP_ENV")
	}

	if os.Getenv("APP_PORT") != "" {
		APP_PORT = os.Getenv("APP_PORT")
	}

	if os.Getenv("DB_DSN") != "" {
		DB_DSN = os.Getenv("DB_DSN")
	}

	if os.Getenv("APP_KEY") != "" {
		APP_KEY = os.Getenv("APP_KEY")
	}

	if os.Getenv("APP_TOKEN_KEY") != "" {
		APP_TOKEN_KEY = os.Getenv("APP_TOKEN_KEY")
	}

	if os.Getenv("APP_THEME") != "" {
		APP_THEME = os.Getenv("APP_THEME")
	}

	if os.Getenv("GITHUB_CLIENT_ID") != "" {
		GITHUB_CLIENT_ID = os.Getenv("GITHUB_CLIENT_ID")
	}

	if os.Getenv("GITHUB_CLIENT_SECRET") != "" {
		GITHUB_CLIENT_SECRET = os.Getenv("GITHUB_CLIENT_SECRET")
	}

	if os.Getenv("COOKIE_UUID") != "" {
		COOKIE_UUID = os.Getenv("COOKIE_UUID")
	}

	if os.Getenv("DB_QUERY_LOGGING") != "" {
		DB_QUERY_LOGGING = strings.ToLower(os.Getenv("DB_QUERY_LOGGING")) == "true"
	}

	if os.Getenv("SHOW_TETUA_BLOCK") != "" {
		SHOW_TETUA_BLOCK = strings.ToLower(os.Getenv("SHOW_TETUA_BLOCK")) == "true"
	}
}
