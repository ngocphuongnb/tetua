package entrepository

import (
	"context"
	"fmt"
	"log"

	"ariga.io/sqlcomment"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/logger"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/packages/entrepository/ent"
)

var Client *ent.Client

type Config struct {
	DB_DSN    string
	DB_DRIVER string
}

type Repository struct {
	Client *ent.Client
}

func New(cfg Config) repositories.Repositories {
	driverName := cfg.DB_DRIVER
	if driverName == "" {
		driverName = "mysql"
	}
	db, err := sql.Open(driverName, cfg.DB_DSN)

	if err != nil {
		log.Fatal(err)
	}

	if config.DB_QUERY_LOGGING {
		// Create sqlcomment driver which wraps sqlite driver.
		drv := sqlcomment.NewDriver(db,
			sqlcomment.WithDriverVerTag(),
			sqlcomment.WithTags(sqlcomment.Tags{
				sqlcomment.KeyApplication: "tetua",
				sqlcomment.KeyFramework:   "net/http",
			}),
		)
		drv = dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
			requestID := fmt.Sprintf("%v", ctx.Value("request_id"))
			logger.Get().WithContext(logger.Context{"request_id": requestID}).Debug(i...)
		})
		Client = ent.NewClient(ent.Driver(drv))
	} else {
		Client = ent.NewClient(ent.Driver(db))
	}

	if err := Client.Schema.Create(context.Background()); err != nil {
		log.Fatal(err)
	}

	return repositories.Repositories{
		File:       CreateFileRepository(Client),
		User:       CreateUserRepository(Client),
		Post:       CreatePostRepository(Client),
		Page:       CreatePageRepository(Client),
		Role:       CreateRoleRepository(Client),
		Topic:      CreateTopicRepository(Client),
		Comment:    CreateCommentRepository(Client),
		Setting:    &SettingRepository{&Repository{Client: Client}},
		Permission: CreatePermissionRepository(Client),
	}
}
