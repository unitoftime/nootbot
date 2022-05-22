package cmd

import (
	"github.com/bwmarrin/discordgo"
)

// This is the general nooter interface. If you want your command to work with all platforms, Noot this way
type ApiNooter interface {
	NootMessage(msg string)
}

// This is a discord specific nooter. If you want to use special discord features, Noot this way
type DiscordNooter struct {
	channel string
	session *discordgo.Session
}

func NewDiscordNooter(channel string, session *discordgo.Session) *DiscordNooter {
	return &DiscordNooter{
		channel: channel,
		session: session,
	}
}

func (d *DiscordNooter) NootMessage(msg string) {
	d.session.ChannelMessageSend(d.channel, msg)
}

func (d *DiscordNooter) NootComplexMessage(complexMessage *discordgo.MessageSend) *discordgo.Message {
	m, _ := d.session.ChannelMessageSendComplex(d.channel, complexMessage)
	return m
}

func (d *DiscordNooter) NootDeleteMessage(msgId string) {
	d.session.ChannelMessageDelete(d.channel, msgId)
}

func (d *DiscordNooter) NootReact(msgId string, reaction string) {
	d.session.MessageReactionAdd(d.channel, msgId, reaction)
}

type User struct {
	Id   string // This is the channel id. I think this can't change?
	Name string // This is the users current display name. I think this can change
}

type ParsedMessage struct {
	Command string // This is the command that was detected
	Prefix  string // This is string before the command
	Postfix string // This is the string after the command
}

type Message struct {
	Id     string        // Message id
	Author User          // This is the person who sent the Message
	Parsed ParsedMessage // This is the parsed message
}

type Command struct {
	Name    string // The command string to search for
	Description string // The command's description and usage
	Handler Commander
}

type Commander interface {
	Handle(ApiNooter, Message)
}
