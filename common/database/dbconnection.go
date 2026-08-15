package database

import (
	"log"
	"os"
	"time"

	"github.com/mannanpatel/Restro/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func DBConnection(cfg *config.Config) {
	var err error
	DB, err = gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: newGormLogger(cfg.Server.Env),
	})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("MySQL connected successfully!")
}

func GetDB() *gorm.DB {
	return DB
}

func newGormLogger(env string) gormlogger.Interface {
	level := gormlogger.Warn
	if env == "development" {
		level = gormlogger.Info
	}

	return gormlogger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  level,
			IgnoreRecordNotFoundError: true,
			Colorful:                  env == "development",
		},
	)
}
