package auth_test

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ngocphuongnb/tetua/app/auth"
	"github.com/ngocphuongnb/tetua/app/cache"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/mock"
	mockrepository "github.com/ngocphuongnb/tetua/app/mock/repository"
	"github.com/ngocphuongnb/tetua/app/repositories"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/test"
	ga "github.com/ngocphuongnb/tetua/packages/auth"
	"github.com/ngocphuongnb/tetua/packages/fiberserver"
	fiber "github.com/ngocphuongnb/tetua/packages/fiberserver"
	"github.com/stretchr/testify/assert"
)

func init() {
	mock.CreateRepositories()
	repositories.User.Create(context.Background(), mock.RootUser)
	cache.Roles = []*entities.Role{auth.ROLE_ADMIN, auth.ROLE_USER, auth.ROLE_GUEST}
}

func TestProvider(t *testing.T) {
	providers := map[string]auth.NewProviderFn{
		"local": ga.NewLocal,
		"mock1": mock.NewAuth,
		"mock2": mock.NewAuth,
	}
	providerMap := map[string]map[string]string{
		"local": nil,
		"mock1": nil,
		"mock2": nil,
	}

	config.Auth = nil
	auth.New(providers)
	assert.Equal(t, 0, len(auth.Providers()))

	config.Auth = &config.AuthConfig{
		EnabledProviders: []string{"local"},
		Providers:        providerMap,
	}
	auth.New(providers)
	assert.Equal(t, 1, len(auth.Providers()))

	config.Auth = &config.AuthConfig{
		EnabledProviders: []string{"local", "mock1", "mock2"},
		Providers:        providerMap,
	}
	auth.New(providers)
	assert.Equal(t, 2, len(auth.Providers()))
	assert.Equal(t, providers["mock1"](nil).Name(), auth.GetProvider("mock").Name())
	assert.Equal(t, nil, auth.GetProvider("invalid_provider"))
}

func TestActionConfigs(t *testing.T) {
	testConfig := auth.Config(&server.AuthConfig{
		Action:       "test",
		Value:        entities.PERM_ALL,
		DefaultValue: entities.PERM_ALL,
	})

	assert.Equal(t, testConfig, auth.GetAuthConfig("test"))
	assert.Equal(t, (*server.AuthConfig)(nil), auth.GetAuthConfig("test2"))

	defer test.RecoverPanic(t, "Duplicate action config: test", "duplicate action")
	auth.Config(&server.AuthConfig{
		Action:       "test",
		Value:        entities.PERM_ALL,
		DefaultValue: entities.PERM_ALL,
	})
}

func TestHelpers(t *testing.T) {
	ctx := &fiberserver.Context{}
	role1Permissions := []*entities.PermissionValue{{
		Action: "test",
		Value:  entities.PERM_ALL,
	}}
	cache.RolesPermissions = []*entities.RolePermissions{{
		RoleID:      1,
		Permissions: role1Permissions,
	}}

	assert.Equal(t, role1Permissions, auth.GetRolePermissions(1).Permissions)
	assert.Equal(t, []*entities.PermissionValue{}, auth.GetRolePermissions(2).Permissions)
	assert.Equal(t, role1Permissions[0], auth.GetRolePermission(1, "test"))
	assert.Equal(t, &entities.PermissionValue{}, auth.GetRolePermission(1, "test2"))
	assert.Equal(t, []*entities.Role{cache.Roles[0]}, auth.GetRolesFromIDs([]int{1}))
	assert.Equal(t, []*entities.Role{}, auth.GetRolesFromIDs([]int{4}))
	assert.Equal(t, true, auth.AllowAll(ctx))
	assert.Equal(t, false, auth.AllowNone(ctx))

	s := mock.CreateServer()

	s.Get("/test", func(c server.Context) error {
		assert.Equal(t, false, auth.AllowLoggedInUser(c))
		c.Locals("user", &entities.User{ID: 1})
		assert.Equal(t, true, auth.AllowLoggedInUser(c))
		err := auth.GetFile(c)
		assert.Equal(t, true, entities.IsNotFound(err))
		return nil
	})

	mock.GetRequest(s, "/test")
	repositories.File.Create(context.Background(), &entities.File{
		ID:     1,
		UserID: 1,
	})

	s.Get("/files/:id", func(c server.Context) error {
		err := auth.GetFile(c)

		if c.Param("id") == "new" {
			assert.Equal(t, true, auth.FileOwnerCheck(c))
		}

		if c.ParamInt("id") > 1 {
			assert.Equal(t, true, entities.IsNotFound(err))
			assert.Equal(t, nil, c.Locals("file"))
		} else if c.ParamInt("id") == 1 {
			assert.Equal(t, nil, err)

			if f, ok := c.Locals("file").(*entities.File); ok {
				assert.Equal(t, 1, f.ID)
			} else {
				assert.Fail(t, "local file is not a file")
			}

			c.Locals("user", nil)
			assert.Equal(t, false, auth.FileOwnerCheck(c))

			c.Locals("user", &entities.User{ID: 2})
			assert.Equal(t, false, auth.FileOwnerCheck(c))

			c.Locals("user", &entities.User{ID: 1})
			assert.Equal(t, true, auth.FileOwnerCheck(c))

			c.Locals("file", 1)
			assert.Equal(t, false, auth.FileOwnerCheck(c))
		}

		return nil
	})

	mock.GetRequest(s, "/files/new")
	mock.GetRequest(s, "/files/1")
	mock.GetRequest(s, "/files/2")

	repositories.Post.Create(context.Background(), &entities.Post{
		ID:     1,
		Name:   "post 1",
		UserID: 1,
	})
	repositories.Post.Create(context.Background(), &entities.Post{
		ID:     2,
		Name:   "post 2",
		UserID: 2,
	})

	s.Get("/posts/:id", func(c server.Context) error {
		err := auth.GetPost(c)

		if c.Param("id") == "new" {
			assert.Equal(t, nil, err)
			assert.Equal(t, nil, c.Locals("post"))
			assert.Equal(t, true, auth.PostOwnerCheck(c))
		} else if c.ParamInt("id") == 1 {
			assert.Equal(t, nil, err)

			if p, ok := c.Locals("post").(*entities.Post); ok {
				assert.Equal(t, 1, p.ID)
			} else {
				assert.Fail(t, "local post is not a post")
			}

			c.Locals("user", &entities.User{ID: 1})
			assert.Equal(t, true, auth.PostOwnerCheck(c))

			c.Locals("user", &entities.User{ID: 2})
			assert.Equal(t, false, auth.PostOwnerCheck(c))

			c.Locals("user", nil)
			assert.Equal(t, false, auth.PostOwnerCheck(c))
		} else if c.ParamInt("id") > 2 {
			assert.Equal(t, true, entities.IsNotFound(err))
			assert.Equal(t, nil, c.Locals("post"))
		}

		return nil
	})

	mock.GetRequest(s, "/posts/new")
	mock.GetRequest(s, "/posts/1")
	mock.GetRequest(s, "/posts/2")
	mock.GetRequest(s, "/posts/3")

	repositories.Comment.Create(context.Background(), &entities.Comment{
		ID:     1,
		UserID: 1,
	})

	s.Get("/comments/:id", func(c server.Context) error {
		err := auth.GetComment(c)

		if c.Param("id") == "new" {
			assert.Equal(t, nil, err)
			assert.Equal(t, nil, c.Locals("comment"))
			assert.Equal(t, true, auth.CommentOwnerCheck(c))
		} else if c.ParamInt("id") == 1 {
			assert.Equal(t, nil, err)

			if c, ok := c.Locals("comment").(*entities.Comment); ok {
				assert.Equal(t, 1, c.ID)
			} else {
				assert.Fail(t, "local comment is not a comment")
			}

			c.Locals("user", &entities.User{ID: 1})
			assert.Equal(t, true, auth.CommentOwnerCheck(c))

			c.Locals("user", &entities.User{ID: 2})
			assert.Equal(t, false, auth.CommentOwnerCheck(c))

			c.Locals("user", nil)
			assert.Equal(t, false, auth.CommentOwnerCheck(c))

			c.Locals("comment", nil)
			assert.Equal(t, false, auth.CommentOwnerCheck(c))

		} else {
			assert.Equal(t, true, entities.IsNotFound(err))
			assert.Equal(t, nil, c.Locals("comment"))
		}

		return nil
	})

	mock.GetRequest(s, "/comments/new")
	mock.GetRequest(s, "/comments/0")
	mock.GetRequest(s, "/comments/1")
	mock.GetRequest(s, "/comments/2")
}

func TestRoutes(t *testing.T) {
	logger := mock.CreateLogger()
	s := mock.CreateServer()
	auth.Routes(s)

	body, resp := mock.GetRequest(s, "/auth/invalid_provider")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "Invalid provider", string(body))

	body, resp = mock.GetRequest(s, "/auth/mock")

	assert.Equal(t, true, strings.HasPrefix(resp.Header["Location"][0], "https://auth.mock-service.local/auth"))
	assert.Equal(t, http.StatusFound, resp.StatusCode)

	body, _ = mock.GetRequest(s, "/auth/mock/callback")
	assert.Equal(t, "Something went wrong", string(body))
	assert.Equal(t, []*mock.MockLoggerMessage{{
		Type:   "Error",
		Params: []interface{}{errors.New("code is empty")},
	}}, logger.Messages)

	_, resp = mock.GetRequest(s, "/auth/mock/callback?code=mock_callback_code")
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/", resp.Header["Location"][0])

	header := http.Header{}
	header.Add("Cookie", strings.Split(resp.Header["Set-Cookie"][0], ";")[0])
	request := http.Request{Header: header}
	cookies := request.Cookies()
	token, err := jwt.ParseWithClaims(
		cookies[0].Value,
		&entities.UserJwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.APP_KEY), nil
		},
	)

	assert.Equal(t, nil, err)
	claims, ok := token.Claims.(*entities.UserJwtClaims)
	assert.Equal(t, true, ok)
	assert.Equal(t, true, token.Valid)
	assert.Equal(t, mock.NormalUser2.ID, claims.User.ID)
	assert.Equal(t, mock.NormalUser2.Provider, claims.User.Provider)
	assert.Equal(t, mock.NormalUser2.Username, claims.User.Username)
	assert.Equal(t, mock.NormalUser2.Email, claims.User.Email)
	assert.Equal(t, mock.NormalUser2.DisplayName, claims.User.DisplayName)
	assert.Equal(t, mock.NormalUser2.RoleIDs, claims.User.RoleIDs)
	assert.Equal(t, mock.NormalUser2.AvatarImageUrl, claims.User.AvatarImageUrl)

	s2 := fiber.New(fiber.Config{JwtSigningKey: config.APP_KEY})
	s2.Use(func(c server.Context) error {
		c.Locals("jwt_header", map[string]interface{}{
			"invalid_header": math.Inf(1),
		})
		return c.Next()
	})
	auth.Routes(s2)

	_, resp = mock.GetRequest(s2, "/auth/mock/callback?code=mock_callback_code")
	assert.Equal(t, "Error setting login info", logger.Messages[1].Params[0])
	assert.Equal(t, http.StatusBadGateway, resp.StatusCode)

	mockrepository.ErrorCreateIfNotExistsByProvider = true

	body, _ = mock.GetRequest(s, "/auth/mock/callback?code=mock_callback_code")
	assert.Equal(t, "Something went wrong", string(body))
	assert.Equal(t, mock.MockLoggerMessage{
		Type:   "Error",
		Params: []interface{}{errors.New("CreateIfNotExistsByProvider error")},
	}, logger.Last())

	mockrepository.ErrorCreateIfNotExistsByProvider = false
}

func TestAssignUserInfo(t *testing.T) {
	logger := mock.CreateLogger()
	s := mock.CreateServer()
	s.Use(auth.AssignUserInfo)

	s.Get("/nocookie", func(c server.Context) error {
		assert.Equal(t, auth.GUEST_USER, c.User())
		return nil
	})
	mock.GetRequest(s, "/nocookie")

	s.Get("/invalidcookie", func(c server.Context) error {
		assert.Equal(t, auth.GUEST_USER, c.User())
		return nil
	})

	mock.GetRequest(s, "/invalidcookie", map[string]string{
		"cookie": config.APP_TOKEN_KEY + "=aaaa",
	})
	assert.Equal(t, jwt.NewValidationError("token contains an invalid number of segments", 0x1), logger.Messages[0].Params[0])

	s.Get("/validcookie", func(c server.Context) error {
		assert.Equal(t, &entities.User{
			ID:             mock.NormalUser2.ID,
			Provider:       mock.NormalUser2.Provider,
			Username:       mock.NormalUser2.Username,
			Email:          mock.NormalUser2.Email,
			DisplayName:    mock.NormalUser2.DisplayName,
			RoleIDs:        mock.NormalUser2.RoleIDs,
			Active:         mock.NormalUser2.Active,
			Roles:          auth.GetRolesFromIDs(mock.NormalUser2.RoleIDs),
			AvatarImageUrl: mock.NormalUser2.Avatar(),
		}, c.User())
		return nil
	})

	exp := time.Now().Add(time.Hour * 100 * 365 * 24)
	jwtToken, _ := mock.NormalUser2.JwtClaim(exp)
	mock.GetRequest(s, "/validcookie", map[string]string{
		"cookie": config.APP_TOKEN_KEY + "=" + jwtToken,
	})
	assert.Equal(t, jwt.NewValidationError("token contains an invalid number of segments", 0x1), logger.Messages[0].Params[0])
}

func TestAuthCheck(t *testing.T) {
	s := mock.CreateServer()
	s.Use(auth.Check)

	s.Get("/noauthconfig", func(c server.Context) error {
		return c.SendString("noauthconfig")
	})
	body, _ := mock.GetRequest(s, "/noauthconfig")
	assert.Equal(t, "noauthconfig", string(body))

	s.Get("/withauthconfigprepareerror", func(c server.Context) error {
		return c.SendString("withauthconfigprepareerror")
	}, auth.Config(&server.AuthConfig{
		Action:       "withauthconfigprepareerror",
		DefaultValue: entities.PERM_ALL,
		Prepare: func(c server.Context) error {
			return errors.New("prepare error")
		},
	}))

	exp := time.Now().Add(time.Hour * 100 * 365 * 24)
	jwtToken, _ := mock.NormalUser2.JwtClaim(exp)
	validAuthHeader := map[string]string{"cookie": config.APP_TOKEN_KEY + "=" + jwtToken}
	body, resp := mock.GetRequest(s, "/withauthconfigprepareerror", validAuthHeader)
	assert.Equal(t, "prepare error", body)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func createServerWithAuthConfig(action string) server.Server {
	s := mock.CreateServer()
	s.Get("/posts/:id", func(c server.Context) error {
		return c.SendString("View post: " + c.Param("id"))
	}, auth.Config(&server.AuthConfig{
		Action:       action,
		DefaultValue: entities.PERM_NONE,
		Prepare: func(c server.Context) error {
			return auth.GetPost(c)
		},
		OwnCheckFN: auth.PostOwnerCheck,
	}))
	return s
}

func TestGuestViewPostAllActionConfigs(t *testing.T) {
	var body = ""
	var s server.Server
	var resp *http.Response

	s = createServerWithAuthConfig("guest.post.view.perm_none")
	_, resp = mock.GetRequest(s, "/posts/3")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	_, resp = mock.GetRequest(s, "/posts/1")
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/login?back=%2Fposts%2F1", resp.Header["Location"][0])

	s = createServerWithAuthConfig("guest.post.view.perm_own")
	_, resp = mock.GetRequest(s, "/posts/1")
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/login?back=%2Fposts%2F1", resp.Header["Location"][0])

	cache.RolesPermissions = []*entities.RolePermissions{{
		RoleID: 3, // Guest
		Permissions: []*entities.PermissionValue{{
			Action: "guest.post.view.perm_all",
			Value:  entities.PERM_ALL,
		}},
	}}
	s = createServerWithAuthConfig("guest.post.view.perm_all")
	body, resp = mock.GetRequest(s, "/posts/1")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "View post: 1", body)
}

func TestInactiveUser(t *testing.T) {
	var s server.Server
	var resp *http.Response
	exp := time.Now().Add(time.Hour * 100 * 365 * 24)

	mock.NormalUser2.Active = false
	jwtTokenNormalUser2, _ := mock.NormalUser2.JwtClaim(exp)
	authHeaderNormalUser2 := map[string]string{"cookie": config.APP_TOKEN_KEY + "=" + jwtTokenNormalUser2}

	s = createServerWithAuthConfig("inactiveuser.post.view.perm_none")
	cache.RolesPermissions = []*entities.RolePermissions{{
		RoleID: 2, // User
		Permissions: []*entities.PermissionValue{{
			Action: "inactiveuser.post.view.perm_none",
			Value:  entities.PERM_NONE,
		}},
	}}

	_, resp = mock.GetRequest(s, "/posts/2", authHeaderNormalUser2)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/inactive", resp.Header["Location"][0])
	mock.NormalUser2.Active = true
}

func TestUserViewPostAllActionConfigs(t *testing.T) {
	var body = ""
	var s server.Server
	var resp *http.Response
	exp := time.Now().Add(time.Hour * 100 * 365 * 24)

	jwtTokenNormalUser2, _ := mock.NormalUser2.JwtClaim(exp)
	authHeaderNormalUser2 := map[string]string{"cookie": config.APP_TOKEN_KEY + "=" + jwtTokenNormalUser2}

	jwtTokenNormalUser3, _ := mock.NormalUser3.JwtClaim(exp)
	authHeaderNormalUser3 := map[string]string{"cookie": config.APP_TOKEN_KEY + "=" + jwtTokenNormalUser3}

	// Action that allow no one to access
	s = createServerWithAuthConfig("user.post.view.perm_none")
	cache.RolesPermissions = []*entities.RolePermissions{{
		RoleID: 2, // User
		Permissions: []*entities.PermissionValue{{
			Action: "user.post.view.perm_none",
			Value:  entities.PERM_NONE,
		}},
	}}

	// Access from user who is the post owner
	body, resp = mock.GetRequest(s, "/posts/2", authHeaderNormalUser2)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.Equal(t, "Insufficient permission", body)

	// Access from user who is not the post owner
	body, resp = mock.GetRequest(s, "/posts/2", authHeaderNormalUser3)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.Equal(t, "Insufficient permission", body)

	// Action that only allow the owner to access
	s = createServerWithAuthConfig("user.post.view.perm_own")
	cache.RolesPermissions = []*entities.RolePermissions{{
		RoleID: 2, // User
		Permissions: []*entities.PermissionValue{{
			Action: "user.post.view.perm_own",
			Value:  entities.PERM_OWN,
		}},
	}}

	// Access from user who is not the post owner
	body, resp = mock.GetRequest(s, "/posts/2", authHeaderNormalUser3)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.Equal(t, "Insufficient permission", body)

	// Access from user who is the post owner
	body, resp = mock.GetRequest(s, "/posts/2", authHeaderNormalUser2)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "View post: 2", body)

	// Action that allow all users to access
	s = createServerWithAuthConfig("user.post.view.perm_all")
	cache.RolesPermissions = []*entities.RolePermissions{{
		RoleID: 2, // User
		Permissions: []*entities.PermissionValue{{
			Action: "user.post.view.perm_all",
			Value:  entities.PERM_ALL,
		}},
	}}

	// Access from user who is not the post owner
	body, resp = mock.GetRequest(s, "/posts/2", authHeaderNormalUser3)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "View post: 2", body)

	// Access from user who is the post owner
	body, resp = mock.GetRequest(s, "/posts/2", authHeaderNormalUser2)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "View post: 2", body)
}

func TestRootUserAllActionConfigs(t *testing.T) {
	var body = ""
	var s server.Server
	var resp *http.Response
	exp := time.Now().Add(time.Hour * 100 * 365 * 24)

	jwtTokenRootUser, _ := mock.RootUser.JwtClaim(exp)
	authHeaderRootUser := map[string]string{"cookie": config.APP_TOKEN_KEY + "=" + jwtTokenRootUser}

	s = createServerWithAuthConfig("root.post.view.perm_none")

	body, resp = mock.GetRequest(s, "/posts/2", authHeaderRootUser)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "View post: 2", body)
}
