package api

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/unitoftime/nootbot/pkg/live"

	"github.com/bwmarrin/discordgo"

	"github.com/unitoftime/nootbot/cmd"
)

type Discord struct {
	session  *discordgo.Session
	commands []cmd.Command

	liveBanners *live.BannerSystem
}

func NewDiscord(token string, commands []cmd.Command) *Discord {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	// Message content is needed to read commands and to export threads
	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent

	discord := &Discord{
		liveBanners: live.NewBannerSystem(session, live.LiveBanners, 5 * time.Minute),
		session:     session,
		commands:    commands,
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
	d.session.AddHandler(d.handleInteractions)
	d.registerSlashCommands()

	// goroutine because the scheduler will block once started
	go d.liveBanners.Listen()
}

// Handles MessageCreate events
func (d *Discord) handleMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	d.dispatch(s, m.Message, "new")
}

// Handles MessageUpdate events, so edited messages run commands too
func (d *Discord) handleMessageUpdates(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m == nil {
		return
	}
	d.dispatch(s, m.Message, "edited")
}

// Runs every command found in the message. kind says whether the message is new or edited.
func (d *Discord) dispatch(s *discordgo.Session, m *discordgo.Message, kind string) {
	// Ignore if the session or the message are nil
	if s == nil || m == nil || m.Author == nil || s.State == nil || s.State.User == nil {
		return
	}
	// Ignore all messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	defer logPanic()

	for _, command := range d.commands {
		prefix, postfix, found := strings.Cut(m.Content, command.Name)
		if !found {
			continue
		}

		// Logged before running, so the request is on record even if the author deletes it
		log.Printf("%s requested by %s (%s) in channel %s, %s message %s: %q",
			command.Name, m.Author.Username, m.Author.ID, m.ChannelID, kind, m.ID, m.Content)

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

		nooter := cmd.NewDiscordNooter(m.ChannelID, m.Author.ID, d.session)
		command.Handler.Handle(nooter, message)
	}
}

// Deferred by event handlers, so one bad command can't take down the bot
func logPanic() {
	if r := recover(); r != nil {
		fmt.Println("Recovering from panic: ", r)
	}
}
