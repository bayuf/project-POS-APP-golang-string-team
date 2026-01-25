package database

import (
	"fmt"
	"log"
	"os"

	"time"

	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(config utils.DatabaseCofig, zapLog *zap.Logger) (*gorm.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s",
		config.Username, config.Password, config.Name, config.SSL, config.Host)

	// Setup logger for GORM
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // enable color
		},
	)

	conn, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger:                 newLogger,
		PrepareStmt:            true,
		SkipDefaultTransaction: true, //disabled default transaction
	})

	if err != nil {
		zapLog.Error("failed to connect to database", zap.Error(err))
		return nil, fmt.Errorf("parse config: %w", err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		zapLog.Error("failed to get sql db", zap.Error(err))
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// ---- Pool settings ----
	sqlDB.SetConnMaxIdleTime(time.Duration(config.MaxConn) * time.Minute)
	sqlDB.SetConnMaxLifetime(time.Duration(config.MaxConn) * time.Minute)
	sqlDB.SetMaxIdleConns(config.MaxIdle)
	sqlDB.SetMaxOpenConns(config.MaxOpen)

	return conn, nil
}
