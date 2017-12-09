package image

import (
	"github.com/disiqueira/MySlackBotTwo/pkg/bot"
	"github.com/disiqueira/MySlackBotTwo/pkg/provider"
	"strings"
	"github.com/disiqueira/MySlackBotTwo/pkg/plugins/web"
	"fmt"
)

const (
	seeUsage = "Invalid args, see usage with: !help puppet."
)

var (
	api provider.ImageRecognition
)

func init() {
	bot.RegisterCommand(
		"image",
		"Interact with images",
		"tags url",
		image)
}

func getProvider() provider.ImageRecognition {
	if api == nil {
		api = provider.NewImageRecognition(bot.Configs().ClarifaiToken())
	}
	return api
}

func image(command *bot.Cmd) (string, error) {
	if !argsValid(command) {
		return seeUsage, nil
	}
	link := web.ExtractURL(command.Raw)

	res, err := getProvider().Analyze([]string{link})
	if err != nil {
		if err == provider.ErrNoConceptsFound {
			return "", nil
		}
		return "",err
	}

	var plainTextList []string
	for key, value := range res[link] {
		plainTextList = append(plainTextList, fmt.Sprintf("%s *%6.4f*", key, value))
	}
	return strings.Join(plainTextList, "\n"), nil
}

func argsValid(command *bot.Cmd) bool {
	return validCommand(command.Args[0]) && web.ExtractURL(command.Raw) != ""
}

func validCommand(cmd string) bool {
	return cmd == "tags"
}