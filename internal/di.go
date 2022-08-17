package internal

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

var Module = fx.Module("app",
	fx.Provide(provideLogger),
	fx.Provide(provideGormDB),
	fx.Provide(provideRepository),
	fx.Provide(provideService),
	fx.Invoke(invokeMigrate),
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

func invokeMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&product{})
	if err != nil {
		return
	}
}

func provideRepository(db *gorm.DB) iRepository {
	return &repository{db: db}
}

func provideService(repository iRepository) IService {
	return &Service{repository: repository}
}
