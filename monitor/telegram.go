package monitor

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (pm *ProcessMonitor) registerBotHandlers() {
	pm.bot.Handle("/myid", func(message *tb.Message) {
		_, err := pm.bot.Send(message.Sender, fmt.Sprintf("Your Telegram ID is: %d", message.Sender.ID))
		if err != nil {
			pm.logger.Println(err)
		}
	})

	pm.bot.Handle("/status", pm.wrapAuthorizedHandler(pm.botHandlerStatus))
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
		_, err := pm.bot.Send(message.Sender, err)
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
