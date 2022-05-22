package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

type DogHttp struct {
	Message string
	status  string
}

type DogoCommander struct{}

func (c DogoCommander) Handle(s ApiNooter, m Message) {
	n, ok := s.(*DiscordNooter)
	if !ok { return } // Only works on discord

	image := &DogHttp{}
	GetJson("https://dog.ceo/api/breeds/image/random", image)

	body, _ := ReadFile(image.Message)
	imageSend := bytes.NewReader(body)
	fileSend := &discordgo.File{Name: "image.png", ContentType: "image/png", Reader: imageSend}

	response := discordgo.MessageSend{Content: "Random Dogo", Files: []*discordgo.File{fileSend}}

	n.NootComplexMessage(&response)
}

func GetJson(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func ReadFile(URL string) ([]byte, error) {
	//Get the response bytes from the url
	resp, err := http.Get(URL)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return []byte{}, errors.New("Received non 200 response code")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, errors.New("Read all failed")
	}

	return body, nil
}
