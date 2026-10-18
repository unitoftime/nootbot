package cmd

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// This is the general nooter interface. If you want your command to work with all platforms, Noot this way
type ApiNooter interface {
	NootMessage(msg string)
}

// This is a discord specific nooter. If you want to use special discord features, Noot this way
type DiscordNooter struct {
	channel   string
	requester string // Id of the user whose message is being answered
	session   *discordgo.Session
}

func NewDiscordNooter(channel string, requester string, session *discordgo.Session) *DiscordNooter {
	return &DiscordNooter{
		channel:   channel,
		requester: requester,
		session:   session,
	}
}

func (d *DiscordNooter) NootMessage(msg string) {
	d.NootComplexMessage(&discordgo.MessageSend{Content: msg})
}

// Every discord message goes out through here, so commands that repeat user text can't be used
// to post images: link previews are suppressed (attachments, like !random's, still show), only
// user mentions ping, and the message names who it is answering.
func (d *DiscordNooter) NootComplexMessage(complexMessage *discordgo.MessageSend) *discordgo.Message {
	complexMessage.Content = fmt.Sprintf("<@%s> %s", d.requester, complexMessage.Content)
	complexMessage.AllowedMentions = &discordgo.MessageAllowedMentions{
		Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
	}

	m, err := d.send(complexMessage)
	if err != nil {
		log.Println("Unable to send message to channel", d.channel, ":", err)
	}
	return m
}

// discordgo v0.25.0 can't set flags on a new message, so this adds them to its MessageSend
type flaggedMessageSend struct {
	*discordgo.MessageSend
	Flags discordgo.MessageFlags `json:"flags,omitempty"`
}

// Mirrors discordgo's ChannelMessageSendComplex, but sends the message with embeds suppressed
func (d *DiscordNooter) send(complexMessage *discordgo.MessageSend) (*discordgo.Message, error) {
	data := flaggedMessageSend{MessageSend: complexMessage, Flags: discordgo.MessageFlagsSuppressEmbeds}
	endpoint := discordgo.EndpointChannelMessages(d.channel)

	var response []byte
	var err error
	if len(complexMessage.Files) > 0 {
		contentType, body, encodeErr := discordgo.MultipartBodyWithJSON(data, complexMessage.Files)
		if encodeErr != nil {
			return nil, encodeErr
		}
		response, err = d.session.RequestWithLockedBucket("POST", endpoint, contentType, body, d.session.Ratelimiter.LockBucket(endpoint), 0)
	} else {
		response, err = d.session.RequestWithBucketID("POST", endpoint, data, endpoint)
	}
	if err != nil {
		return nil, err
	}

	var m *discordgo.Message
	err = discordgo.Unmarshal(response, &m)
	return m, err
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
	Name        string // The command string to search for
	Description string // The command's description and usage
	Handler     Commander
}

type Commander interface {
	Handle(ApiNooter, Message)
}
