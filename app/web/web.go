package web

import (
	"path"

	"github.com/ngocphuongnb/tetua/app/asset"
	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/middlewares"
	"github.com/ngocphuongnb/tetua/app/server"
	webcomment "github.com/ngocphuongnb/tetua/app/web/comment"
	"github.com/ngocphuongnb/tetua/app/web/manage"
	webpost "github.com/ngocphuongnb/tetua/app/web/post"
	websetting "github.com/ngocphuongnb/tetua/app/web/setting"
	websitemap "github.com/ngocphuongnb/tetua/app/web/sitemap"
	webuser "github.com/ngocphuongnb/tetua/app/web/user"
	fiber "github.com/ngocphuongnb/tetua/packages/fiberserver"
)

type Config struct {
	JwtSigningKey string
	Theme         string
}

var (
	authPostCompose = auth.Config(&server.AuthConfig{
		Action:       "post.compose",
		DefaultValue: entities.PERM_OWN,
		Prepare:      auth.GetPost,
		OwnCheckFN:   auth.PostOwnerCheck,
	})

	authPostSave = auth.Config(&server.AuthConfig{
		Action:       "post.save",
		DefaultValue: entities.PERM_OWN,
		Prepare:      auth.GetPost,
		OwnCheckFN:   auth.PostOwnerCheck,
	})

	authPostDelete = auth.Config(&server.AuthConfig{
		Action:       "post.delete",
		DefaultValue: entities.PERM_OWN,
		Prepare:      auth.GetPost,
		OwnCheckFN:   auth.PostOwnerCheck,
	})

	authPostList = auth.Config(&server.AuthConfig{
		Action:       "post.list",
		DefaultValue: entities.PERM_OWN,
		OwnCheckFN:   auth.AllowLoggedInUser,
	})

	authPostView = auth.Config(&server.AuthConfig{
		Action:       "post.view",
		DefaultValue: entities.PERM_ALL,
	})

	authCommentList = auth.Config(&server.AuthConfig{
		Action:       "comment.list",
		DefaultValue: entities.PERM_OWN,
		OwnCheckFN:   auth.AllowLoggedInUser,
	})

	authCommentSave = auth.Config(&server.AuthConfig{
		Action:       "comment.save",
		DefaultValue: entities.PERM_OWN,
		OwnCheckFN:   auth.CommentOwnerCheck,
	})

	authCommentDelete = auth.Config(&server.AuthConfig{
		Action:       "comment.delete",
		DefaultValue: entities.PERM_OWN,
		OwnCheckFN:   auth.CommentOwnerCheck,
	})

	authFileUpload = auth.Config(&server.AuthConfig{
		Action:       "file.upload",
		DefaultValue: entities.PERM_OWN,
		OwnCheckFN:   auth.AllowLoggedInUser,
	})

	authFileList = auth.Config(&server.AuthConfig{
		Action:       "file.list",
		DefaultValue: entities.PERM_OWN,
		OwnCheckFN:   auth.AllowLoggedInUser,
	})

	authFileDelete = auth.Config(&server.AuthConfig{
		Action:       "file.delete",
		DefaultValue: entities.PERM_OWN,
		Prepare:      auth.GetFile,
		OwnCheckFN:   auth.FileOwnerCheck,
	})

	authUserProfile = auth.Config(&server.AuthConfig{
		Action:       "user.profile",
		DefaultValue: entities.PERM_ALL,
	})

	authUserSettingCompose = auth.Config(&server.AuthConfig{
		Action:       "user.setting.compose",
		DefaultValue: entities.PERM_OWN,
		OwnCheckFN:   auth.AllowLoggedInUser,
	})

	authUserSettingSave = auth.Config(&server.AuthConfig{
		Action:       "user.setting.save",
		DefaultValue: entities.PERM_OWN,
		OwnCheckFN:   auth.AllowLoggedInUser,
	})

	authTopicView = auth.Config(&server.AuthConfig{
		Action:       "topic.view",
		DefaultValue: entities.PERM_ALL,
	})

	authTopicFeed = auth.Config(&server.AuthConfig{
		Action:       "topic.feed",
		DefaultValue: entities.PERM_ALL,
	})
)

func NewServer(cfg Config) server.Server {
	s := fiber.New(fiber.Config{
		JwtSigningKey: cfg.JwtSigningKey,
		AppName:       config.Setting("app_name"),
	})
	s.Register(auth.Routes)
	s.Static("/", path.Join(config.WD, "public"))

	for _, assetFile := range asset.All() {
		assetPath := path.Join("assets", assetFile.Name)
		if config.DEVELOPMENT {
			s.Static(assetPath, assetFile.Path)
		} else {
			if assetFile.DisableInline {
				func(assetPath string, assetFile *asset.StaticAsset) {
					s.Get(assetPath, func(c server.Context) error {
						return c.SendString(assetFile.Content)
					})
				}(assetPath, assetFile)
			}
		}
	}

	s.Use(middlewares.All()...)
	manage.RegisterRoutes(s)

	compose := s.Group("/posts/:id")
	compose.Get("", webpost.Compose, authPostCompose)
	compose.Post("", webpost.Save, authPostSave)
	compose.Delete("", webpost.Delete, authPostDelete)

	comment := s.Group("/comments")
	comment.Get("", webcomment.List, authCommentList)
	comment.Post("/:id", webcomment.Save, authCommentSave)
	comment.Delete("/:id", webcomment.Delete, authCommentDelete)

	file := s.Group("/files")
	file.Post("/upload", Upload, authFileUpload)
	file.Get("", FileList, authFileList)
	file.Delete("/:id", FileDelete, authFileDelete)

	profile := s.Group("/u")
	profile.Get("/:username", webuser.Profile, authUserProfile)

	s.Get("", Index)
	s.Get("/search", Search)
	s.Get("/feed", Feed)
	s.Get("/inactive", webuser.Inactive)
	s.Get("/login", webuser.Login)
	s.Post("/login", webuser.PostLogin)
	s.Get("/logout", webuser.Logout)
	s.Get("/sitemap/index.xml", websitemap.Index)
	s.Get("/sitemap/topics.xml", websitemap.Topic)
	s.Get("/sitemap/users.xml", websitemap.User)
	s.Get("/sitemap/posts-:page.xml", websitemap.Post)
	s.Get("/settings", websetting.Index, authUserSettingCompose)
	s.Post("/settings", websetting.Save, authUserSettingSave)

	s.Get("/posts", webpost.List, authPostList)
	s.Get("/:slug.html", webpost.View, authPostView)
	s.Get("/:slug", TopicView, authTopicView)
	s.Get("/:slug/feed", TopicFeed, authTopicFeed)

	return s
}
