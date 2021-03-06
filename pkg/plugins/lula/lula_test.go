package lula

import (
	"strings"
	"testing"

	"github.com/disiqueira/MySlackBotTwo/pkg/bot"
)

func TestLulaWhenTheTextDoesNotMatchLula(t *testing.T) {
	cmd := &bot.PassiveCmd{}
	cmd.Raw = "My name is go-bot, I am awesome."
	got, err := lula(cmd)

	if err != nil {
		t.Errorf("Error should be nil => %s", err)
	}
	if got != "" {
		t.Errorf("Test failed. Expected a empty return, got:  '%s'", got)
	}
}

func TestLulaWhenTheTextMatchLula(t *testing.T) {
	cmd := &bot.PassiveCmd{}
	cmd.Raw = "eu não votei na lula!"
	got, err := lula(cmd)

	if err != nil {
		t.Errorf("Error should be nil => %s", err)
	}
	if !strings.HasPrefix(got, ":lula: ") {
		t.Errorf("Test failed. Should return a Lula quote")
	}
}
