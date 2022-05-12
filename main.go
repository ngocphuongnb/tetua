package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"sort"

	_ "ariga.io/sqlcomment"
	_ "github.com/Joker/hpp"
	_ "github.com/davecgh/go-spew/spew"
	"github.com/ngocphuongnb/tetua/app/asset"
	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/cache"
	"github.com/ngocphuongnb/tetua/app/cmd"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/fs"
	"github.com/ngocphuongnb/tetua/app/logger"
	"github.com/ngocphuongnb/tetua/app/repositories"

	"github.com/ngocphuongnb/tetua/app/web"
	ga "github.com/ngocphuongnb/tetua/packages/auth"
	ent "github.com/ngocphuongnb/tetua/packages/entrepository"
	"github.com/ngocphuongnb/tetua/packages/rclonefs"
	zap "github.com/ngocphuongnb/tetua/packages/zaplogger"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func prepare(workingDir string) {
	config.Init(workingDir)
	themeDir := path.Join(config.ROOT_DIR, "app/themes", config.APP_THEME)
	logger.New(zap.New(zap.Config{
		Development: config.DEVELOPMENT,
		LogFile:     path.Join(config.PRIVATE_DIR, "logs/tetua.log"),
	}))
	repositories.New(ent.New(ent.Config{DB_DSN: config.DB_DSN}))
	config.Settings(repositories.Setting.All(context.Background()))
	asset.Load(themeDir, false)
	fs.New(
		config.STORAGES.DefaultDisk,
		rclonefs.NewFromConfig(config.STORAGES),
	)
	auth.New(
		ga.NewLocal(),
		ga.NewGithub(&oauth2.Config{
			ClientID:     config.GITHUB_CLIENT_ID,
			ClientSecret: config.GITHUB_CLIENT_SECRET,
			RedirectURL:  config.Url("/auth/github/callback"),
			Endpoint:     github.Endpoint,
		}),
	)

	if err := cache.All(); err != nil {
		log.Fatal("Cache error", err)
	}
}

func getWd(c *cli.Context) string {
	workingDir := c.Args().First()

	if workingDir == "" {
		workingDir = "."
	}

	if _, err := os.Stat(workingDir); os.IsNotExist(err) {
		fmt.Println("Application directory not found")
		os.Exit(1)
	}

	// if err := os.Chdir(workingDir); err != nil {
	// 	fmt.Println("Can't change directory", err)
	// 	os.Exit(1)
	// }

	return workingDir
}

func main() {
	app := &cli.App{
		Name:  "tetua",
		Usage: "Easy blogging cms!",
		Commands: []*cli.Command{
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "Start tetua server",
				Action: func(c *cli.Context) error {
					prepare(getWd(c))
					web.NewServer(web.Config{
						JwtSigningKey: config.APP_KEY,
						Theme:         config.APP_THEME,
					}).Listen(":3000")
					return nil
				},
			},
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "Initialize the Tetua application",
				Action: func(c *cli.Context) error {
					return config.CreateConfigFile(getWd(c))
				},
			},
			{
				Name:    "setup",
				Aliases: []string{"s"},
				Usage:   "Setup the tetua application",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "username",
						Aliases:  []string{"u"},
						Usage:    "Admin username",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "password",
						Aliases:  []string{"p"},
						Usage:    "Admin password",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					prepare(getWd(c))
					return cmd.Setup(c.String("username"), c.String("password"))
				},
			},
			{
				Name:  "bundlestatic",
				Usage: "Bundle static files",
				Action: func(c *cli.Context) error {
					prepare(getWd(c))
					themeDir := path.Join(config.ROOT_DIR, "app/themes", config.APP_THEME)
					asset.Load(themeDir, true)
					return cmd.BundleStaticAssets()
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
