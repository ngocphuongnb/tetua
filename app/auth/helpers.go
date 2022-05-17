package auth

import (
	"github.com/ngocphuongnb/tetua/app/cache"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
)

func GetRolePermissions(roleID int) *entities.RolePermissions {
	for _, rolePermission := range cache.RolesPermissions {
		if rolePermission.RoleID == roleID {
			return rolePermission
		}
	}

	return &entities.RolePermissions{
		RoleID:      roleID,
		Permissions: []*entities.PermissionValue{},
	}
}

func GetRolePermission(roleID int, action string) *entities.PermissionValue {
	rolePermissions := GetRolePermissions(roleID)

	for _, permission := range rolePermissions.Permissions {
		if permission.Action == action {
			return permission
		}
	}

	return &entities.PermissionValue{}
}

func GetRolesFromIDs(IDs []int) []*entities.Role {
	result := []*entities.Role{}

	for _, role := range cache.Roles {
		for _, id := range IDs {
			if role.ID == id {
				result = append(result, role)
			}
		}
	}

	return result
}

func GetFile(c server.Context) error {
	fileID := c.ParamInt("id")
	file, err := repositories.File.ByID(c.Context(), fileID)

	if err != nil {
		return err
	}

	c.Locals("file", file)

	return nil
}

func GetPost(c server.Context) error {
	postIDParam := c.Param("id")

	if postIDParam == "new" {
		return nil
	}

	post, err := repositories.Post.ByID(c.Context(), c.ParamInt("id"))

	if err != nil {
		return err
	}

	c.Post(post)

	return nil
}

func GetComment(c server.Context) error {
	commentIDParam := c.Param("id")

	if commentIDParam == "new" {
		return nil
	}

	comment, err := repositories.Comment.ByID(c.Context(), c.ParamInt("id"))

	if err != nil {
		return err
	}

	c.Locals("comment", comment)

	return nil
}

func FileOwnerCheck(c server.Context) bool {
	if c.Param("id") == "new" {
		return true
	}

	user := c.User()
	file, ok := c.Locals("file").(*entities.File)

	if user == nil || file == nil {
		return false
	}

	if !ok {
		return false
	}

	if file.UserID != user.ID {
		return false
	}

	return true
}

func PostOwnerCheck(c server.Context) bool {
	if c.Param("id") == "new" {
		return true
	}

	user := c.User()
	post := c.Post()

	if user == nil || post == nil {
		return false
	}

	if post.UserID != user.ID {
		return false
	}

	c.Post(post)
	return true
}

func CommentOwnerCheck(c server.Context) bool {
	if c.Param("id") == "new" {
		return true
	}

	user := c.User()
	comment, ok := c.Locals("comment").(*entities.Comment)

	if user == nil || comment == nil {
		return false
	}

	if !ok {
		return false
	}

	if comment.UserID != user.ID {
		return false
	}

	return true
}

func AllowLoggedInUser(c server.Context) bool {
	user := c.User()
	if user == nil || user.ID == 0 {
		return false
	}

	return true
}

func AllowAll(c server.Context) bool {
	return true
}

func AllowNone(c server.Context) bool {
	return false
}
