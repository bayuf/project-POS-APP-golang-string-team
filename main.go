package main

import (
	"fmt"
	"log"

	"github.com/bayuf/project-POS-APP-golang-string-team/cmd"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/wire"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/database"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// init config
	config, err := utils.ReadConfiguration()
	if err != nil {
		log.Fatal("error :", err)
	}

	//init logger
	logger, err := utils.InitLogger(config.PathLogging, config.Debug)
	if err != nil {
		log.Fatal("error :", err)
	}

	// set gin mode
	gin.SetMode(config.GinMode)

	// init database
	dbPool, err := database.InitDB(config.DB, logger)
	if err != nil {
		logger.Error("cant init database :", zap.Error(err))
		log.Fatal("cant init database :", err)
	}

	// migrate database
	if config.DB.DBMigrate {
		if err := data.Migrate(dbPool); err != nil {
			logger.Error("cant migrate database :", zap.Error(err))
			log.Fatal("cant migrate database :", err)
		}
	}

	// Seeder
	if config.DB.DBSeeder {
		fmt.Println("seeder runnn...")
		if err := data.SeedAll(dbPool, logger); err != nil {
			logger.Error("cant seed database :", zap.Error(err))
			log.Fatal("cant seed database :", err)
		} else {
			fmt.Println("seed stop")
		}
	}

	// init layer
	repo := repository.NewRepository(dbPool, logger)
	app := wire.Wiring(dbPool, repo, logger, config)

	// start app
	fmt.Println(config.AppName, "is starting...")
	cmd.APIserver(app, config, logger)
}
