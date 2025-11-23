package main

import (
	"context"
	"log"

	"github.com/o-ga09/ecs-express-mode-api/internal/server"
	"github.com/o-ga09/ecs-express-mode-api/pkg/config"
)

func main() {
	ctx := context.Background()

	ctx, err := config.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// ロガーの設定
	server.Logger(ctx)

	srv := server.NewServer(ctx)
	if err := srv.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
