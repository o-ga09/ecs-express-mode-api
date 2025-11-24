package middleware

import (
	"time"

	"github.com/labstack/echo"
	"github.com/o-ga09/ecs-express-mode-api/pkg/logger"
)

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			res := c.Response()
			req := c.Request()
			start := time.Now()

			logger.Info(
				req.Context(),
				"Request started",
				"method", req.Method,
				"path", req.URL.Path,
				"remote_addr", c.RealIP(),
				"user_agent", req.UserAgent(),
			)

			err := next(c)

			elapsed := time.Since(start)
			logger.Info(
				req.Context(),
				"Request completed",
				"method", req.Method,
				"path", req.URL.Path,
				"status", res.Status,
				"elapsed_ms", elapsed.Milliseconds(),
				"bytes_out", res.Size,
			)
			return err
		}
	}
}
