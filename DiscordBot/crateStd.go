package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type DogHttp struct {
	Message string
	status  string
}

type NootCommander struct{}

func (c NootCommander) Init() {}

func (c NootCommander) Handle(s *Noot, m *Message) {
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Noot Noot! <@%s>", m.Author.ID))
}

type DogoCommander struct{}

func (c DogoCommander) Init() {}

func (c DogoCommander) Handle(s *Noot, m *Message) {
	image := &DogHttp{}
	GetJson("https://dog.ceo/api/breeds/image/random", image)

	body, _ := ReadFile(image.Message)
	imageSend := bytes.NewReader(body)
	fileSend := &discordgo.File{Name: "image.png", ContentType: "image/png", Reader: imageSend}

	response := discordgo.MessageSend{Content: "Random Dogo", Files: []*discordgo.File{fileSend}}

	s.ChannelMessageSendComplex(m.ChannelID, &response)
}

type RandomCommander struct{}

func (c RandomCommander) Init() {}

func (c RandomCommander) Handle(s *Noot, m *Message) {

	min, err := strconv.Atoi(m.Args[0])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Command: %s expects numbers", m.Command))
		return
	}

	max, err := strconv.Atoi(m.Args[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Command: %s expects numbers", m.Command))
		return
	}

	if min > max {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Command: %s expects first arg be smaller number", m.Command))
		return
	}

	if min == 0 || max == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Command: %s args have to be bigger then 0", m.Command))
		return
	}

	value := rand.Intn(max-min) + min

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Random number is %d", value))
}

type InfoCommander struct {
	Content string
}

func (c *InfoCommander) Init() {
	content := ""
	for _, crate := range crates {
		content += fmt.Sprintf("📦%s - %s \n", crate.Name, crate.Author)
		content += "Commands:\n"
		for _, command := range crate.Commands {
			content += fmt.Sprintf("		%s\n", command.Name)
		}
		content += "\n"
	}

	c.Content = content
}

func (c *InfoCommander) Handle(s *Noot, m *Message) {
	embed := &discordgo.MessageEmbed{Type: "rich", Title: "!Info", Description: c.Content, Color: 0x00FFFF}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

type TestCommander struct{}

func (c TestCommander) Init() {}

func (c TestCommander) Handle(s *Noot, m *Message) {
	// {
	// 	"channel_id": `${context.params.event.channel_id}`,
	// 	"content": "",
	// 	"tts": false,
	// 	"embeds": [
	// 		{
	// 			"type": "rich",
	// 			"title": `Hello world`,
	// 			"description": `test test`,
	// 			"color": 0x00FFFF
	// 		}
	// 	]
	// }

	// embed := &discordgo.MessageEmbed{Type: "rich", Title: "hello world", Description: "test test", Color: 0x00FFFF}
	// s.ChannelMessageSendEmbed(m.ChannelID, embed)
	s.ChannelMessageSend(m.ChannelID, "Test!")
}
