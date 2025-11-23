package route

import "github.com/labstack/echo"

func SetUpRouters(e *echo.Echo) {
	root := e.Group("/v1/api")
	{
		root.GET("/health", healthCheckHandler)
	}
}

func healthCheckHandler(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"status": "ok",
	})
}
