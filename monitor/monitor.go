package monitor

import (
	"log"
	"time"

	"github.com/atomicptr/blackdesert-monitor/system"
)

type ProcessMonitor struct {
	config Config

	logger     *log.Logger
	ticker     *time.Ticker
	tickerChan chan bool
	appErrors  chan error
}

func New(config Config, logger *log.Logger) *ProcessMonitor {
	return &ProcessMonitor{
		config:    config,
		logger:    logger,
		ticker:    time.NewTicker(config.PollInterval),
		appErrors: make(chan error, 1),
	}
}

func (pm *ProcessMonitor) Start() error {
	pm.tickerChan = make(chan bool)

	for {
		select {
		case err := <-pm.appErrors:
			return err
		case <-pm.tickerChan:
			return nil
		case <-pm.ticker.C:
			pm.tick()
		}
	}
}

func (pm *ProcessMonitor) tick() {
	processId, err := system.FindProcessByName(pm.config.ProcessName)
	if err != nil {
		pm.logger.Println(err)
		return
	}

	conns, err := system.GetConnectionsWithPid(processId)
	if err != nil {
		pm.logger.Println(err)
		return
	}

	log.Printf("BDO PID: %d, Connections: %d, RUNNING? %v\n", processId, len(conns), len(conns) > 0)
}

func (pm *ProcessMonitor) Stop() {
	pm.ticker.Stop()
	close(pm.tickerChan)
}
