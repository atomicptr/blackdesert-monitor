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

	filePath := "./settings.yaml"
	config, err := monitor.ConfigFromFile(filePath)
	if err != nil {
		return errors.Wrapf(err, "could not open file \"%s\"", filePath)
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
