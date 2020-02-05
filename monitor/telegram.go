package monitor

import (
	"fmt"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/atomicptr/blackdesert-monitor/system"
)

func (pm *ProcessMonitor) registerBotHandlers() {
	pm.bot.Handle("/myid", func(message *tb.Message) {
		_, err := pm.bot.Send(message.Sender, fmt.Sprintf("Your Telegram ID is: %d", message.Sender.ID))
		if err != nil {
			pm.logger.Println(err)
		}
	})

	pm.bot.Handle("/status", pm.wrapAuthorizedHandler(pm.botHandlerStatus))
	pm.bot.Handle("/quit", pm.wrapAuthorizedHandler(pm.botHandlerQuit))
}

func (pm *ProcessMonitor) sendMessageToUser(message string) error {
	user, err := pm.bot.ChatByID(strconv.Itoa(pm.config.Telegram.UserId))
	if err != nil {
		return err
	}

	_, err = pm.bot.Send(user, message)
	if err != nil {
		return err
	}

	return nil
}

func (pm *ProcessMonitor) wrapAuthorizedHandler(handlerFunc func(*tb.Message)) func(*tb.Message) {
	return func(message *tb.Message) {
		pm.logger.Printf("authorized method call from: %s (%d)\n", message.Sender.Username, message.Sender.ID)
		if message.Sender.ID != pm.config.Telegram.UserId {
			_, err := pm.bot.Send(message.Sender, "You're not authorized to do this, add your Telegram ID (/myid) to my settings.")
			if err != nil {
				pm.logger.Println(err)
			}
			return
		}

		handlerFunc(message)
	}
}

func (pm *ProcessMonitor) botHandlerStatus(message *tb.Message) {
	err := pm.Status()
	if err != nil {
		_, err := pm.bot.Send(message.Sender, fmt.Sprintf("Black Desert is not running properly.\n\nReason:\n%s", err))
		if err != nil {
			pm.logger.Println(err)
		}
		return
	}

	_, err = pm.bot.Send(message.Sender, "Everything is fine!")
	if err != nil {
		pm.logger.Println(err)
	}
}

func (pm *ProcessMonitor) botHandlerQuit(message *tb.Message) {
	pm.logger.Printf("/quit received, will kill process %d\n", pm.lastKnownProcessId)
	err := system.Kill(pm.lastKnownProcessId)
	if err != nil {
		pm.logger.Println(err)
		_, err := pm.bot.Send(message.Sender, fmt.Sprintf("Could not kill process: %s", err))
		if err != nil {
			pm.logger.Println(err)
		}
		return
	}

	_, err = pm.bot.Send(message.Sender, "Successfully quit Black Desert")
	if err != nil {
		pm.logger.Println(err)
	}
}
