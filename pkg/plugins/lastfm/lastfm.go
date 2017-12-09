package lastfm

import (
	"fmt"
	"github.com/disiqueira/MySlackBotTwo/pkg/bot"
	"github.com/disiqueira/MySlackBotTwo/pkg/provider"
)

const (
	answerFormat = "%s is listening to \"%s\" by %s from the album %s."
	answerNotFound = "User not found."
)

var (
	api provider.LastFM
)

func init() {
	bot.RegisterCommand(
		"lastfm",
		"Returns the last music the user was listening.",
		"username",
		lastFM)
}

func getProvider() provider.LastFM {
	if api == nil {
		api = provider.NewLastFM(bot.Configs().LastFMToken())
	}
	return api
}

func lastFM(command *bot.Cmd) (string, error) {
	resp, err := getProvider().ByUser(command.Args[0])
	if err != nil {
		return "", err
	}

	lastTrack, err := resp.LastTrack()
	if err != nil {
		return answerNotFound, nil
	}

	answer := fmt.Sprintf(answerFormat, resp.Recenttracks.Attr.User, lastTrack.Name, lastTrack.Artist.Text, lastTrack.Album.Text)

	return answer, nil
}