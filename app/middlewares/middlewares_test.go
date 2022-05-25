package middlewares_test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/ngocphuongnb/tetua/app/config"
	"github.com/ngocphuongnb/tetua/app/logger"
	"github.com/ngocphuongnb/tetua/app/middlewares"
	"github.com/ngocphuongnb/tetua/app/mock"
	"github.com/ngocphuongnb/tetua/app/server"
	fiber "github.com/ngocphuongnb/tetua/packages/fiberserver"
	"github.com/stretchr/testify/assert"
)

func TestCookieMiddleware(t *testing.T) {
	mockServer := mock.CreateServer()
	mockServer.Use(middlewares.Cookie)
	mockServer.Get("/test", func(c server.Context) error {
		return c.SendString("ok")
	})
	body, resp := mock.GetRequest(mockServer, "/test")
	assert.Equal(t, "ok", string(body))
	assert.Equal(t, true, strings.HasPrefix(resp.Header["Set-Cookie"][0], config.COOKIE_UUID+"="))
}

func TestRecoverMiddleware(t *testing.T) {
	mockLogger := mock.CreateLogger(true)
	mockServer := mock.CreateServer()
	mockServer.Use(middlewares.Recover)
	mockServer.Get("/test", func(c server.Context) error {
		panic("test recover")
	})
	body, _ := mock.GetRequest(mockServer, "/test")
	assert.Equal(t, `{"error":"test recover"}`, body)
	param0Str := fmt.Sprintf("%v", mockLogger.Last().Params[0])
	assert.Equal(t, true, strings.HasPrefix(param0Str, "test recover goroutine"))
}

func TestRecoverMiddlewareResponseError(t *testing.T) {
	mockLogger := mock.CreateLogger(true)
	mockServer := fiber.New(fiber.Config{
		JwtSigningKey: "sesj5JYrRxrB2yUWkBFM7KKWCY2ykxBw",
		JSONEncoder: func(v interface{}) ([]byte, error) {
			return nil, errors.New("json error")
		},
	})

	mockServer.Use(middlewares.Recover)
	mockServer.Get("/test", func(c server.Context) error {
		panic("test recover")
	})
	mock.GetRequest(mockServer, "/test")
	assert.Equal(t, errors.New("json error"), mockLogger.Last().Params[0])
}

func TestRequestIDMiddlewares(t *testing.T) {
	mockServer := mock.CreateServer()
	mockServer.Use(middlewares.RequestID)
	mockServer.Get("/test", func(c server.Context) error {
		assert.Equal(t, true, c.RequestID() != "")
		return c.SendString("ok")
	})
	mockServer.Get("/test2", func(c server.Context) error {
		assert.Equal(t, "my-request-id", c.RequestID())
		return c.SendString("ok")
	})

	_, resp := mock.GetRequest(mockServer, "/test")
	assert.Equal(t, true, resp.Header["X-Request-Id"][0] != "")

	_, resp = mock.GetRequest(mockServer, "/test2", map[string]string{
		middlewares.HEADER_REQUEST_ID: "my-request-id",
	})
	assert.Equal(t, "my-request-id", resp.Header["X-Request-Id"][0])
}

func TestRequestLogMiddlewares(t *testing.T) {
	mockLogger := mock.CreateLogger(true)
	mockServer := mock.CreateServer()
	mockServer.Use(middlewares.RequestLog)
	mockServer.Get("/test", func(c server.Context) error {
		c.Status(http.StatusBadRequest)
		return errors.New("Test request error")
	})

	mock.GetRequest(mockServer, "/test")

	msg := mockLogger.Last().Params[0]
	ctx, ok := mockLogger.Last().Params[1].(logger.Context)

	assert.Equal(t, true, ok)
	assert.Equal(t, "Request completed", msg)
	assert.Equal(t, http.StatusBadRequest, ctx["status"])
	assert.Equal(t, "GET", ctx["method"])
	assert.Equal(t, "/test", ctx["path"])
	assert.Equal(t, "0.0.0.0", ctx["ip"])
	assert.Equal(t, "Test request error", ctx["error"])
	assert.Equal(t, true, ctx["latency"] != "")
}

func TestGetAllMiddlewares(t *testing.T) {
	assert.Equal(t, 6, len(middlewares.All()))
}
