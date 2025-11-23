package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"cloud.google.com/go/logging"
	"github.com/o-ga09/ecs-express-mode-api/pkg/config"
	"go.opentelemetry.io/otel/trace"
)

var privateHeaders = []string{"Authorization", "Cookie", "X-CSRF-Token"}

// cloud logging の Log level 定義
var (
	Severitydefault = slog.Level(logging.Default)
	SeverityInfo    = slog.Level(logging.Info)
	SeverityWarn    = slog.Level(logging.Warning)
	SeverityError   = slog.Level(logging.Error)
	SeverityNotice  = slog.Level(logging.Notice)
)

// traceId , spanId 追加
type traceHandler struct {
	slog.Handler
	projectID string
	env       string
}

// RequestInfo はリクエスト情報を表す構造体
type RequestInfo struct {
	status                                            int
	contents_length                                   int64
	method, path, sourceIP, query, user_agent, errors string
	elapsed                                           time.Duration
}

// traceHandler 実装
func (h *traceHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.Handler.Enabled(ctx, l)
}

func (h *traceHandler) Handle(ctx context.Context, r slog.Record) error {
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		trace := fmt.Sprintf("projects/%s/traces/%s", h.projectID, sc.TraceID().String())
		r.AddAttrs(slog.String("logging.googleapis.com/trace", trace),
			slog.String("logging.googleapis.com/spanId", sc.SpanID().String()))
	}

	return h.Handler.Handle(ctx, r)
}

func (h *traceHandler) WithAttr(attrs []slog.Attr) slog.Handler {
	return &traceHandler{h.Handler.WithAttrs(attrs), h.projectID, h.env}
}

func (h *traceHandler) WithGroup(g string) slog.Handler {
	return h.Handler.WithGroup(g)
}

// logger 生成関数
func Logger(ctx context.Context) {
	replacer := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.MessageKey {
			a.Key = "message"
		}

		if a.Key == slog.LevelKey {
			a.Key = "severity"
			a.Value = slog.StringValue(logging.Severity(a.Value.Any().(slog.Level)).String())
		}

		if a.Key == slog.SourceKey {
			a.Key = "logging.googleapis.com/sourceLocation"
		}

		return a
	}
	env := config.GetCtxEnv(ctx)
	h := traceHandler{slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, ReplaceAttr: replacer}), env.ProjectID, env.Env}

	logger := slog.New(&h)
	slog.SetDefault(logger)
}
