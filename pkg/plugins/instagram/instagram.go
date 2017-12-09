package image

import (
	"github.com/disiqueira/MySlackBotTwo/pkg/bot"
	"github.com/disiqueira/MySlackBotTwo/pkg/provider"
	"strings"
)

const (
	seeUsage = "Invalid args, see usage with: !help puppet."
	invalidParams         = "Invalid parameters"
	numPhotos = 5
)

var (
	api provider.Instagram
)

func init() {
	bot.RegisterCommand(
		"instagram",
		"Interact with Instagram",
		"(last|follow|stories|best) profile",
		image)
}

func getProvider() (provider.Instagram, error) {
	var err error
	if api == nil {
		api, err = provider.NewInstagram(bot.Configs().InstagramUsername(), bot.Configs().InstagramPassword())
	}
	return api, err
}

func last(command *bot.Cmd) (string, error) {
	ig, err := getProvider()
	if err != nil {
		return "",err
	}

	photos, err := ig.LastPhotos(command.Args[1], numPhotos)
	if err != nil {
		return "",err
	}

	var plainTextList []string
	for _, value := range photos {
		plainTextList = append(plainTextList, value)
	}
	return strings.Join(plainTextList, " \n "), nil
}

func image(command *bot.Cmd) (string, error) {
	if len(command.Args) != 2 {
		return invalidParams, nil
	}

	switch command.Args[0] {
	case "last":
		return last(command)
	default:
		return invalidParams, nil
	}
}
