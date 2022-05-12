package repositories

import (
	"context"

	"github.com/ngocphuongnb/tetua/app/entities"
)

// Repository will manipulate CRUD into the database only, not related to business rules. In microservices architecture, repositories can connect to microservices to get data.

var (
	File       FileRepository
	Role       RoleRepository
	Post       PostRepository
	Topic      TopicRepository
	User       UserRepository
	Permission PermissionRepository
	Comment    CommentRepository
	Setting    SettingRepository
)

type Repository[E entities.Entity, F entities.EntityFilter] interface {
	All(ctx context.Context) ([]*E, error)
	ByID(ctx context.Context, id int) (*E, error)
	DeleteByID(ctx context.Context, id int) error
	Create(ctx context.Context, data *E) (*E, error)
	Update(ctx context.Context, data *E) (*E, error)
	Count(ctx context.Context, filters ...*F) (int, error)
	Find(ctx context.Context, filters ...*F) ([]*E, error)
	Paginate(ctx context.Context, filters ...*F) (*entities.Paginate[E], error)
}

type Repositories struct {
	File       FileRepository
	User       UserRepository
	Post       PostRepository
	Role       RoleRepository
	Topic      TopicRepository
	Comment    CommentRepository
	Setting    SettingRepository
	Permission PermissionRepository
}

func New(config Repositories) {
	File = config.File
	User = config.User
	Post = config.Post
	Role = config.Role
	Topic = config.Topic
	Comment = config.Comment
	Setting = config.Setting
	Permission = config.Permission
}
