package url

import (
	"github.com/disiqueira/MySlackBotTwo/pkg/bot"
	"github.com/disiqueira/MySlackBotTwo/pkg/plugins/web"
	"html"
	"regexp"
	"strings"
)

var (
	re = regexp.MustCompile("<title>\\n*?(.*?)\\n*?<\\/title>")
)

func urlTitle(cmd *bot.PassiveCmd) (string, error) {
	link := web.ExtractURL(cmd.Raw)
	if link == "" {
		return "", nil
	}

	body, err := web.GetBody(link)
	if err != nil {
		return "", err
	}

	title := re.FindString(string(body))
	if title == "" {
		return "", nil
	}

	title = strings.Replace(title, "\n", "", -1)
	title = title[strings.Index(title, ">")+1 : strings.LastIndex(title, "<")]

	return html.UnescapeString(title), nil
}

func init() {
	bot.RegisterPassiveCommand(
		"url",
		urlTitle)
}
