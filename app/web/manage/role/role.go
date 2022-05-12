package managerole

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/cache"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/ngocphuongnb/tetua/views"
)

func Index(c server.Context) (err error) {
	status := http.StatusOK
	roles, err := repositories.Role.All(c.Context())
	c.Meta().Title = "Manage roles"

	if err != nil {
		status = http.StatusBadRequest
		c.WithError("Load all roles error", err)
	}

	return c.Status(status).Render(views.ManageRoleIndex(roles))
}

func Compose(c server.Context) (err error) {
	return composeView(c, &entities.RoleMutation{}, false)
}

func Save(c server.Context) (err error) {
	var role *entities.Role
	data := getRoleSaveData(c)

	if role, err = getProcessingRole(c); err != nil {
		c.WithError("Query editting role error", err)
	}

	if c.Messages().Length() > 0 {
		return composeView(c, data, true)
	}

	role.Root = data.Root
	role.Name = data.Name
	role.Description = data.Description

	if role.ID > 0 {
		role, err = repositories.Role.Update(c.Context(), role)
	} else {
		role, err = repositories.Role.Create(c.Context(), role)
	}

	if err != nil {
		c.WithError("Error saving role", err)
		return composeView(c, data, true)
	}

	if !role.Root {
		if err := repositories.Role.SetPermissions(c.Context(), role.ID, data.Permissions); err != nil {
			c.WithError("Error saving role", err)
			return composeView(c, data, true)
		}
	}

	if err := cache.CachePermissions(c.Context()); err != nil {
		c.WithError("Error caching permissions", err)
		return composeView(c, data, true)
	}
	return c.Redirect("/manage/roles/" + strconv.Itoa(role.ID))
}

func Delete(c server.Context) error {
	role, err := getProcessingRole(c)

	if role.ID < 4 {
		return c.Status(http.StatusBadRequest).SendString("Error deleting role")
	}

	if err != nil {
		c.Logger().Error("Error deleting role", err)
		return c.Status(http.StatusBadRequest).SendString("Error deleting role")
	}

	if err := repositories.Role.DeleteByID(c.Context(), role.ID); err != nil {
		c.Logger().Error("Error deleting role", err)
		return c.Status(http.StatusBadRequest).SendString("Error deleting role")
	}

	return c.Status(http.StatusOK).SendString("Success")
}

func getProcessingRole(c server.Context) (role *entities.Role, err error) {
	if c.Param("id") == "new" {
		return &entities.Role{}, nil
	}

	return repositories.Role.ByID(c.Context(), c.ParamInt("id"))
}

func composeView(c server.Context, data *entities.RoleMutation, isSave bool) (err error) {
	role, err := getProcessingRole(c)
	c.Meta().Title = "Create Role"

	if err != nil {
		c.WithError("Query editting role error", err)
	} else if !isSave {
		data.Root = role.Root
		data.Name = role.Name
		data.Description = role.Description
	}

	if role.ID > 0 {
		c.Meta().Title = "Edit Role: " + role.Name
	}

	var rolePermissions []*entities.PermissionValue
	var rolePermissionsByActions = map[string]string{}

	for _, permission := range role.Permissions {
		rolePermissionsByActions[permission.Action] = permission.Value
	}

	for _, permissionValue := range auth.ActionConfigs {
		if value, ok := rolePermissionsByActions[permissionValue.Action]; ok {
			rolePermissions = append(rolePermissions, &entities.PermissionValue{
				Action: permissionValue.Action,
				Value:  entities.GetPermTypeValue(value),
			})
		} else {
			rolePermissions = append(rolePermissions, &entities.PermissionValue{
				Action: permissionValue.Action,
				Value:  entities.PermType(permissionValue.DefaultValue),
			})
		}
	}

	return c.Render(views.ManageRoleCompose(role.ID, data, rolePermissions))
}

func getRoleSaveData(c server.Context) *entities.RoleMutation {
	data := &entities.RoleMutation{}
	if err := c.BodyParser(data); err != nil {
		c.WithError("Error parsing body", err)
		return data
	}

	data.Name = utils.SanitizePlainText(strings.TrimSpace(data.Name))
	data.Description = utils.SanitizePlainText(strings.TrimSpace(data.Description))

	if data.Name == "" || len(data.Name) > 250 {
		c.Messages().AppendError("Name is required and can't be more than 250 characters")
	}

	return data
}
