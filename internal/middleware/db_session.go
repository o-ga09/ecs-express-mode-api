package middleware

import (
	"context"

	"github.com/labstack/echo"
	"github.com/o-ga09/ecs-express-mode-api/internal/database/mysql"
	"gorm.io/gorm"
)

func SetDB() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := mysql.New(c.Request().Context())
			db := ctx.Value(mysql.CtxKey).(*gorm.DB)
			ctx = context.WithValue(ctx, mysql.CtxKey, db)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
