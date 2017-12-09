package weather

import (
	"fmt"
	"github.com/disiqueira/MySlackBotTwo/pkg/bot"
	"github.com/disiqueira/MySlackBotTwo/pkg/provider"
)

const (
	weatherAnswerFormat = "%s, %s - Current: %s %-2.0fC, Humidity: %d%% High: %-2.0fC, Low: %-2.0fC"
)

var (
	api provider.Weather
)

func init() {
	bot.RegisterCommand(
		"weather",
		"Returns the actual weather of a region.",
		"city",
		weather)
}

func getProvider() provider.Weather {
	if api == nil {
		api = provider.NewWeather(bot.Configs().OpenWeatherToken())
	}
	return api
}

func weather(command *bot.Cmd) (string, error) {
	resp, err := getProvider().ByName(command.Args[0])
	if err != nil {
		return "", err
	}

	answer := fmt.Sprintf(weatherAnswerFormat, resp.Name, resp.Sys.Country, resp.DescriptionTotal(), resp.Main.Temp, resp.Main.Humidity, resp.Main.TempMax, resp.Main.TempMin)

	return answer, nil
}