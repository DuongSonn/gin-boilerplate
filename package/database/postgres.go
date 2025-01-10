package database

import (
	"fmt"
	"log"
	"oauth-server/config"
	"oauth-server/package/tracing"
	"oauth-server/utils"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var postgresDB *gorm.DB

func InitPostgres() {
	conf := config.GetConfiguration().PostgresDB
	serverConf := config.GetConfiguration().Server

	// DB logging config
	logLevel := logger.Info
	if serverConf.Mode == utils.RELEASE_MODE {
		logLevel = logger.Silent
	}
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logLevel,               // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,                  // Enable color
		},
	)

	// DB connection
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		conf.Host,
		conf.User,
		conf.Password,
		conf.Database,
		conf.Port,
	)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		log.Fatalf("error_connecting_to_database: %v", err)
	}

	if err := conn.Use(tracing.NewSentryPlugin()); err != nil {
		log.Fatalf("error_connecting_to_database: %v", err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		log.Fatalf("error_get_database: %v", err)
	}
	sqlDB.SetConnMaxIdleTime(time.Hour)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	postgresDB = conn
}

func GetPostgres() *gorm.DB {
	return postgresDB
}

func BeginPostgresTransaction() *gorm.DB {
	return GetPostgres().Begin()
}
