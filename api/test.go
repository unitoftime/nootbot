package api

import (
	"fmt"
	"strings"

	"github.com/unitoftime/nootbot/cmd"
)

type Test struct {
	commands []cmd.Command
}
func NewTest(commands []cmd.Command) *Test {
	return &Test{
		commands: commands,
	}
}


func (t *Test) NootMessage(message string) {
	fmt.Println("Sending Message: ", message)
}

func (t *Test) Listen() {
	fmt.Println("Running Test mode!")
	displayMessage := "!echo hello worldddd"

	for _, command := range t.commands {
		prefix, postfix, found := strings.Cut(displayMessage, command.Name)
		if !found {
			continue
		}

		message := cmd.Message{
			Author: cmd.User{
				Id:   "123455",
				Name: "UnitOfTime",
			},
			Parsed: cmd.ParsedMessage{
				Command: command.Name,
				Prefix:  prefix,
				Postfix: postfix,
			},
		}

		command.Handler.Handle(t, message)
	}
}
