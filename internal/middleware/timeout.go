package middleware

import (
	"context"
	"time"

	"github.com/labstack/echo"
)

// TimeoutConfig はタイムアウト設定を返すカスタムミドルウェア
func TimeoutConfig() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(c.Request().Context(), 30*time.Second)
			defer cancel()

			// リクエストのコンテキストを更新
			c.SetRequest(c.Request().WithContext(ctx))

			done := make(chan error, 1)
			go func() {
				done <- next(c)
			}()

			select {
			case err := <-done:
				return err
			case <-ctx.Done():
				return echo.NewHTTPError(408, "Request Timeout")
			}
		}
	}
}
