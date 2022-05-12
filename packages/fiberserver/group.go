package fiberserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngocphuongnb/tetua/app/server"
)

type Group struct {
	*fiber.App
	FGroup      *fiber.Group
	middlewares []server.Handler
}

func (g *Group) Use(handler server.Handler) {
	g.middlewares = append(g.middlewares, handler)
}

func (g *Group) Group(prefix string, handlers ...server.Handler) server.Group {
	var fiberHandlers []fiber.Handler

	for _, handler := range handlers {
		fiberHandlers = append(fiberHandlers, func(c *fiber.Ctx) error {
			return handler(createContext(c))
		})
	}

	gg := g.FGroup.Group(prefix, fiberHandlers...).(*fiber.Group)

	return &Group{
		FGroup:      gg,
		App:         g.App,
		middlewares: g.middlewares,
	}
}

func (g *Group) Get(path string, handler server.Handler, authConfigs ...*server.AuthConfig) {
	authConfigs = append(authConfigs, &server.AuthConfig{Action: ""})
	g.FGroup.Get(path, transformHandlers([]server.Handler{handler}, g.middlewares...)...).Name(authConfigs[0].Action)
}

func (g *Group) Post(path string, handler server.Handler, authConfigs ...*server.AuthConfig) {
	authConfigs = append(authConfigs, &server.AuthConfig{Action: ""})
	g.FGroup.Post(path, transformHandlers([]server.Handler{handler}, g.middlewares...)...).Name(authConfigs[0].Action)
}

func (g *Group) Put(path string, handler server.Handler, authConfigs ...*server.AuthConfig) {
	authConfigs = append(authConfigs, &server.AuthConfig{Action: ""})
	g.FGroup.Put(path, transformHandlers([]server.Handler{handler}, g.middlewares...)...).Name(authConfigs[0].Action)
}

func (g *Group) Delete(path string, handler server.Handler, authConfigs ...*server.AuthConfig) {
	authConfigs = append(authConfigs, &server.AuthConfig{Action: ""})
	g.FGroup.Delete(path, transformHandlers([]server.Handler{handler}, g.middlewares...)...).Name(authConfigs[0].Action)
}

func (g *Group) Patch(path string, handler server.Handler, authConfigs ...*server.AuthConfig) {
	authConfigs = append(authConfigs, &server.AuthConfig{Action: ""})
	g.FGroup.Patch(path, transformHandlers([]server.Handler{handler}, g.middlewares...)...).Name(authConfigs[0].Action)
}

func (g *Group) Options(path string, handler server.Handler, authConfigs ...*server.AuthConfig) {
	authConfigs = append(authConfigs, &server.AuthConfig{Action: ""})
	g.FGroup.Options(path, transformHandlers([]server.Handler{handler}, g.middlewares...)...).Name(authConfigs[0].Action)
}

func (g *Group) Head(path string, handler server.Handler, authConfigs ...*server.AuthConfig) {
	authConfigs = append(authConfigs, &server.AuthConfig{Action: ""})
	g.FGroup.Head(path, transformHandlers([]server.Handler{handler}, g.middlewares...)...).Name(authConfigs[0].Action)
}
