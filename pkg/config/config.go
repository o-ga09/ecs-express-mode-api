package config

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/caarlos0/env/v6"
	"github.com/o-ga09/ecs-express-mode-api/pkg/errors"
	"google.golang.org/api/option"
)

type Env string

const CtxEnvKey Env = "env"

type Config struct {
	Env                       string `env:"ENV" envDefault:"dev"`
	Port                      string `env:"PORT" envDefault:"8080"`
	Database_url              string `env:"DATABASE_URL" envDefult:""`
	Sentry_DSN                string `env:"SENTRY_DSN" envDefult:""`
	ProjectID                 string `env:"PROJECT_ID" envDefult:""`
	CLOUDFLARE_R2_ACCOUNT_ID  string `env:"CLOUDFLARE_R2_ACCOUNT_ID" envDefult:""`
	CLOUDFLARE_R2_ACCESSKEY   string `env:"CLOUDFLARE_R2_ACCESSKEY" envDefult:""`
	CLOUDFLARE_R2_SECRETKEY   string `env:"CLOUDFLARE_R2_SECRETKEY" envDefult:""`
	CLOUDFLARE_R2_BUCKET_NAME string `env:"CLOUDFLARE_R2_BUCKET_NAME" envDefult:""`
	COOKIE_DOMAIN             string `env:"COOKIE_DOMAIN" envDefault:"localhost"`
}

func New(ctx context.Context) (context.Context, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, errors.Wrap(ctx, err)
	}

	return context.WithValue(ctx, CtxEnvKey, cfg), nil
}

// GetFirebaseApp はFirebaseアプリケーションインスタンスを取得します
func GetFirebaseApp(ctx context.Context) (*firebase.App, error) {
	var app *firebase.App
	var err error
	var credentialsPath string

	env := GetCtxEnv(ctx)

	if env.Env == "local" {
		// 環境変数からサービスアカウントのパスを取得
		credentialsPath = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		opt := option.WithCredentialsFile(credentialsPath)
		app, err = firebase.NewApp(ctx, nil, opt)
	} else {
		// Cloud Run環境: デフォルトの認証情報を使用
		app, err = firebase.NewApp(ctx, nil)
	}

	if err != nil {
		return nil, errors.Wrap(ctx, err)
	}

	return app, nil
}

func GetCtxEnv(ctx context.Context) *Config {
	var cfg *Config
	var ok bool
	if cfg, ok = ctx.Value(CtxEnvKey).(*Config); !ok {
		log.Fatal("config not found")
	}
	return cfg
}
