package cmd

import (
	"fmt"
	"strings"
	// "strconv"
)

type EchoCommander struct {
}

func (c EchoCommander) Handle(n ApiNooter, msg Message) {
	n.NootMessage(msg.Parsed.Prefix + msg.Parsed.Postfix)
}

type RecursionCommander struct {
}

func (c RecursionCommander) Handle(n ApiNooter, msg Message) {
	str := msg.Parsed.Command + " " + msg.Parsed.Prefix + msg.Parsed.Command + msg.Parsed.Postfix

	stackSize := 5

	count := strings.Count(str, msg.Parsed.Command)
	if count > stackSize {
		n.NootMessage("Stack Overflow!")
		return
	}

	n.NootMessage(str)
}

type NootCommander struct{}

func (c NootCommander) Handle(s ApiNooter, m Message) {
	s.NootMessage(fmt.Sprintf("Noot Noot! <@%s>", m.Author.Id))
}

// type RandomCommander struct{}

// func (c RandomCommander) Handle(s ApiNooter, m *Message) {

// 	min, err := strconv.Atoi(m.Args[0])
// 	if err != nil {
// 		s.NootMessage(fmt.Sprintf("Command: %s expects numbers", m.Command))
// 		return
// 	}

// 	max, err := strconv.Atoi(m.Args[1])
// 	if err != nil {
// 		s.NootMessage(fmt.Sprintf("Command: %s expects numbers", m.Command))
// 		return
// 	}

// 	if min > max {
// 		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Command: %s expects first arg be smaller number", m.Command))
// 		return
// 	}

// 	if min == 0 || max == 0 {
// 		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Command: %s args have to be bigger then 0", m.Command))
// 		return
// 	}

// 	value := rand.Intn(max-min) + min

// 	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Random number is %d", value))
// }

type InfoCommander struct {
	Content string
}

func NewInfoCommander(commands []Command) *InfoCommander {
	// content := ""
	// content += fmt.Sprintf("📦%s - %s \n", crate.Name, crate.Author)
	content := "Commands:\n"
	for _, command := range commands {
		content += fmt.Sprintf("		%s %s\n", command.Name, command.Description)
	}
	content += "\n"

	return &InfoCommander{
		Content: content,
	}
}

func (c *InfoCommander) Handle(s ApiNooter, m Message) {
	// TODO add back complex messages
	// embed := &discordgo.MessageEmbed{Type: "rich", Title: "!Info", Description: c.Content, Color: 0x00FFFF}
	// s.ChannelMessageSendEmbed(m.ChannelID, embed)
	s.NootMessage(c.Content)
}
