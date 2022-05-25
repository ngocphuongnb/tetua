package manage

import (
	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/server"
	managecomment "github.com/ngocphuongnb/tetua/app/web/manage/comment"
	managefile "github.com/ngocphuongnb/tetua/app/web/manage/file"
	managepage "github.com/ngocphuongnb/tetua/app/web/manage/page"
	managepost "github.com/ngocphuongnb/tetua/app/web/manage/post"
	managerole "github.com/ngocphuongnb/tetua/app/web/manage/role"
	managesetting "github.com/ngocphuongnb/tetua/app/web/manage/setting"
	managetopic "github.com/ngocphuongnb/tetua/app/web/manage/topic"
	manageuser "github.com/ngocphuongnb/tetua/app/web/manage/user"
)

func manageAuthConfig(action string) *server.AuthConfig {
	return auth.Config(&server.AuthConfig{
		Action:       action,
		DefaultValue: entities.PERM_NONE,
		OwnCheckFN:   auth.AllowNone,
	})
}

var (
	authManage               = manageAuthConfig("manage")
	authManageTopicList      = manageAuthConfig("manage.topic.list")
	authManageTopicCompose   = manageAuthConfig("manage.topic.compose")
	authManageTopicSave      = manageAuthConfig("manage.topic.save")
	authManageTopicDelete    = manageAuthConfig("manage.topic.delete")
	authManagePostList       = manageAuthConfig("manage.post.list")
	authManagePostApprove    = manageAuthConfig("manage.post.approve")
	authManagePageList       = manageAuthConfig("manage.page.list")
	authManagePageCompose    = manageAuthConfig("manage.page.compose")
	authManagePageSave       = manageAuthConfig("manage.page.save")
	authManagePageDelete     = manageAuthConfig("manage.page.delete")
	authManageRoleList       = manageAuthConfig("manage.role.list")
	authManageRoleCompose    = manageAuthConfig("manage.role.compose")
	authManageRoleSave       = manageAuthConfig("manage.role.save")
	authManageRoleDelete     = manageAuthConfig("manage.role.delete")
	authManageUserList       = manageAuthConfig("manage.user.list")
	authManageUserCompose    = manageAuthConfig("manage.user.compose")
	authManageUserSave       = manageAuthConfig("manage.user.save")
	authManageuserdelete     = manageAuthConfig("manage.user.delete")
	authManageSettingCompose = manageAuthConfig("manage.setting.compose")
	authManageSettingSave    = manageAuthConfig("manage.setting.save")
	authManageCommentList    = manageAuthConfig("manage.comment.list")
	authManageFileList       = manageAuthConfig("manage.file.list")
)

func RegisterRoutes(s server.Server) {
	manage := s.Group("/manage")
	manage.Get("", Manage, authManage)

	topic := manage.Group("/topics")
	topic.Get("", managetopic.Index, authManageTopicList)
	topic.Get("/:id", managetopic.Compose, authManageTopicCompose)
	topic.Post("/:id", managetopic.Save, authManageTopicSave)
	topic.Delete("/:id", managetopic.Delete, authManageTopicDelete)

	post := manage.Group("/posts")
	post.Get("", managepost.Index, authManagePostList)
	post.Post("/:id/approve", managepost.Approve, authManagePostApprove)

	page := manage.Group("/pages")
	page.Get("", managepage.Index, authManagePageList)
	page.Get("/:id", managepage.Compose, authManagePageCompose)
	page.Post("/:id", managepage.Save, authManagePageSave)
	page.Delete("/:id", managepage.Delete, authManagePageDelete)

	role := manage.Group("/roles")
	role.Get("", managerole.Index, authManageRoleList)
	role.Get("/:id", managerole.Compose, authManageRoleCompose)
	role.Post("/:id", managerole.Save, authManageRoleSave)
	role.Delete("/:id", managerole.Delete, authManageRoleDelete)

	user := manage.Group("/users")
	user.Get("", manageuser.Index, authManageUserList)
	user.Get("/:id", manageuser.Compose, authManageUserCompose)
	user.Post("/:id", manageuser.Save, authManageUserSave)
	user.Delete("/:id", manageuser.Delete, authManageuserdelete)

	setting := manage.Group("/settings")
	setting.Get("", managesetting.Settings, authManageSettingCompose)
	setting.Post("", managesetting.Save, authManageSettingSave)

	comment := manage.Group("/comments")
	comment.Get("", managecomment.Index, authManageCommentList)

	file := manage.Group("/files")
	file.Get("", managefile.Index, authManageFileList)
}
