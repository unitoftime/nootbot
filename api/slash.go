package api

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Slash commands are picked from Discord's command menu instead of typed as chat text, so running one
// posts nothing to the channel and the bot can reply privately. They are discord only, unlike cmd.Command.
type slashCommand struct {
	definition slashCommandDefinition
	handle     func(s *discordgo.Session, i *discordgo.Interaction)
}

// discordgo v0.25.0's ApplicationCommand has no default_member_permissions, so this adds it
type slashCommandDefinition struct {
	*discordgo.ApplicationCommand
	DefaultMemberPermissions string `json:"default_member_permissions,omitempty"`
}

// As DefaultMemberPermissions: Discord only lets server admins see and run the command
const adminOnly = "0"

var slashCommands = []slashCommand{
	exportCommand,
}

// Replaces the bot's registered slash commands with slashCommands. Needs an open session.
func (d *Discord) registerSlashCommands() {
	definitions := make([]slashCommandDefinition, len(slashCommands))
	for i, c := range slashCommands {
		definitions[i] = c.definition
	}

	// A bot's application id is its user id
	endpoint := discordgo.EndpointApplicationGlobalCommands(d.session.State.User.ID)
	_, err := d.session.RequestWithBucketID("PUT", endpoint, definitions, endpoint)
	if err != nil {
		log.Println("Unable to register slash commands:", err)
	}
}

// Handles InteractionCreate events by running the slash command they name
func (d *Discord) handleInteractions(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i == nil || i.Interaction == nil || i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	defer logPanic()

	name := i.ApplicationCommandData().Name
	for _, c := range slashCommands {
		if c.definition.Name != name {
			continue
		}

		user := interactionUser(i.Interaction)
		log.Printf("/%s requested by %s (%s) in channel %s", name, user.Username, user.ID, i.ChannelID)

		c.handle(s, i.Interaction)
		return
	}
}

// Discord fills Member for commands run in a server and User for commands run in DMs
func interactionUser(i *discordgo.Interaction) *discordgo.User {
	if i.Member != nil {
		return i.Member.User
	}
	return i.User
}

// Answers i with a private "thinking..." placeholder. Discord drops interactions that aren't answered within
// 3 seconds, so slow commands defer first and then send their reply with followUpPrivately.
func deferPrivately(s *discordgo.Session, i *discordgo.Interaction) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Flags: uint64(discordgo.MessageFlagsEphemeral)},
	})
}

// Replaces deferPrivately's placeholder with reply, which only the user who ran the command can see
func followUpPrivately(s *discordgo.Session, i *discordgo.Interaction, reply *discordgo.WebhookParams) {
	reply.Flags = uint64(discordgo.MessageFlagsEphemeral)
	if _, err := s.FollowupMessageCreate(i, false, reply); err != nil {
		log.Println("Unable to reply to interaction in channel", i.ChannelID, ":", err)
	}
}
