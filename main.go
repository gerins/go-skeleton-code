package main

import (
	"os"
	"os/signal"
	"runtime"

	gin "go-skeleton-code/cmd/gin"
	"go-skeleton-code/config"
	"go-skeleton-code/internal/app"
	"go-skeleton-code/pkg/log"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		cfg  = config.ParseConfigFile("config.yaml")
		http = gin.NewHTTPServer(cfg)
	)

	// Init log
	log.InitWithConfig(log.Config{
		LogToTerminal:     cfg.App.Logging.LogToTerminal,
		LogToFile:         cfg.App.Logging.LogToFile,
		Location:          cfg.App.Logging.Location,
		FileLogName:       cfg.App.Logging.FileLogName,
		MaxAge:            cfg.App.Logging.MaxAge,
		RotationFile:      cfg.App.Logging.RotationFile,
		HideSensitiveData: cfg.App.Logging.HideSensitiveData,
	})

	// Init app
	appExitSignal := app.Init(http.Server, cfg)

	// Run server
	httpExitSignal := http.Run()

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt)
	for range interruptSignal {
		appExitSignal <- true
		httpExitSignal <- true
		<-appExitSignal  // Finish
		<-httpExitSignal // Finish
		return           // Now we can safely exit the app
	}
}
