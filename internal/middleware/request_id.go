package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

const (
	RequestIDHeader = "X-Request-ID"
)

// RequestID はリクエストごとにユニークなIDを生成するミドルウェア
func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			// 既存のリクエストIDを取得、なければ新規生成
			requestID := req.Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// リクエストとレスポンスにIDを設定
			req.Header.Set(RequestIDHeader, requestID)
			res.Header().Set(RequestIDHeader, requestID)

			return next(c)
		}
	}
}
