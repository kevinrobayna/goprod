package di

import (
	"fmt"
	"github.com/kevinrobayna/goprod/internal/models"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
	"os"
)

var ApplicationModule = fx.Module("app",
	fx.Provide(provideConfig, provideDbConfig),
	fx.Provide(provideLogger),
	fx.Provide(provideGormDB),
	fx.Invoke(invokeMigrate),
)

func provideLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}

type AppConfig struct {
	DBConfig DBConfig `yaml:"db_config"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"database"`
	SslMode  string `yaml:"sslmode"`
}

func (config DBConfig) GetDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.DbName,
		config.Port,
		config.SslMode,
	)
}

func provideConfig() (AppConfig, error) {
	f, err := os.Open("config.yml")
	if err != nil {
		return AppConfig{}, err
	}
	var cfg AppConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	return cfg, nil
}

func provideDbConfig(config AppConfig) DBConfig {
	return config.DBConfig
}

func FxEvent(logger *zap.Logger) fxevent.Logger {
	return &fxevent.ZapLogger{Logger: logger}
}

func provideGormDB(config DBConfig, logger *zap.Logger) (*gorm.DB, error) {
	gormLogger := zapgorm2.New(logger)
	gormLogger.SetAsDefault() // configure gorm to use this zapgorm.Logger for callbacks
	db, err := gorm.Open(postgres.Open(config.GetDsn()), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func invokeMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.Product{})
	if err != nil {
		return
	}

	db.Create(&models.Product{Code: "D42", Price: 100})
}
