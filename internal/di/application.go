package di

import (
	"github.com/kevinrobayna/goprod/internal/products"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

var ApplicationModule = fx.Module("app",
	fx.Provide(
		provideLogger,
		provideGormDB,
		provideRepository,
		provideService,
	),
	fx.Invoke(products.InvokeMigrate),
)

func provideLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

func FxEvent(logger *zap.Logger) fxevent.Logger {
	return &fxevent.ZapLogger{Logger: logger}
}

func provideGormDB(logger *zap.Logger) (*gorm.DB, error) {
	gormLogger := zapgorm2.New(logger)
	gormLogger.SetAsDefault() // configure gorm to use this zapgorm.Logger for callbacks
	db, err := gorm.Open(sqlite.Open("./models.sqlite"), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func provideRepository(db *gorm.DB) products.IRepository {
	return products.NewRepository(db)
}

func provideService(repository products.IRepository) products.IService {
	return products.NewService(repository)
}
