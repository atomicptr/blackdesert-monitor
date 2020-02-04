package monitor

import (
	"github.com/pkg/errors"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/atomicptr/blackdesert-monitor/system"
)

type ProcessMonitor struct {
	config             *Config
	bot                *tb.Bot
	logger             *log.Logger
	ticker             *time.Ticker
	tickerChan         chan bool
	appErrors          chan error
	lastKnownProcessId int
}

func New(config *Config, logger *log.Logger) (*ProcessMonitor, error) {
	bot, err := tb.NewBot(tb.Settings{
		Token:  config.Telegram.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not init telegram bot")
	}

	pm := &ProcessMonitor{
		config:    config,
		bot:       bot,
		logger:    logger,
		ticker:    time.NewTicker(config.PollInterval),
		appErrors: make(chan error, 1),
	}

	pm.registerBotHandlers()

	return pm, nil
}

func (pm *ProcessMonitor) Start() error {
	go pm.bot.Start()
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
	err := pm.Status()
	if err != nil {
		pm.logger.Println(err)
		return
	}

	pm.logger.Printf("Black Desert (pid: %d) is running fine...\n", pm.lastKnownProcessId)
}

func (pm *ProcessMonitor) Status() error {
	processId, err := system.FindProcessByName(pm.config.ProcessName)
	if err != nil {
		return err
	}
	pm.lastKnownProcessId = processId

	connections, err := system.GetConnectionsWithPid(processId)
	if err != nil {
		return err
	}

	numConnections := len(connections)

	if numConnections == 0 {
		// no connections -> dc!?
		return errors.New("disconnected")
	}

	return nil
}

func (pm *ProcessMonitor) Stop() {
	pm.bot.Stop()
	pm.ticker.Stop()
	close(pm.tickerChan)
}
