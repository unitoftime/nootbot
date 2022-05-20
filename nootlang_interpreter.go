package main

import (
	"fmt"
	"time"

	"github.com/jomy10/nootlang/interpreter"
	"github.com/jomy10/nootlang/parser"
)

type NootlangCommander struct{}

func (c NootlangCommander) Handle(n *Noot, msg Message) {
	str := msg.Parsed.Postfix

	tokens, err := parser.Tokenize(str)
	if err != nil {
		n.SendMessage(fmt.Sprintf("Error while tokenizing source code: %s", err))
	}

	nodes, err := parser.Parse(tokens)
	if err != nil {
		n.SendMessage(fmt.Sprintf("Error while parsing source code: %s", err))
	}

	stdout := make(chan string)
	stderr := make(chan string)
	eop := make(chan int, 1)

	defer close(stdout)
	defer close(stderr)
	defer close(eop)

	go interpreter.Interpret(nodes, stdout, stderr, eop)
	go outHandler(n, stdout)
	go outHandler(n, stderr)

	exitCode := <-eop
	time.Sleep(2 * time.Millisecond)
	fmt.Printf("Noot program exited with exit code %s\n", exitCode)
}

func outHandler(n *Noot, out chan string) {
	n.SendMessage(<-out)
}
