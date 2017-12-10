package image

import (
	"github.com/disiqueira/MySlackBotTwo/pkg/bot"
	"github.com/disiqueira/MySlackBotTwo/pkg/provider"
	"strings"
	"fmt"
)

const (
	invalidParams         = "Invalid parameters"
	numLastPhotos = 5
	numBestPhotos       = 100
	minConceptValue = 0.9
)

var (
	apiInstagram provider.Instagram
	apiImageRecognition provider.ImageRecognition
	conceptsDefault  = []string{"bikini", "lingerie", "sexy", "pretty", "glamour", "seduction"}
)

func init() {
	bot.RegisterCommand(
		"instagram",
		"Interact with Instagram",
		"(last|follow|stories|best|bikini) profile",
		image)
}

func image(command *bot.Cmd) (string, error) {
	if len(command.Args) != 2 {
		return invalidParams, nil
	}

	switch command.Args[0] {
	case "last":
		return lastCmd(command)
	case "best":
		return bestCmd(command)
	case "bikini":
		return bikiniCmd(command)
	case "sexy":
		return sexyCmd(command)
	case "follow":
		return followCmd(command)
	default:
		return invalidParams, nil
	}
}

func lastCmd(command *bot.Cmd) (string, error) {
	photos, err := lastPhotos(command, numLastPhotos)
	if err != nil {
		return "",err
	}

	var plainTextList []string
	for _, value := range photos {
		plainTextList = append(plainTextList, value)
	}
	return strings.Join(plainTextList, " \n "), nil
}

func bestCmd(command *bot.Cmd) (string, error) {
	return conceptImages(command, conceptsDefault)
}

func bikiniCmd(command *bot.Cmd) (string, error) {
	return conceptImages(command, []string{"bikini"})
}

func sexyCmd(command *bot.Cmd) (string, error) {
	return conceptImages(command, []string{"sexy"})
}

func followCmd(command *bot.Cmd) (string, error) {
	ig, err := instagram()
	if err != nil {
		return "", err
	}
	if err = ig.Follow(command.Args[1]); err != nil {
		return fmt.Sprintf("Follow: Err: %s", err.Error()), err
	}

	return fmt.Sprintf("Following: %s", command.Args[1]), nil
}

func lastPhotos(command *bot.Cmd, num int) ([]string, error) {
	ig, err := instagram()
	if err != nil {
		return nil,err
	}

	return ig.LastPhotos(command.Args[1], num)
}

func instagram() (provider.Instagram, error) {
	var err error
	if apiInstagram == nil {
		apiInstagram, err = provider.NewInstagram(bot.Configs().InstagramUsername(), bot.Configs().InstagramPassword())
	}
	return apiInstagram, err
}

func conceptImages(command *bot.Cmd, listConcepts []string) (string, error) {
	photos, err := lastPhotos(command, numBestPhotos)
	if err != nil {
		return "",err
	}

	concepts, err := imageRecognition().Analyze(photos)
	if err != nil {
		return "", err
	}

	var final []string
	for url, concept := range concepts {
		if shouldInsertImage(concept, listConcepts, minConceptValue) {
			final = append(final, url)
		}
	}
	return strings.Join(final, "\n"), nil
}

func imageRecognition() (provider.ImageRecognition) {
	if apiImageRecognition == nil {
		apiImageRecognition = provider.NewImageRecognition(bot.Configs().ClarifaiToken())
	}
	return apiImageRecognition
}

func shouldInsertImage(concept provider.Concepts, validateConcepts []string, minValue float64) bool {
	for _, conceptFilter := range validateConcepts {
		val, ok := concept[conceptFilter]
		if ok && val > minValue {
			return true
		}
	}
	return false
}
