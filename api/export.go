package api

import (
	"log"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"

	"github.com/unitoftime/nootbot/pkg/transcript"
)

// /export privately sends the admin who runs it the thread they're in as a text file
var exportCommand = slashCommand{
	definition: slashCommandDefinition{
		ApplicationCommand: &discordgo.ApplicationCommand{
			Name:        "export",
			Description: "Privately sends you this thread as a text file",
		},
		DefaultMemberPermissions: adminOnly,
	},
	handle: handleExport,
}

func handleExport(s *discordgo.Session, i *discordgo.Interaction) {
	if err := deferPrivately(s, i); err != nil {
		log.Println("Unable to acknowledge /export in channel", i.ChannelID, ":", err)
		return
	}
	followUpPrivately(s, i, exportThread(s, i.ChannelID))
}

// Returns the reply to /export run in channelID: the thread's transcript, or why there isn't one
func exportThread(s *discordgo.Session, channelID string) *discordgo.WebhookParams {
	channel, err := s.Channel(channelID)
	if err != nil {
		return &discordgo.WebhookParams{Content: "Couldn't look up this channel: " + err.Error()}
	}
	// Limited to threads so a misplaced run can't page through a channel's entire history
	if !channel.IsThread() {
		return &discordgo.WebhookParams{Content: "/export only works inside a thread."}
	}

	msgs, err := transcript.Fetch(s, channelID)
	if err != nil {
		return &discordgo.WebhookParams{Content: "Couldn't read this thread's messages: " + err.Error()}
	}

	return &discordgo.WebhookParams{
		Content: "Here's **" + channel.Name + "**.",
		Files: []*discordgo.File{{
			Name:        fileName(channel.Name),
			ContentType: "text/plain",
			Reader:      strings.NewReader(transcript.Format(channel.Name, msgs)),
		}},
	}
}

// Turns a thread name like "v1.23 Bugs and Feedback" into "v1.23-Bugs-and-Feedback.txt"
func fileName(threadName string) string {
	name := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '-' || r == '_' {
			return r
		}
		return '-'
	}, threadName)
	return name + ".txt"
}
