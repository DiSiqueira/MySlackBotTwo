package random

import (
	"fmt"
	"github.com/disiqueira/MySlackBotTwo/pkg/bot"
	"math/rand"
	"strconv"
)

const (
	maxValueDefault = 100
)

func init() {
	bot.RegisterCommand(
		"random",
		"Return a random number between zero and the specified value.",
		"{maxValue}",
		random)
}

func random(command *bot.Cmd) (string, error) {
	val := maxValueDefault
	if len(command.Args) > 0 {
		val = maxValue(command.Args[0])
	}
	return fmt.Sprint(rand.Intn(val)), nil
}

func maxValue(text string) int {
	val, err := strconv.Atoi(text)
	if err != nil || val == 0 {
		return maxValueDefault
	}
	return val
}
