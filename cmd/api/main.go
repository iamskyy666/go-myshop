package main

import (
	"github.com/gin-gonic/gin"
	"github.com/iamskyy666/go-myshop/internal/config"
	"github.com/iamskyy666/go-myshop/internal/database"
	"github.com/iamskyy666/go-myshop/internal/logger"
)

func main() {
	log := logger.New()
	cfg,err:=config.Load()

	if err!=nil{
		log.Fatal().Err(err).Msg("Failed to load config!")
	}

	db,err:=database.New(cfg.Database)

	if err!=nil{
		log.Fatal().Err(err).Msg("Failed to connect to DATABASE!")
	}

	mainDB,err:= db.DB()

	if err!=nil{
		log.Fatal().Err(err).Msg("Failed to get database-connection!")
	}

	defer mainDB.Close()

	gin.SetMode(cfg.Server.GinMode)

	log.Info().Msg("🟢 Starting server...")

}