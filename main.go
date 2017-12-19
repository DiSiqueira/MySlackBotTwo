package main

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"github.com/disiqueira/MySlackBotTwo/pkg/config"
	"github.com/disiqueira/MySlackBotTwo/pkg/slack"
)

func main() {
	logrus.Info("Starting MySlackBot")

	logrus.Info("Loading configs")
	var cfgs config.Specs
	if err := envconfig.Process("msb", &cfgs); err != nil {
		logrus.Fatal(err.Error())
	}
	logrus.Infof("Configs: %v", cfgs)

	log.SetOutput(os.Stdout)
	logrus.Info("Starting Slack")

	slack.Run(&cfgs)
}
