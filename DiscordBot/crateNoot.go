package main

import (
	"fmt"
	"time"

	"github.com/jomy10/nootlang/interpreter"
	"github.com/jomy10/nootlang/parser"
)

type NootlangCommander struct{}

func (c NootlangCommander) Init() {}

func (c NootlangCommander) Handle(s *Noot, m *Message) {
	str := m.Args[0]

	tokens, err := parser.Tokenize(str)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error while tokenizing source code: %s", err))
	}

	nodes, err := parser.Parse(tokens)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error while parsing source code: %s", err))
	}

	stdout := make(chan string)
	stderr := make(chan string)
	eop := make(chan int, 1)

	defer close(stdout)
	defer close(stderr)
	defer close(eop)

	go interpreter.Interpret(nodes, stdout, stderr, eop)
	go outHandler(s, m, stdout)
	go outHandler(s, m, stderr)

	exitCode := <-eop
	time.Sleep(2 * time.Millisecond)
	fmt.Printf("Noot program exited with exit code %d\n", exitCode)
}

func outHandler(s *Noot, m *Message, out chan string) {
	s.ChannelMessageSend(m.ChannelID, <-out)
}
