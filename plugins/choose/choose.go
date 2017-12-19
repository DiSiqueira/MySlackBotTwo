package main

import (
	"math/rand"
)

type choose struct {

}

func (c *choose) Init() {
}

func (c *choose) Command() string {
	return "choose"
}

func (c *choose) Description() string {
	return "Randomly picks one of the options."
}

func (c *choose) ExampleArgs() string {
	return "option1 option2 {option3 ...}"
}

func (c *choose) Execute(args []string,channel string) (string, string, error) {
	return args[rand.Intn(len(args))],channel, nil
}

var CustomPlugin choose
