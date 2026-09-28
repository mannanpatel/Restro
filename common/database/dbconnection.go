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

// var DB *gorm.DB

func DBConnection(cfg *config.Config) *gorm.DB {
	var err error
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: newGormLogger(cfg.Server.Env),
	})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	// if err := DB.AutoMigrate(
	// 	&models.User{},
	// 	&models.Product{},
	// 	&models.Cart{},
	// 	&models.CartItem{},
	// 	&models.Order{},
	// 	&models.OrderItem{},
	// 	&models.Payment{},
	// ); err != nil {
	// 	log.Fatalf("Database migration failed: %v", err)
	// }

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("MySQL connected successfully!")

	return db
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
