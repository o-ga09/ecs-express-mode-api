package route

import (
	"github.com/labstack/echo"
	"github.com/o-ga09/ecs-express-mode-api/internal/database/mysql"
	"github.com/o-ga09/ecs-express-mode-api/pkg/errors"
)

func SetUpRouters(e *echo.Echo) {
	root := e.Group("/v1/api")
	{
		root.GET("/health", healthCheckHandler)
		root.GET("/health/db", dbHealthCheckHandler)
	}
}

func healthCheckHandler(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"status": "ok",
	})
}

func dbHealthCheckHandler(c echo.Context) error {
	ctx := c.Request().Context()
	db := mysql.CtxFromDB(ctx)
	err := db.Select("1").Error
	if err != nil {
		errors.Wrap(ctx, err)
	}
	return c.JSON(200, map[string]string{
		"db_status": "connected",
	})
}
