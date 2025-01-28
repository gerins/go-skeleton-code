package app

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"

	"go-skeleton-code/pkg/log"
	"go-skeleton-code/config"
	"go-skeleton-code/internal/app/domains/order"
	"go-skeleton-code/internal/app/domains/user"
	"go-skeleton-code/pkg/gorm"
	"go-skeleton-code/pkg/kafka"
	"go-skeleton-code/pkg/redis"
)

func Init(gin *gin.Engine, g *grpc.Server, cfg *config.Config) chan bool {
	var (
		exitSignal         = make(chan bool)
		validator          = validator.New()
		apiTimeout         = cfg.App.HTTP.CtxTimeout
		redis              = redis.Init(cfg.Dependencies.Cache)
		readDatabase       = gorm.InitPostgres(cfg.Dependencies.Database.Read)
		writeDatabase      = gorm.InitPostgres(cfg.Dependencies.Database.Write)
		matchOrderConsumer = kafka.NewConsumer(cfg.Dependencies.MessageBroker, cfg.Dependencies.MessageBroker.Consumer.Topic.MatchOrder)
		producer, writer   = kafka.NewProducer(cfg.Dependencies.MessageBroker.Brokers)
	)

	// Init http router
	{
		// Repository
		userRepository := user.NewRepository(readDatabase, writeDatabase)
		orderRepository := order.NewRepository(readDatabase, writeDatabase)

		// Usecase
		userUsecase := user.NewUsecase(cfg.Security, validator, userRepository)
		orderUsecase := order.NewUsecase(writeDatabase, producer, validator, orderRepository, userRepository)

		// Handler
		api := gin.Group("/api")
		user.NewHTTPHandler(userUsecase, apiTimeout).InitRoutes(api)
		order.NewHTTPHandler(orderUsecase, apiTimeout, cfg.Security).InitRoutes(api)

		// Queue
		order.NewQueueHandler(matchOrderConsumer, orderUsecase, apiTimeout).StartConsumer()
	}

	// Graceful shutdown
	go func() {
		<-exitSignal // Receive exit signal
		log.Info("disconnecting service dependencies")

		if err := matchOrderConsumer.Close(); err != nil {
			log.Error(err)
		}

		if err := writer.Close(); err != nil {
			log.Error(err)
		}

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
