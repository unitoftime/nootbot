package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var noot = &Noot{}
var crates = []Crate{
	{
		Name:      "Std",
		Author:    "Dracula",
		Prefix:    "!",
		ChannelID: "661170976528007189",
		Commands: []Command{
			{Name: "Info", Args: 0, Handler: &InfoCommander{}},
			{Name: "Noot", Args: 0, Handler: NootCommander{}},
			{Name: "Dogo", Args: 0, Handler: DogoCommander{}},
			{Name: "Random", Args: 2, Handler: RandomCommander{}},
			{Name: "Test", Args: 0, Handler: TestCommander{}},
		},
	},
	{
		Name:      "NootLang",
		Author:    "Jomy",
		Prefix:    "!",
		ChannelID: "661170976528007189",
		Commands: []Command{
			{Name: "Eval", Args: -1, Handler: NootlangCommander{}},
		},
	},
}

// load token and other variables from env file
func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Init all commands
	for _, crate := range crates {
		for _, command := range crate.Commands {
			command.Handler.Init()
		}
	}
}

func main() {
	Init()
	token := os.Getenv("TOKEN_DISCORD")

	dg, err := discordgo.New("Bot " + token)

	//setup noot
	noot.dg = dg

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	message := &Message{
		Command: "", Args: []string{}, ChannelID: m.ChannelID,
		Author: Author{ID: m.Author.ID, Username: m.Author.Username, Avatar: m.Author.Avatar},
	}

	for _, crate := range crates {
		if crate.ChannelID != m.ChannelID {
			continue
		}

		for _, command := range crate.Commands {
			commandName := crate.Prefix + command.Name
			if strings.HasPrefix(m.Content, commandName) {

				// generate args
				if command.Args > 0 {
					params := strings.Split(m.Content, " ")
					if (len(params) - 1) < int(command.Args) {
						s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Command: %s expects %d args", commandName, command.Args))
						continue
					}
					message.Args = params[1:]
				} else if command.Args == -1 {
					arg := strings.Split(m.Content, commandName)
					message.Args = []string{arg[1]}
				}

				message.Command = commandName

				// then do Handler
				command.Handler.Handle(noot, message)
			}
		}
	}
}
