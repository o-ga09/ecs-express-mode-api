package middleware

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// CORSConfig はCORS設定を返すカスタムミドルウェア
func CORSConfig() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.PATCH,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			RequestIDHeader,
		},
		ExposeHeaders: []string{
			RequestIDHeader,
		},
		AllowCredentials: true,
		MaxAge:           3600,
	})
}
