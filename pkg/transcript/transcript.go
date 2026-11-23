// Package transcript turns a Discord channel's history into plain text, ready to read
package transcript

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const pageSize = 100 // Discord's maximum messages per history request

const timeLayout = "2006-01-02 15:04"

// Fetch returns every message in the channel, oldest first.
func Fetch(s *discordgo.Session, channelID string) ([]*discordgo.Message, error) {
	var msgs []*discordgo.Message
	before := ""
	for {
		// Pages come newest first, walking back from before
		page, err := s.ChannelMessages(channelID, pageSize, before, "", "")
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, page...)
		if len(page) < pageSize {
			break
		}
		before = page[len(page)-1].ID
	}

	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}
	return msgs, nil
}

// Format renders msgs (oldest first) as a transcript titled title. Only messages people wrote are kept:
// joins, pins, renames and other system events are dropped.
func Format(title string, msgs []*discordgo.Message) string {
	kept := make([]*discordgo.Message, 0, len(msgs))
	for _, m := range msgs {
		if m := written(m); m != nil {
			kept = append(kept, m)
		}
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n", title)
	if len(kept) > 0 {
		fmt.Fprintf(&b, "# %d messages, %s to %s UTC\n", len(kept),
			kept[0].Timestamp.UTC().Format(timeLayout), kept[len(kept)-1].Timestamp.UTC().Format(timeLayout))
	} else {
		b.WriteString("# No messages\n")
	}

	for _, m := range kept {
		b.WriteByte('\n')
		writeMessage(&b, m)
	}
	return b.String()
}

// written returns the message a person wrote that m represents, or nil if m is a system event
func written(m *discordgo.Message) *discordgo.Message {
	switch m.Type {
	case discordgo.MessageTypeDefault, discordgo.MessageTypeReply:
		return m
	case discordgo.MessageTypeThreadStarterMessage:
		// A thread started from a channel message opens with a stub pointing at that message
		return m.ReferencedMessage
	}
	return nil
}

// Writes one message as:
//
//	[2026-11-03 14:02] alice (replying to bob): message text,
//	  continued lines indented
//	  attachments: https://...
//	  reactions: 👍 7, :custom: 2
func writeMessage(b *strings.Builder, m *discordgo.Message) {
	fmt.Fprintf(b, "[%s] %s", m.Timestamp.UTC().Format(timeLayout), m.Author.Username)
	if m.Type == discordgo.MessageTypeReply && m.ReferencedMessage != nil && m.ReferencedMessage.Author != nil {
		fmt.Fprintf(b, " (replying to %s)", m.ReferencedMessage.Author.Username)
	}
	b.WriteString(": ")
	b.WriteString(strings.ReplaceAll(m.ContentWithMentionsReplaced(), "\n", "\n  "))
	b.WriteByte('\n')

	if len(m.Attachments) > 0 {
		b.WriteString("  attachments:")
		for i, a := range m.Attachments {
			if i > 0 {
				b.WriteByte(',')
			}
			b.WriteByte(' ')
			b.WriteString(a.URL)
		}
		b.WriteByte('\n')
	}

	if len(m.Reactions) > 0 {
		b.WriteString("  reactions:")
		for i, r := range m.Reactions {
			if i > 0 {
				b.WriteByte(',')
			}
			fmt.Fprintf(b, " %s %d", emojiName(r.Emoji), r.Count)
		}
		b.WriteByte('\n')
	}
}

// Unicode emoji are named by the emoji itself, custom server emoji by their :name:
func emojiName(e *discordgo.Emoji) string {
	if e.ID != "" {
		return ":" + e.Name + ":"
	}
	return e.Name
}
