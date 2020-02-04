package monitor

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/atomicptr/blackdesert-monitor/system"
)

type ProcessMonitor struct {
	config                    *Config
	bot                       *tb.Bot
	logger                    *log.Logger
	ticker                    *time.Ticker
	tickerChan                chan bool
	appErrors                 chan error
	lastKnownProcessId        int
	unavailabilityCounter     int
	unavailabilityMessageSent bool
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
		config:                    config,
		bot:                       bot,
		logger:                    logger,
		ticker:                    time.NewTicker(config.PollInterval),
		appErrors:                 make(chan error, 1),
		unavailabilityCounter:     0,
		unavailabilityMessageSent: false,
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

		if !pm.unavailabilityMessageSent {
			pm.unavailabilityCounter++
			pm.logger.Printf("unavailability counter raised to %d/%d\n", pm.unavailabilityCounter, pm.config.UnavailabilityThreshold)
		}

		if !pm.unavailabilityMessageSent && pm.unavailabilityCounter >= pm.config.UnavailabilityThreshold {
			// inform user that the Black Desert connection has been cut
			pm.logger.Println("Unavailability threshold reached, will inform user...")
			err := pm.sendMessageToUser(
				fmt.Sprintf("Connection to Black Desert Online has been lost, Reason: %s", err),
			)
			if err != nil {
				pm.logger.Println(err)
			}

			pm.unavailabilityMessageSent = true

			if pm.config.CloseBlackDesertWhenUnavailable {
				pm.logger.Println("...and close Black Desert")
				err := system.Kill(pm.lastKnownProcessId)
				if err != nil {
					pm.logger.Println(err)
				}
			}
		}

		return
	}

	if pm.unavailabilityMessageSent {
		pm.logger.Println("connection has been re-established")
	}

	// reset unavailability stuff
	pm.unavailabilityCounter = 0
	pm.unavailabilityMessageSent = false

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
		return errors.New("connection has been lost")
	}

	return nil
}

func (pm *ProcessMonitor) Stop() {
	pm.bot.Stop()
	pm.ticker.Stop()
	close(pm.tickerChan)
}
