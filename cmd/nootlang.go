package cmd

import (
	"bytes"
	"fmt"
	// "time"

	"github.com/jomy10/nootlang/interpreter"
	"github.com/jomy10/nootlang/parser"
)

type NootlangCommander struct{}

func (c NootlangCommander) Handle(n ApiNooter, msg Message) {
	str := msg.Parsed.Postfix

	tokens, err := parser.Tokenize(str)
	if err != nil {
		n.NootMessage(fmt.Sprintf("Error while tokenizing source code: %s", err))
		return
	}

	nodes, err := parser.Parse(tokens)
	if err != nil {
		n.NootMessage(fmt.Sprintf("Error while parsing source code: %s", err))
		return
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	stdin := new(bytes.Reader)

	if err := interpreter.Interpret(nodes, stdout, stderr, stdin); err != nil {
		n.NootMessage(fmt.Sprintf("[Runtime error] %s", err.Error()))
	}

	// TODO: "realtime" output in discord + stderr
	n.NootMessage(stdout.String())
}
