package middleware

import (
	"github.com/labstack/echo"
	Ctx "github.com/o-ga09/ecs-express-mode-api/pkg/context"
	"github.com/o-ga09/ecs-express-mode-api/pkg/uuid"
)

const (
	RequestIDHeader = "X-Request-ID"
)

// RequestID はリクエストごとにユニークなIDを生成するミドルウェア
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			c.SetRequest(c.Request().WithContext(Ctx.SetRequestID(ctx)))
			req := c.Request()
			res := c.Response()

			// 既存のリクエストIDを取得、なければ新規生成
			requestID := req.Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = uuid.GenerateID()
			}

			// リクエストとレスポンスにIDを設定
			req.Header.Set(RequestIDHeader, requestID)
			res.Header().Set(RequestIDHeader, requestID)

			return next(c)
		}
	}
}
