package puppet

import (
	"strings"

	"github.com/disiqueira/MySlackBotTwo/pkg/bot"
)

const (
	seeUsage = "Invalid args, see usage with: !help puppet."
)

func sendMessage(command *bot.Cmd) (result bot.CmdResult, err error) {
	result = bot.CmdResult{}

	if !validCommand(command.Args[0]) {
		result.Message = seeUsage
		return
	}

	result.Channel = command.Args[1]
	result.Message = strings.Join(command.Args[2:], " ")
	return
}

func validCommand(cmd string) bool {
	return cmd == "say" || cmd == "act"
}

func init() {
	puppet := bot.RegisterCommandV2(
		"puppet",
		"Allows you to send messages through the bot",
		"say|act #channel message",
		sendMessage)

	puppet.SetMinArgs(3)
}
