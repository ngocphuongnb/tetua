package fiberserver

import (
	"bufio"
	"context"
	"fmt"
	"mime/multipart"
	"runtime"
	"strconv"

	"github.com/gofiber/fiber/v2"
	fiberUtils "github.com/gofiber/utils"
	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/entities"
	"github.com/ngocphuongnb/tetua/app/logger"
	"github.com/ngocphuongnb/tetua/app/server"
	"github.com/ngocphuongnb/tetua/app/utils"
	"github.com/valyala/fasthttp"
)

type Context struct {
	*fiber.Ctx
	MetaData *entities.Meta
}

type Response struct {
	*fasthttp.Response
}

func (c *Context) RequestID() string {
	return fmt.Sprintf("%v", c.Locals("request_id"))
}

func (c *Context) Hostname() string {
	return c.Ctx.Hostname()
}

func (c *Context) BaseUrl() string {
	return c.Ctx.Protocol() + "://" + c.Ctx.Hostname()
}

func (c *Context) RouteName() string {
	return c.Ctx.Route().Name
}

func (c *Context) Context() context.Context {
	return context.WithValue(context.Background(), "request_id", c.RequestID())
}

func (c *Context) File(name string) (*multipart.FileHeader, error) {
	return c.Ctx.FormFile(name)
}

func (c *Context) Logger() logger.Logger {
	return logger.Get().WithContext(logger.Context{"request_id": c.RequestID()})
}

func (c *Context) Messages(ms ...*entities.Messages) *entities.Messages {
	if len(ms) > 0 {
		c.MetaData.Messages = ms[0]
		return ms[0]
	}
	return c.MetaData.Messages
}

func (c *Context) Meta(metas ...*entities.Meta) *entities.Meta {
	if len(metas) > 0 {
		c.MetaData = metas[0]
		c.MetaData.User = c.User()

		if c.MetaData.Messages == nil {
			c.MetaData.Messages = &entities.Messages{}
		}
		return c.MetaData
	}

	if c.MetaData == nil {
		appName := config.Setting("app_name")
		c.MetaData = &entities.Meta{
			Title:       appName,
			Description: appName,
			User:        c.User(),
			Messages:    &entities.Messages{},
		}
	}

	return c.MetaData
}

func (c *Context) Render(fn func(meta *entities.Meta, wr *bufio.Writer)) error {
	if c.Meta().Canonical == "" {
		c.Meta().Canonical = utils.Url(c.Path())
	}

	if c.Meta().Type == "" {
		c.Meta().Type = "website"
	}

	c.Response().Header("content-type", "text/html; charset=utf-8")
	requestID := fiberUtils.ImmutableString(c.RequestID())
	c.Ctx.Response().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer func(requestID string) {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				stack := make([]byte, 4<<10)
				length := runtime.Stack(stack, true)
				msg := fmt.Sprintf("%v %s\n", err, stack[:length])
				logger.Get().WithContext(logger.Context{"request_id": requestID, "recovered": true}).Error(msg)
				if _, err := w.WriteString("Something went wrong"); err != nil {
					fmt.Println(err)
				}
			}
		}(requestID)
		fn(c.Meta(), w)
		w.Flush()
	})

	return nil
}

func (c *Context) User() *entities.User {
	if user, ok := c.Locals("user").(*entities.User); ok {
		return user
	}

	return nil
}

func (c *Context) Post(posts ...*entities.Post) *entities.Post {
	if len(posts) > 0 {
		c.Locals("post", posts[0])
	}

	if post, ok := c.Locals("post").(*entities.Post); ok {
		return post
	}

	return nil
}

func (c *Context) Json(v interface{}) error {
	return c.Ctx.JSON(v)
}

func (c *Context) Cookies(key string, defaultValue ...string) string {
	return c.Ctx.Cookies(key, defaultValue...)
}

func (c *Context) Cookie(v *server.Cookie) {
	cookie := fiber.Cookie{
		Name:     v.Name,
		Value:    v.Value,
		Path:     v.Path,
		Domain:   v.Domain,
		Expires:  v.Expires,
		Secure:   v.Secure,
		HTTPOnly: v.HTTPOnly,
		SameSite: v.SameSite,
	}
	c.Ctx.Cookie(&cookie)
}

func (c *Context) Status(v int) server.Context {
	c.Ctx.Status(v)
	return c
}
func (c *Context) Next() error {
	return c.Ctx.Next()
}

func (c *Context) Locals(key string, value ...interface{}) (val interface{}) {
	return c.Ctx.Locals(key, value...)
}

func (c *Context) Param(key string) string {
	return c.Ctx.Params(key)
}

func (c *Context) ParamInt(key string, values ...int) int {
	if p, err := strconv.Atoi(c.Ctx.Params(key)); err == nil {
		return p
	}

	if len(values) > 0 {
		return values[0]
	}

	return 0
}

func (c *Context) WithError(msg string, err error) {
	c.Messages().AppendError(msg)
	c.Logger().Error(msg, err)
}

func (c *Context) Query(key string, values ...string) string {
	if value := c.Ctx.Query(key); value != "" {
		return value
	}

	if len(values) > 0 {
		return values[0]
	}

	return ""
}

func (c *Context) QueryInt(key string, values ...int) int {
	if p, err := strconv.Atoi(c.Ctx.Query(key)); err == nil {
		return p
	}

	if len(values) > 0 {
		return values[0]
	}

	return 0
}

func (c *Context) Response() server.Response {
	return &Response{c.Ctx.Response()}
}

func (c *Context) Method() string {
	return c.Ctx.Method()
}

func (c *Context) Send(data []byte) error {
	return c.Ctx.Send(data)
}

func (c *Context) SendString(data string) error {
	return c.Ctx.SendString(data)
}

func (c *Context) OriginalURL() string {
	return c.Ctx.OriginalURL()
}

func (c *Context) Path() string {
	return c.Ctx.Path()
}

func (c *Context) Redirect(path string) error {
	return c.Ctx.Redirect(path)
}

func (c *Context) RedirectToRoute(name string, params ...map[string]interface{}) error {
	if len(params) > 0 {
		return c.Ctx.RedirectToRoute(name, params[0])
	}
	return c.Ctx.RedirectToRoute(name, fiber.Map{})
}

func (c *Context) Header(key string, vals ...string) string {

	if len(vals) > 0 {
		c.Ctx.Set(key, vals[0])
		return vals[0]
	}

	return c.Ctx.Get(key)
}

func (c *Context) BodyParser(v interface{}) error {
	return c.Ctx.BodyParser(v)
}

func (r *Response) Header(key string, vals ...string) string {
	if len(vals) > 0 {
		r.Response.Header.Add(key, vals[0])
		return vals[0]
	}

	return string(r.Response.Header.Peek(key))
}
