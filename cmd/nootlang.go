package cmd

import (
	"bytes"
	"context"
	"fmt"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func(ctx context.Context) {
		if err := interpreter.Interpret(nodes, stdout, stderr, stdin); err != nil {
			n.NootMessage(fmt.Sprintf("[Runtime error] %s", err.Error()))
		}

		// TODO: "realtime" output in discord + stderr
		n.NootMessage(stdout.String())
	}(ctx)

	// Wait for execution or timeout
	<-ctx.Done()

	if ctx.Err() != nil {
		n.NootMessage(ctx.Err().Error())
	}
}
