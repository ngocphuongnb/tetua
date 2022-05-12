package fiberserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/server"
)

func transformHandlers(handlers []server.Handler, middlewares ...server.Handler) []fiber.Handler {
	var fiberHandlers []fiber.Handler

	for i := 0; i < len(middlewares); i++ {
		func(i int) {
			fiberHandlers = append(fiberHandlers, func(c *fiber.Ctx) error {
				return middlewares[i](createContext(c))
			})
		}(i)
	}

	for i := 0; i < len(handlers); i++ {
		func(i int) {
			fiberHandlers = append(fiberHandlers, func(c *fiber.Ctx) error {
				return handlers[i](createContext(c))
			})
		}(i)
	}

	return fiberHandlers
}

func createContext(c *fiber.Ctx) server.Context {
	var user *entities.User

	if c.Locals("user") != nil {
		user = c.Locals("user").(*entities.User)
	}

	return &Context{
		Ctx: c,
		MetaData: &entities.Meta{
			Messages: &entities.Messages{},
			User:     user,
		},
	}
}
