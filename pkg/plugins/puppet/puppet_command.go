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

	if !argsValid(command.Args) {
		result.Message = seeUsage
		return
	}

	result.Channel = command.Args[1]
	result.Message = strings.Join(command.Args[2:], " ")
	return
}

func argsValid(args []string) bool {
	return len(args) >= 3 && validCommand(args[0])
}

func validCommand(cmd string) bool {
	return cmd == "say" || cmd == "act"
}

func init() {
	bot.RegisterCommandV2(
		"puppet",
		"Allows you to send messages through the bot",
		"say|act #channel your message",
		sendMessage)
}
