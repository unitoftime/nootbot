package main

import (
	"os"
  "fmt"
  "log"
	"time"
	"strings"
  "io/ioutil"


  "golang.org/x/net/context"
  "golang.org/x/oauth2/google"
  "google.golang.org/api/youtube/v3"
)

func main() {
	//"https://youtu.be/qzjGLgezuQk"
	if len(os.Args) < 2 {
		log.Fatal("Must provide URL as first argument")
	}
	livestreamId := os.Args[1]
	log.Println(livestreamId)

  ctx := context.Background()

  b, err := ioutil.ReadFile("token.json")
  if err != nil {
    log.Fatalf("Unable to read client secret file: %v", err)
  }

  // If modifying these scopes, delete your previously saved credentials
  // at ~/.credentials/youtube-go-quickstart.json
  config, err := google.ConfigFromJSON(b, youtube.YoutubeScope)
  if err != nil {
    log.Fatalf("Unable to parse client secret file to config: %v", err)
  }
  client := getClient(ctx, config)
  service, err := youtube.New(client)

  handleError(err, "Error creating YouTube client")

	chatIds := fetchChatIds([]string{livestreamId}, service)

	liveChatId, ok := chatIds[livestreamId]
	if !ok { panic("Failed to get chat ID") }

	fmt.Println("Live chat Id", liveChatId)
	nooter := Noot{
		service: service,
		liveChatId: liveChatId,
	}

	message := "Hello World"
	nooter.SendMessage(message)

	nooter.Listen()
}

type Noot struct {
	service *youtube.Service
	liveChatId string
}

func (n *Noot) SendMessage(message string) {
	call := n.service.LiveChatMessages.Insert([]string{"snippet"}, &youtube.LiveChatMessage{
		Snippet: &youtube.LiveChatMessageSnippet{
			LiveChatId: n.liveChatId,
			Type:       "textMessageEvent",
			TextMessageDetails: &youtube.LiveChatTextMessageDetails{
				MessageText: message,
			},
		},
	})
	_, err := call.Do()
	if err != nil {
		fmt.Println("Error sending message: ", message, " On Channel: ", n.liveChatId, " Error Was: ", err)
	}
}

func (n *Noot) Listen() {
	currentToken := ""
	allPreviousMessagesRead := false

	for {
		call := n.service.LiveChatMessages.List(n.liveChatId, []string{"snippet,authorDetails"})
		if currentToken != "" {
			call.PageToken(currentToken)
		}

		resp, err := call.Do()
		if err != nil {
			fmt.Println("Error with LiveChatMessages.List", err)
			time.Sleep(5 * time.Second)
			continue
		}

		currentToken = resp.NextPageToken

		// Read until the length of the response items is 0, then we know we've read all of the pre-existing messages
		log.Println("Read messages:", len(resp.Items))

		if !allPreviousMessagesRead && len(resp.Items) == 0 {
			log.Println("Finished reading all previous messages")
			allPreviousMessagesRead = true
		}
		if !allPreviousMessagesRead {
			log.Println("Skipping Old messages", len(resp.Items))
			continue
		}

		commands := []Command{
			Command{
				Name: "!echo",
				Handler: EchoCommander{},
			},
			Command{
				Name: "!recursion",
				Handler: RecursionCommander{},
			},
		}


		for _, item := range resp.Items {
			for _, cmd := range commands {
				prefix, postfix, found := strings.Cut(item.Snippet.DisplayMessage, cmd.Name)
				if !found { continue }

				message := Message{
					Author: User{
						Id: item.AuthorDetails.ChannelId,
						Name: item.AuthorDetails.DisplayName,
					},
					Parsed: ParsedMessage{
						Command: cmd.Name,
						Prefix: prefix,
						Postfix: postfix,
					},
				}

				cmd.Handler.Handle(n, message)
			}
		}

		// time.Sleep(time.Duration(resp.PollingIntervalMillis) * time.Millisecond)
		time.Sleep(3 * time.Second)
	}
}

type User struct {
	Id string // This is the channel id. I think this can't change?
	Name string // This is the users current display name. I think this can change
}

type ParsedMessage struct {
	Command string // This is the command that was detected
	Prefix string  // This is string before the command
	Postfix string // This is the string after the command
}

type Message struct {
	Author User // This is the person who sent the Message
	Parsed ParsedMessage // This is the parsed message
}

type Command struct {
	Name string // The command string to search for
	Handler Commander
}

type Commander interface {
	Handle(*Noot, Message)
}

type EchoCommander struct {
}

func (c EchoCommander) Handle(n *Noot, msg Message) {
	n.SendMessage(msg.Parsed.Prefix + msg.Parsed.Postfix)
	// n.SendMessage(fmt.Sprintf("%s %s", msg.Prefix, msg.Postfix)
}

type RecursionCommander struct {
}

func (c RecursionCommander) Handle(n *Noot, msg Message) {
	str := msg.Parsed.Command + " " + msg.Parsed.Prefix + msg.Parsed.Command + msg.Parsed.Postfix

	stackSize := 5

	count := strings.Count(str, msg.Parsed.Command)
	if count > stackSize {
		n.SendMessage("Stack Overflow!")
		return
	}

	n.SendMessage(str)
}
