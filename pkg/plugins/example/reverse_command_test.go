package example

import (
	"github.com/disiqueira/MySlackBotTwo/pkg/bot"
	"testing"
)

func TestReverseString(t *testing.T) {
	arg := "Hello world"
	want := "dlrow olleH"
	bot := &bot.Cmd{
		Command: "reverse",
		RawArgs: arg,
	}

	got, error := reverse(bot)

	if got != want {
		t.Errorf("Expected '%v' got '%v'", want, got)
	}

	if error != nil {
		t.Errorf("Expected '%v' got '%v'", nil, error)
	}
}

func TestReverseEmptyString(t *testing.T) {
	arg := ""
	want := ""
	bot := &bot.Cmd{
		Command: "reverse",
		RawArgs: arg,
	}
	got, error := reverse(bot)

	if got != want {
		t.Errorf("Expected '%v' got '%v'", want, got)
	}

	if error != nil {
		t.Errorf("Expected '%v' got '%v'", nil, error)
	}
}
