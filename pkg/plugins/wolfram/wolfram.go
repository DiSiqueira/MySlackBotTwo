package wolfram

import (
	"fmt"
	"github.com/disiqueira/MySlackBotTwo/pkg/bot"
	"github.com/disiqueira/MySlackBotTwo/pkg/provider"
	"strings"
)

const (
	multivac = "INSUFFICIENT DATA FOR A MEANINGFUL ANSWER"
)

var (
	api provider.Wolfram
)

func init() {
	bot.RegisterCommand(
		"wolfram",
		"Answer any question in natural language",
		"question",
		wolfram)
}

func getProvider() provider.Wolfram {
	if api == nil {
		api = provider.NewWolfram(bot.Configs().WolframToken())
	}
	return api
}

func wolfram(command *bot.Cmd) (string, error) {
	wolframAnswer, err := getProvider().Ask(command.RawArgs)
	if err != nil {
		return "",err
	}

	var plainTextList []string
	for _, pod := range wolframAnswer.Queryresult.Pods {
		plainTextList = filterAndAddToList(plainTextList,fmt.Sprintf("*%s:*", pod.Title))
		for _, subPod := range pod.Subpods {
			plainTextList = filterAndAddToList(plainTextList, fmt.Sprintf("*%s:*", subPod.Title))
			plainTextList = filterAndAddToList(plainTextList, subPod.Plaintext)
		}
	}

	if len(plainTextList) == 0 {
		return multivac,nil
	}

	return strings.Join(plainTextList, "\n"), nil
}

func filterAndAddToList(list []string, item string) []string {
	item = strings.Replace(item, "|", ":", -1)
	item = strings.TrimSpace(item)
	if item != "" && len(item) > 3 {
		return append(list, item)
	}
	return list
}
