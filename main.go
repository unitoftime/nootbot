package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"os"
	"os/signal"
	"syscall"

	"github.com/unitoftime/nootbot/api"
	"github.com/unitoftime/nootbot/cmd"
)

func main() {
	commands := []cmd.Command{
		cmd.Command{
			Name:    "!echo",
			Handler: cmd.EchoCommander{},
		},
		cmd.Command{
			Name:    "!recursion",
			Handler: cmd.RecursionCommander{},
		},
		cmd.Command{
			Name:    "!eval",
			Handler: cmd.NootlangCommander{},
		},
		cmd.Command{
			Name:    "!java",
			Handler: cmd.JavaCommander{},
		},
		cmd.Command{
			Name:    "!noot",
			Handler: cmd.NootCommander{},
		},
		cmd.Command{
			Name:    "!dogo",
			Handler: cmd.DogoCommander{},
		},
		// cmd.Command{
		// 	Name:    "!random",
		// 	Handler: cmd.RandomCommander{},
		// },
	}

	infoHandler := cmd.NewInfoCommander(commands)
	infoCmd := cmd.Command{
		Name: "!commands",
		Handler: infoHandler,
	}
	infoCmd2 := cmd.Command{
		Name: "!info",
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
