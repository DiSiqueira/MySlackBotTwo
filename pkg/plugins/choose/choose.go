package choose

import (
	"github.com/disiqueira/MySlackBotTwo/pkg/bot"
	"math/rand"
)

func init() {
	bot.RegisterCommand(
		"choose",
		"Randomly picks one of the options.",
		"option1 option2 {option3 ...}",
		choose)
}

func choose(command *bot.Cmd) (string, error) {
	return command.Args[rand.Intn(len(command.Args))], nil
}
