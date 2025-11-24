package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/o-ga09/ecs-express-mode-api/internal/middleware"
	"github.com/o-ga09/ecs-express-mode-api/internal/route"
	"github.com/o-ga09/ecs-express-mode-api/pkg/config"
	"github.com/o-ga09/ecs-express-mode-api/pkg/logger"
)

type Server struct {
	Config *config.Config
}

func NewServer(ctx context.Context) *Server {
	return &Server{
		Config: config.GetCtxEnv(ctx),
	}
}

func (s *Server) Run(ctx context.Context) error {
	e := echo.New()

	// カスタムエラーハンドラーを設定
	e.HTTPErrorHandler = middleware.CustomErrorHandler

	// ミドルウェア
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.SetDB())
	e.Use(middleware.CORSConfig())
	e.Use(middleware.TimeoutConfig())

	route.SetUpRouters(e)

	// サーバーの起動
	port := fmt.Sprintf(":%s", s.Config.Port)
	srv := &http.Server{
		Addr:    port,
		Handler: e,
	}

	// サーバーの起動
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(ctx, fmt.Sprintf("Failed to listen and serve: %v", err))
		}
	}()

	logger.Info(ctx, fmt.Sprintf("Server is running on %s", port))
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info(ctx, "graceful shutdown")

	// サーバーのタイムアウト設定
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// サーバーのシャットダウン
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(ctx, fmt.Sprintf("failed to shutdown server: %v", err))
		return err
	}

	return nil
}
