package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/morikuni/failure/v2"
	"github.com/o-ga09/ecs-express-mode-api/pkg/errors"
	"github.com/o-ga09/ecs-express-mode-api/pkg/logger"
)

// ErrorResponse はエラーレスポンスの構造体
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

// CustomErrorHandler はカスタムエラーハンドラー
func CustomErrorHandler(err error, c echo.Context) {
	ctx := c.Request().Context()
	var code int
	var message string
	var errorCode string
	var details string

	stack := getCallstack(err)
	errMessage := errors.GetMessage(err)
	// Echo のHTTPErrorの場合
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message.(string)
		errorCode = string(errors.GetCode(err))
	} else {
		// カスタムエラーの場合
		code = http.StatusInternalServerError
		message = errMessage
		errorCode = string(errors.GetCode(err))

		// エラーログを記録
		logger.Error(
			ctx,
			errMessage,
			"callStack", stack,
			"path", c.Request().URL.Path,
			"method", c.Request().Method,
		)
	}

	// エラーレスポンスを返す
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, ErrorResponse{
				Error:   message,
				Code:    errorCode,
				Details: details,
			})
		}
		if err != nil {
			logger.Error(ctx, "Failed to send error response", "error", err.Error())
		}
	}
}

func getCallstack(err error) string {
	if err == nil {
		return ""
	}
	callstack := failure.CallStackOf(err)
	msg := ""
	for _, frame := range callstack.Frames() {
		msg += fmt.Sprintf("%+v\n", frame)
	}
	return msg
}
