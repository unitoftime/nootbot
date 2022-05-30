package cmd

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type CatHttp struct {
	Url string `json:"url"`
	Id  string `json:"id"`
}

type DogHttp struct {
	Message string
	status  string
}

type RandomCommander struct{}

func (c RandomCommander) Handle(s ApiNooter, m Message) {
	n, ok := s.(*DiscordNooter)
	if !ok {
		return
	} // Only works on discord

	response := discordgo.MessageSend{Content: "Random Cat"}
	arg := strings.ReplaceAll(m.Parsed.Postfix, " ", "")
	url := ""
	title := ""

	if arg == "" {
		choose := []string{"dog", "cat"}
		n := rand.Int() % len(choose)
		arg = choose[n]
	}

	switch arg {
	case "dog":
		image := &DogHttp{}
		GetJson("https://dog.ceo/api/breeds/image/random", image)

		url = image.Message
		title = "Random Dog"

	case "cat":
		image := &[]CatHttp{}
		GetJson("https://api.thecatapi.com/v1/images/search", image)

		url = (*image)[0].Url
		title = "Random Cat"
	}

	body, _ := ReadFile(url)
	imageSend := bytes.NewReader(body)
	imageType := strings.Split(url, ".")

	contentType := fmt.Sprintf("image/%s", imageType[len(imageType)-1])
	name := fmt.Sprintf("image.%s", imageType[len(imageType)-1])

	fileSend := &discordgo.File{Name: name, ContentType: contentType, Reader: imageSend}

	response = discordgo.MessageSend{Content: title, Files: []*discordgo.File{fileSend}}

	n.NootComplexMessage(&response)
}
