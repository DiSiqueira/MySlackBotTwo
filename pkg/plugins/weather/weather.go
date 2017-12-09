package weather

import (
	"fmt"
	"github.com/disiqueira/MySlackBotTwo/pkg/bot"
	"github.com/disiqueira/MySlackBotTwo/pkg/provider"
	"math/rand"
	"regexp"
	"strconv"
)

var (
	re  = regexp.MustCompile(pattern)
	api provider.Weather
)

const (
	pattern          = "(?i)\\b(pokemon|poke|pikachu)[s|z]{0,1}\\b"
	maxPokemon       = 500
	pokeAnswerFormat = "%s (%d) %s"
)

func pokemon(command *bot.PassiveCmd) (string, error) {
	if !re.MatchString(command.Raw) {
		return "", nil
	}
	poke, err := api.Search(strconv.Itoa(rand.Intn(maxPokemon)))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(pokeAnswerFormat, poke.Name, poke.ID, poke.Sprites.FrontDefault), nil
}

func init() {
	api = provider.NewWeather()
	bot.RegisterPassiveCommand(
		"pokemon",
		pokemon)
}
