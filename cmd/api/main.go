package main

import (
	"context"
	"log"

	"github.com/o-ga09/ecs-express-mode-api/internal/server"
	"github.com/o-ga09/ecs-express-mode-api/pkg/config"
	"github.com/o-ga09/ecs-express-mode-api/pkg/logger"
)

func main() {
	ctx := context.Background()

	ctx, err := config.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// ロガーの設定
	logger.Logger(ctx)

	srv := server.NewServer(ctx)
	if err := srv.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
