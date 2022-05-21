package main

import "github.com/bwmarrin/discordgo"

type Crate struct {
	Name      string
	Author    string
	Prefix    string
	ChannelID string
	Commands  []Command
}

type Command struct {
	Name    string
	Args    int8
	Handler Commander
}

type Commander interface {
	Init()
	Handle(s *Noot, m *Message)
}

type Author struct {
	ID       string
	Username string
	Avatar   string
}

type Message struct {
	Author    Author
	Command   string
	Args      []string
	ChannelID string
}

type Noot struct {
	dg *discordgo.Session
}

func (n Noot) ChannelMessageSend(c string, m string) {
	n.dg.ChannelMessageSend(c, m)
}

func (n Noot) ChannelMessageSendComplex(c string, m *discordgo.MessageSend) {
	n.dg.ChannelMessageSendComplex(c, m)
}

func (n Noot) ChannelMessageSendEmbed(c string, m *discordgo.MessageEmbed) {
	n.dg.ChannelMessageSendEmbed(c, m)
}
