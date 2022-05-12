package fiberserver

import "github.com/gofiber/fiber/v2"

type Route struct {
	*fiber.Route
}

func (r *Route) Method() string {
	return r.Route.Method
}

func (r *Route) Path() string {
	return r.Route.Path
}

func (r *Route) Name() string {
	return r.Route.Name
}
