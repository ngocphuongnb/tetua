package mock

import (
	"github.com/ngocphuongnb/tetua/app/entities"
	repo "github.com/ngocphuongnb/tetua/app/mock/repository"
	"github.com/ngocphuongnb/tetua/app/repositories"
)

func Repositories() repositories.Repositories {
	return repositories.Repositories{
		File:       &repo.FileRepository{Repository: &repo.Repository[entities.File]{Name: "file"}},
		Post:       &repo.PostRepository{Repository: &repo.Repository[entities.Post]{Name: "post"}},
		Comment:    &repo.CommentRepository{Repository: &repo.Repository[entities.Comment]{Name: "comment"}},
		Role:       &repo.RoleRepository{Repository: &repo.Repository[entities.Role]{Name: "role"}},
		Topic:      &repo.TopicRepository{Repository: &repo.Repository[entities.Topic]{Name: "topic"}},
		User:       &repo.UserRepository{Repository: &repo.Repository[entities.User]{Name: "user"}},
		Permission: &repo.PermissionRepository{Repository: &repo.Repository[entities.Permission]{Name: "permission"}},
	}
}
func CreateRepositories() {
	repositories.File = &repo.FileRepository{Repository: &repo.Repository[entities.File]{Name: "file"}}
	repositories.Post = &repo.PostRepository{Repository: &repo.Repository[entities.Post]{Name: "post"}}
	repositories.Comment = &repo.CommentRepository{Repository: &repo.Repository[entities.Comment]{Name: "comment"}}
	repositories.Role = &repo.RoleRepository{Repository: &repo.Repository[entities.Role]{Name: "role"}}
	repositories.Topic = &repo.TopicRepository{Repository: &repo.Repository[entities.Topic]{Name: "topic"}}
	repositories.User = &repo.UserRepository{Repository: &repo.Repository[entities.User]{Name: "user"}}
	repositories.Permission = &repo.PermissionRepository{Repository: &repo.Repository[entities.Permission]{Name: "permission"}}
}
