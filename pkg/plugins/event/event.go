package event

import (
	"github.com/disiqueira/MySlackBotTwo/pkg/bot"
	"fmt"
	"time"
)

func init() {
	event := bot.RegisterCommandV2(
		"event",
		"Create events for your channels",
		"channel date hour message",
		registerEvent)

	event.MinArgs = 4
}

func registerEvent(command *bot.Cmd) (bot.CmdResult, error) {
	result := bot.CmdResult{}

	date, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", command.Args[1], command.Args[2]))
	if err != nil {
		result.Channel = command.Channel
		result.Message = err.Error()
		return result, err
	}

	config := bot.PeriodicConfig{
		CronSpec: toCron(date),
		Channels: []string{command.Args[0]},
		CmdFunc:  func (_ string) (string, error) {
			return command.Args[3], nil
		},
	}

	bot.EnablePeriodicCommand(config)

	result.Channel = command.Args[0]
	result.Message = "Event registered!"
	return result, nil
}

func toCron(date time.Time) string {
	const format = "0 %d %d %d %d *"
	return fmt.Sprintf(format,
		date.Format("04"),
		date.Format("15"),
		date.Format("02"),
		date.Format("01"))
}