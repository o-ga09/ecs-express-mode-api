package logger

import (
	"context"
	"log/slog"

	"github.com/o-ga09/ecs-express-mode-api/pkg/constant"
	Ctx "github.com/o-ga09/ecs-express-mode-api/pkg/context"
)

func Info(ctx context.Context, msg string, args ...any) {
	allArgs := append([]any{"requestId", Ctx.GetRequestID(ctx)}, args...)
	slog.Log(ctx, constant.SeverityInfo, msg, allArgs...)
}

func Error(ctx context.Context, msg string, args ...any) {
	allArgs := append([]any{"requestId", Ctx.GetRequestID(ctx)}, args...)
	slog.Log(ctx, constant.SeverityError, msg, allArgs...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	allArgs := append([]any{"requestId", Ctx.GetRequestID(ctx)}, args...)
	slog.Log(ctx, constant.SeverityWarn, msg, allArgs...)
}

func Notice(ctx context.Context, msg string, args ...any) {
	allArgs := append([]any{"requestId", Ctx.GetRequestID(ctx)}, args...)
	slog.Log(ctx, constant.SeverityNotice, msg, allArgs...)
}
