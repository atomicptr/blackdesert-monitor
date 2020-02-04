package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"

	"github.com/atomicptr/blackdesert-monitor/monitor"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// logger
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// config
	filePath := "./settings.yaml"

	// arguments found, try to use that as a path
	if len(os.Args) > 1 {
		logger.Printf("Arguments found: %v\n", os.Args[1:])
		filePath = os.Args[1]
	}

	config, err := monitor.ConfigFromFile(filePath)
	if err != nil {
		return errors.Wrapf(err, "could not open file \"%s\"", filePath)
	}

	if err := config.Validate(); err != nil {
		return err
	}

	if config.Telegram.UserId == 0 {
		logger.Println("Your Telegram User ID is empty! You can get it by sending /myid to the bot!")
	}

	// channel to listen for errors coming from the app
	appErrors := make(chan error, 1)

	// app starting
	logger.Printf("main: blackdesert-monitor starting...")
	defer logger.Printf("main: Done")

	// channel to listen for interrupt or terminate signal from OS
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	pm, err := monitor.New(
		config,
		logger,
	)
	if err != nil {
		return err
	}

	go func() {
		appErrors <- pm.Start()
	}()

	select {
	case err := <-appErrors:
		return errors.Wrap(err, "app error")
	case sig := <-shutdown:
		logger.Printf("main: %v shutdown...", sig)
		pm.Stop()
	}

	return nil
}
