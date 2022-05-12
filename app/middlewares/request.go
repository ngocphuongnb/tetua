package middlewares

import (
	"time"

	"github.com/google/uuid"
	"github.com/ngocphuongnb/tetua/app/logger"
	"github.com/ngocphuongnb/tetua/app/server"
)

const HEADER_REQUEST_ID = "X-Request-Id"

func RequestID(c server.Context) error {
	requestId := c.Header(HEADER_REQUEST_ID)

	if requestId == "" {
		requestId = uuid.NewString()
	}

	c.Locals("request_id", requestId)
	c.Header(HEADER_REQUEST_ID, requestId)

	return c.Next()
}

func RequestLog(c server.Context) error {
	start := time.Now()
	err := c.Next()
	latency := time.Since(start).Round(time.Millisecond)
	logContext := logger.Context{
		"latency": latency.String(),
		"status":  c.Response().StatusCode(),
		"method":  c.Method(),
		"path":    c.Path(),
		"ip":      c.IP(),
	}

	if err != nil {
		logContext["error"] = err.Error()
	}

	c.Logger().Info("Request completed", logContext)
	return nil
}
