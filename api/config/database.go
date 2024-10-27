package config

import (
	"fmt"
	"time"

	"github.com/candrairwn/go-pure/api/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type ConfigDefaultDatabase struct {
	Username              string
	Password              string
	Host                  string
	Port                  string
	Database              string
	IdleConnection        int
	MaxConnection         int
	MaxIdleTimeConnection int
	maxLifeTimeConnection int
}

func NewDatabasePostgres(viper *viper.Viper, log *zap.SugaredLogger) (*gorm.DB, error) {
	password, err := utils.ReadFile(viper.GetString("DB_PASSWORD"), log)
	if err != nil {
		log.Fatalf("failed to read db password: %v", err)
	}

	config := ConfigDefaultDatabase{
		Username:              viper.GetString("DB_USERNAME"),
		Password:              password,
		Host:                  viper.GetString("DB_HOST"),
		Port:                  viper.GetString("DB_PORT"),
		Database:              viper.GetString("DB_DATABASE"),
		IdleConnection:        viper.GetInt("DB_IDLE_CONNECTION"),
		MaxConnection:         viper.GetInt("DB_MAX_CONNECTION"),
		MaxIdleTimeConnection: viper.GetInt("DB_MAX_IDLE_TIME_CONNECTION"),
		maxLifeTimeConnection: viper.GetInt("DB_MAX_LIFE_TIME_CONNECTION"),
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.Username, config.Password, config.Database)

	dsnMigrate := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.Username, config.Password, config.Host, config.Port, config.Database)

	if err := MigrateRun(dsnMigrate, log); err != nil {
		log.Fatalf("failed to migrate: %v", err)
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "lsp_cbt.",
		},
		Logger: logger.New(&zapWriter{Logger: log}, logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
		return nil, err
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get db connection: %v", err)
		return nil, err
	}

	connection.SetMaxIdleConns(config.IdleConnection)
	connection.SetMaxOpenConns(config.MaxConnection)
	connection.SetConnMaxIdleTime(time.Duration(config.MaxIdleTimeConnection) * time.Minute)
	connection.SetConnMaxLifetime(time.Duration(config.maxLifeTimeConnection) * time.Minute)

	log.Info("database connected")
	return db, nil
}

func MigrateRun(dsnMigrate string, log *zap.SugaredLogger) error {
	m, err := migrate.New("file://migrations", dsnMigrate)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		log.Infof("migrate: %v", err)
		if err.Error() == "no change" {
			return nil
		} else {
			return err
		}
	}

	log.Info("migrate success")

	return nil
}

type zapWriter struct {
	Logger *zap.SugaredLogger
}

func (zw *zapWriter) Printf(message string, v ...interface{}) {
	zw.Logger.Infof(message, v...)
}
