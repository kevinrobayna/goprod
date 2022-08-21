package internal

import (
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
	"os"
)

var Module = fx.Module("app",
	fx.Provide(provideConfig, provideDbConfig),
	fx.Provide(provideLogger),
	fx.Provide(provideGormDB),
	fx.Provide(provideRepository),
	fx.Provide(provideService),
	fx.Invoke(invokeMigrate),
)

func provideLogger() *zap.Logger {
	logger, _ := zap.NewProduction(
		zap.AddCaller(),
	)
	return logger
}

type AppConfig struct {
	Build    BuildConfig
	DBConfig DBConfig `yaml:"db_config"`
}

type BuildConfig struct {
	ConfigFile string
	Sha        string
	Date       string
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
		config.SslMode)
}

func provideConfig(logger *zap.Logger, build BuildConfig) (AppConfig, error) {
	f, err := os.Open(build.ConfigFile)
	if err != nil {
		logger.Fatal("Unable to open config file", zap.Error(err))
		return AppConfig{}, err
	}
	var cfg AppConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		logger.Fatal("Unable to decode config file", zap.Error(err))
		return AppConfig{}, err
	}
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
		logger.Fatal("Unable to establish connection with database", zap.Error(err))
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
