package cmd

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type PollCommander struct{}

func (c PollCommander) Handle(s ApiNooter, m Message) {
	n, ok := s.(*DiscordNooter)
	if !ok {
		return
	} // Only works on discord

	// parse message
	split := strings.Split(m.Parsed.Postfix, "|| ")
	emojis := strings.Split(split[1], "")

	// delete author message
	n.NootDeleteMessage(m.Id)

	// generate and send message
	response := discordgo.MessageSend{Content: split[0]}
	msg := n.NootComplexMessage(&response)

	// add emojis
	for _, emoji := range emojis {
		n.NootReact(msg.ID, emoji)
		time.Sleep(5 * time.Millisecond)
	}

}
