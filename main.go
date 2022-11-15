package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/unitoftime/nootbot/api"
	"github.com/unitoftime/nootbot/cmd"
)

func main() {
	commands := []cmd.Command{
		cmd.Command{
			Name:        "!echo",
			Description: "[message] - Make NootBot Noot!",
			Handler:     cmd.EchoCommander{},
		},
		cmd.Command{
			Name:        "!recursion",
			Description: "[message] - Make NootBot enter a recursive command.",
			Handler:     cmd.RecursionCommander{},
		},
		cmd.Command{
			Name:        "!eval",
			Description: "[nootlang command] - Evaluate arbitrary nootlang commands.",
			Handler:     cmd.NootlangCommander{},
		},
		cmd.Command{
			Name:        "!java",
			Description: "[None] - Need inspiration for your next Java class?",
			Handler:     cmd.JavaCommander{},
		},
		cmd.Command{
			Name:        "!noot",
			Description: "[None] - Feeling sad? NootBot has a way to make you happy!",
			Handler:     cmd.NootCommander{},
		},
		cmd.Command{
			Name:        "!poll",
			Description: "[question] || [emojisArray] - You have questions, NootBot has answers!",
			Handler:     cmd.PollCommander{},
		},
		cmd.Command{
			Name:        "!weather",
			Description: "[city] | [country code] | [units] - if there are same city names but in different countries, then add a \",\"  after city name in [city] then followed by the country initials for the correct city",
			Handler:     cmd.NewWeatherCommander("weatherApi.token"),
		},
		cmd.Command{
			Name:        "!random",
			Description: "[dog or cat] - Wanna see cute cat or dog image here it is!",
			Handler:     cmd.RandomCommander{},
		},
		cmd.Command{
			Name:        "!notify",
			Description: "[notification] - This command can only be used by the one and only.",
			// For now only non custom emojis are supported
			Handler: cmd.NewNotificationCommander("notify.conf"),
		},
		// cmd.Command{
		// 	Name:    "!random",
		// 	Handler: cmd.RandomCommander{},
		// },
	}

	infoHandler := cmd.NewInfoCommander(commands)
	infoCmd := cmd.Command{
		Name:    "!commands",
		Handler: infoHandler,
	}
	infoCmd2 := cmd.Command{
		Name:    "!info",
		Handler: infoHandler,
	}
	commands = append(commands, infoCmd)
	commands = append(commands, infoCmd2)

	if len(os.Args) < 2 {
		log.Fatal("Must provide URL as first argument")
	}

	if os.Args[1] == "discord" {
		token, err := ioutil.ReadFile("discord.token")
		if err != nil {
			panic(err)
		}
		discord := api.NewDiscord(strings.TrimSuffix(string(token), "\n"), commands)
		discord.Listen()
	} else if os.Args[1] != "test" {
		livestreamId := os.Args[1]
		log.Println(livestreamId)

		ytToken, err := ioutil.ReadFile("token.json")
		if err != nil {
			log.Fatalf("Unable to read client secret file: %v", err)
		}

		youtube := api.NewYoutubeLive(ytToken, commands, livestreamId)
		youtube.Listen()
	} else {
		test := api.NewTest(commands)

		test.NootMessage("Starting Test API")
		test.Listen()
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
