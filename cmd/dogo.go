package cmd

import (
	"bytes"

	"github.com/bwmarrin/discordgo"
)

type DogHttp struct {
	Message string
	status  string
}

type DogoCommander struct{}

func (c DogoCommander) Handle(s ApiNooter, m Message) {
	n, ok := s.(*DiscordNooter)
	if !ok {
		return
	} // Only works on discord

	image := &DogHttp{}
	GetJson("https://dog.ceo/api/breeds/image/random", image)

	body, _ := ReadFile(image.Message)
	imageSend := bytes.NewReader(body)
	fileSend := &discordgo.File{Name: "image.png", ContentType: "image/png", Reader: imageSend}

	response := discordgo.MessageSend{Content: "Random Dogo", Files: []*discordgo.File{fileSend}}

	n.NootComplexMessage(&response)
}
