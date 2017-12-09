package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"github.com/disiqueira/MySlackBotTwo/pkg/config"
	"github.com/disiqueira/MySlackBotTwo/pkg/slack"

	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/9gag"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/catfacts"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/catgif"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/choose"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/chucknorris"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/cmd"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/cnpj"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/cotacao"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/cpf"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/crypto"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/dilma"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/encoding"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/example"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/gif"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/godoc"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/guid"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/image"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/instagram"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/jira"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/lastfm"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/lula"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/megasena"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/puppet"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/random"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/treta"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/uptime"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/url"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/web"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/weather"
	_ "github.com/disiqueira/MySlackBotTwo/pkg/plugins/wolfram"
)

func main() {
	logrus.Info("Starting MySlackBot")

	logrus.Info("Loading configs")
	var cfgs config.Specs
	if err := envconfig.Process("msb", &cfgs); err != nil {
		logrus.Fatal(err.Error())
	}
	logrus.Infof("Configs: %v", cfgs)

	fmt.Println("MySlackBot running!")

	logrus.Info("Starting Slack")

	slack.Run(&cfgs)
}
