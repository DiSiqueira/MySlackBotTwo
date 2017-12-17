package echo

import (
	"github.com/disiqueira/MySlackBotTwo/pkg/bot"
	"github.com/CrowdSurge/banner"
	"strings"
	"errors"
)

func init() {
	bot.RegisterCommand(
		"echo",
		"Write a banner in ASC ART.",
		"text",
		echo)
}

func echo(command *bot.Cmd) (string, error) {
	if len(command.Args) == 0 {
		return "", errors.New("missed argument $text")
	}
	return transform(strings.Join(command.Args, " ")), nil
}

func transform(text string) string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	return "```"+banner.PrintS(text)+"```"
}
