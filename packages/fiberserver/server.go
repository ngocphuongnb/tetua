package fiberserver

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/ngocphuongnb/tetua/app/logger"
	"github.com/ngocphuongnb/tetua/app/server"
)

type Server struct {
	*fiber.App
	middlewares []server.Handler
}

type Config struct {
	JwtSigningKey string
	AppName       string
}

func New(config Config) server.Server {
	app := fiber.New(fiber.Config{
		AppName:               config.AppName,
		StrictRouting:         true,
		CaseSensitive:         true,
		EnablePrintRoutes:     false,
		DisableStartupMessage: false,
		// Prefork:       true,
		// Immutable:     true,
	})

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	return &Server{App: app}
}

func (s *Server) Test(req *http.Request, msTimeout ...int) (resp *http.Response, err error) {
	return s.App.Test(req, msTimeout...)
}

func (s *Server) Listen(address string) {
	logger.Fatal("Listen", logger.Context{"Error": s.App.Listen(address)})
}

func (s *Server) Register(register func(ss server.Server)) server.Server {
	register(s)

	return s
}

func (s *Server) Use(handlers ...server.Handler) {
	s.middlewares = append(s.middlewares, handlers...)
}

func (s *Server) UsePrefix(prefix string, handlers ...server.Handler) {
	var params []interface{} = []interface{}{prefix}

	for _, handler := range transformHandlers(handlers) {
		params = append(params, handler)
	}

	s.App.Use(params...)
}

func (s *Server) Group(prefix string, handlers ...server.Handler) server.Group {
	var fiberHandlers []fiber.Handler

	for _, handler := range handlers {
		fiberHandlers = append(fiberHandlers, func(c *fiber.Ctx) error {
			return handler(createContext(c))
		})
	}

	g := s.App.Group(prefix, fiberHandlers...).(*fiber.Group)

	return &Group{
		FGroup:      g,
		App:         s.App,
		middlewares: s.middlewares,
	}
}

func (s *Server) Static(prefix, root string, configs ...server.StaticConfig) {
	config := fiber.Static{}

	if len(configs) > 0 {
		config = fiber.Static{
			Index:         configs[0].Index,
			Browse:        configs[0].Browse,
			MaxAge:        configs[0].MaxAge,
			Compress:      configs[0].Compress,
			ByteRange:     configs[0].ByteRange,
			Download:      configs[0].Download,
			CacheDuration: configs[0].CacheDuration,
		}
	}
	s.App.Static(prefix, root, config)
}

func (s *Server) Get(path string, handler server.Handler, authConfigs ...*server.AuthConfig) {
	authConfigs = append(authConfigs, &server.AuthConfig{Action: ""})
	s.App.Get(path, transformHandlers([]server.Handler{handler}, s.middlewares...)...).Name(authConfigs[0].Action)
}

func (s *Server) Post(path string, handler server.Handler, authConfigs ...*server.AuthConfig) {
	authConfigs = append(authConfigs, &server.AuthConfig{Action: ""})
	s.App.Post(path, transformHandlers([]server.Handler{handler}, s.middlewares...)...).Name(authConfigs[0].Action)
}

func (s *Server) Put(path string, handler server.Handler, authConfigs ...*server.AuthConfig) {
	authConfigs = append(authConfigs, &server.AuthConfig{Action: ""})
	s.App.Put(path, transformHandlers([]server.Handler{handler}, s.middlewares...)...).Name(authConfigs[0].Action)
}

func (s *Server) Delete(path string, handler server.Handler, authConfigs ...*server.AuthConfig) {
	authConfigs = append(authConfigs, &server.AuthConfig{Action: ""})
	s.App.Delete(path, transformHandlers([]server.Handler{handler}, s.middlewares...)...).Name(authConfigs[0].Action)
}

func (s *Server) Patch(path string, handler server.Handler, authConfigs ...*server.AuthConfig) {
	authConfigs = append(authConfigs, &server.AuthConfig{Action: ""})
	s.App.Patch(path, transformHandlers([]server.Handler{handler}, s.middlewares...)...).Name(authConfigs[0].Action)
}

func (s *Server) Options(path string, handler server.Handler, authConfigs ...*server.AuthConfig) {
	authConfigs = append(authConfigs, &server.AuthConfig{Action: ""})
	s.App.Options(path, transformHandlers([]server.Handler{handler}, s.middlewares...)...).Name(authConfigs[0].Action)
}

func (s *Server) Head(path string, handler server.Handler, authConfigs ...*server.AuthConfig) {
	authConfigs = append(authConfigs, &server.AuthConfig{Action: ""})
	s.App.Options(path, transformHandlers([]server.Handler{handler}, s.middlewares...)...).Name(authConfigs[0].Action)
}
