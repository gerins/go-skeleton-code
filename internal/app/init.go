package app

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"go-skeleton-code/config"
	"go-skeleton-code/internal/app/handler"
	"go-skeleton-code/internal/app/middleware"
	"go-skeleton-code/internal/app/repository"
	"go-skeleton-code/internal/app/usecase"
	"go-skeleton-code/pkg/gorm"
	"go-skeleton-code/pkg/log"
	"go-skeleton-code/pkg/redis"
)

func Init(gin *gin.Engine, cfg *config.Config) chan bool {
	var (
		exitSignal     = make(chan bool)
		validator      = validator.New()
		defaultTimeout = cfg.App.HTTP.CtxTimeout
		redis          = redis.Init(cfg.Dependencies.Cache)
		readDatabase   = gorm.InitPostgres(cfg.Dependencies.Database.Read)
		writeDatabase  = gorm.InitPostgres(cfg.Dependencies.Database.Write)
	)

	/**********************************************
	 *                Repository
	 *********************************************/
	fuelRepository := repository.NewFuelRepository(readDatabase, writeDatabase)

	/**********************************************
	 *                 Usecase
	 *********************************************/
	fuelUsecase := usecase.NewFuelUsecase(cfg.Security, fuelRepository)

	/**********************************************
	 *                 Handler
	 *********************************************/
	master := gin.Group("/v3/master", middleware.SetAPITimeout(defaultTimeout))
	{
		handler.NewFuelHandler(defaultTimeout, validator, fuelUsecase).InitRoutes(master)
	}

	transaction := gin.Group("/v3/transaction")
	{
		handler.NewFuelHandler(defaultTimeout, validator, fuelUsecase).InitRoutes(transaction)
	}

	report := gin.Group("/v3/report")
	{
		handler.NewFuelHandler(defaultTimeout, validator, fuelUsecase).InitRoutes(report)
	}

	// Graceful shutdown
	go func() {
		<-exitSignal // Receive exit signal
		log.Info("disconnecting service dependencies")

		if err := redis.Close(); err != nil {
			log.Error(err)
		}

		if readDatabase, err := readDatabase.DB(); err == nil {
			if err = readDatabase.Close(); err != nil {
				log.Error(err)
			}
		}

		if writeDatabase, err := writeDatabase.DB(); err == nil {
			if err = writeDatabase.Close(); err != nil {
				log.Error(err)
			}
		}

		log.Info("finished disconnecting service dependencies")
		exitSignal <- true // Send signal already finish the job
	}()

	return exitSignal
}
