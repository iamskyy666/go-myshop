package database

import (
	"fmt"

	"github.com/iamskyy666/go-myshop/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(cfg config.DatabaseConfig)(*gorm.DB, error){
	dsn:=fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	// connection
	db,err:=gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err!=nil{
		return nil, fmt.Errorf("database connection failed: %w",err)
	}

	return db,nil
}