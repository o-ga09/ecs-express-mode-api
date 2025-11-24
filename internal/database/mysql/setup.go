package mysql

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/o-ga09/ecs-express-mode-api/pkg/config"
	CtxLogger "github.com/o-ga09/ecs-express-mode-api/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/opentelemetry/tracing"
)

type ctxKey string

const CtxKey ctxKey = "db"
const MAX_RETRY = 10

var (
	db   *gorm.DB
	once sync.Once
)

func New(ctx context.Context) context.Context {
	once.Do(func() {
		ctx, err := config.New(ctx)
		if err != nil {
			log.Fatal(err)
		}
		env := config.GetCtxEnv(ctx)
		logger := NewSentryLogger()

		dialector := mysql.Open(env.Database_url)
		db, err = gorm.Open(dialector, &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: false,
			},
			Logger: logger,
		})
		if err != nil {
			connect(ctx, dialector)
		}

		err = db.Use(tracing.NewPlugin(tracing.WithoutMetrics()))
		if err != nil {
			log.Fatal(err)
		}
		// SQLDBインスタンスを取得
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal(err)
		}

		// コネクションプールの設定
		sqlDB.SetMaxIdleConns(10)           // アイドル状態の最大接続数
		sqlDB.SetMaxOpenConns(100)          // 最大接続数
		sqlDB.SetConnMaxLifetime(time.Hour) // 接続の最大生存期間
	})

	return context.WithValue(ctx, CtxKey, db)
}

func connect(ctx context.Context, dialector gorm.Dialector) context.Context {
	var err error
	for i := 0; i < MAX_RETRY; i++ {
		if db, err = gorm.Open(dialector, &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: false,
			},
			Logger: logger.Default.LogMode(logger.Info),
		}); err == nil {
			return context.WithValue(ctx, CtxKey, db)
		}
		time.Sleep(5 * time.Second)
		CtxLogger.Error(ctx, "Failed to connect to database, retrying...", "attempt", i+1, "error", err)
	}
	return ctx
}

func CtxFromDB(ctx context.Context) *gorm.DB {
	return ctx.Value(CtxKey).(*gorm.DB).WithContext(ctx).Debug()
}
