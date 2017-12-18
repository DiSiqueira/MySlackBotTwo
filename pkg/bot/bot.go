// Package bot provides a simple to use IRC, Slack and Telegram bot
package bot

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/robfig/cron"
)

const (
	// CmdPrefix is the prefix used to identify a command.
	// !hello would be identified as a command
	CmdPrefix = "!"
)

// Bot handles the bot instance
type Bot struct {
	handlers     *Handlers
	cron         *cron.Cron
	disabledCmds []string
}

// ResponseHandler must be implemented by the protocol to handle the bot responses
type ResponseHandler func(target, message string, sender *User)

// Handlers that must be registered to receive callbacks from the bot
type Handlers struct {
	Response ResponseHandler
}

var (
	// std is the name of the standard bot
	std *Bot
)

// New configures a new bot instance
func New(h *Handlers) *Bot {
	std = &Bot{
		handlers: h,
		cron:     cron.New(),
	}
	std.startPeriodicCommands()
	return std
}

func (b *Bot) startPeriodicCommands() {
	for _, cfg := range periodicCommands {
		StartPeriodicCommand(b, cfg)
	}
	if len(b.cron.Entries()) > 0 {
		b.cron.Start()
	}
}

func StartPeriodicCommand(b *Bot, cfg PeriodicConfig) error {
	err := b.cron.AddFunc(cfg.CronSpec, func() {
		for _, channel := range cfg.Channels {
			message, err := cfg.CmdFunc(channel)
			if err != nil {
				log.Println("Periodic command failed ", err)
			} else if message != "" {
				b.handlers.Response(channel, message, nil)
			}
		}
	})

	if err != nil {
		return err
	}

	b.cron.Run()
	return nil
}

// MessageReceived must be called by the protocol upon receiving a message
func (b *Bot) MessageReceived(channel *ChannelData, message *Message, sender *User) {
	command, err := parse(message.Text, channel, sender)
	if err != nil {
		b.handlers.Response(channel.Channel, err.Error(), sender)
		return
	}

	if command == nil {
		b.executePassiveCommands(&PassiveCmd{
			Raw:         message.Text,
			MessageData: message,
			Channel:     channel.Channel,
			ChannelData: channel,
			User:        sender,
		})
		return
	}

	if b.isDisabled(command.Command) {
		return
	}

	switch command.Command {
	case helpCommand:
		b.help(command)
	default:
		b.handleCmd(command)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
