package api

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/unitoftime/nootbot/cmd"
)

type Discord struct {
	session  *discordgo.Session
	commands []cmd.Command
}

func NewDiscord(token string, commands []cmd.Command) *Discord {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages

	discord := &Discord{
		session:  session,
		commands: commands,
	}
	return discord
}

func (d *Discord) Listen() {
	err := d.session.Open()
	if err != nil {
		panic(err)
	}

	d.session.AddHandler(d.handleMessages)
	d.session.AddHandler(d.handleMessageUpdates)
}

// Handles MessageCreate events
func (d *Discord) handleMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovering from panic: ", r)
			}
		}()

		for _, command := range d.commands {
			prefix, postfix, found := strings.Cut(m.Content, command.Name)
			if !found {
				continue
			}

			message := cmd.Message{
				Id: m.ID,
				Author: cmd.User{
					Id:   m.Author.ID,
					Name: m.Author.Username,
				},
				Parsed: cmd.ParsedMessage{
					Command: command.Name,
					Prefix:  prefix,
					Postfix: postfix,
				},
			}

			nooter := cmd.NewDiscordNooter(m.ChannelID, d.session)
			command.Handler.Handle(nooter, message)
		}
	}()
}
func (d *Discord) handleMessageUpdates(s *discordgo.Session, m *discordgo.MessageUpdate) {
	// Ignore all messages edited by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovering from panic: ", r)
			}
		}()

		for _, command := range d.commands {
			prefix, postfix, found := strings.Cut(m.Content, command.Name)
			if !found {
				continue
			}

			message := cmd.Message{
				Id: m.ID,
				Author: cmd.User{
					Id:   m.Author.ID,
					Name: m.Author.Username,
				},
				Parsed: cmd.ParsedMessage{
					Command: command.Name,
					Prefix:  prefix,
					Postfix: postfix,
				},
			}

			nooter := cmd.NewDiscordNooter(m.ChannelID, d.session)
			command.Handler.Handle(nooter, message)
		}
	}()
}
