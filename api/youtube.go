package api

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"

	"github.com/unitoftime/nootbot/cmd"
)

type YoutubeLive struct {
	service    *youtube.Service
	liveChatId string
	commands   []cmd.Command
}

// TODO dynamically find livestream ID
func NewYoutubeLive(token []byte, commands []cmd.Command, livestreamId string) *YoutubeLive {
	ctx := context.Background()
	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(token, youtube.YoutubeScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)
	service, err := youtube.New(client)

	handleError(err, "Error creating YouTube client")

	chatIds := fetchChatIds([]string{livestreamId}, service)

	liveChatId, ok := chatIds[livestreamId]
	if !ok {
		panic("Failed to get chat ID")
	}

	fmt.Println("Live chat Id", liveChatId)
	ytLive := YoutubeLive{
		service:    service,
		liveChatId: liveChatId,
		commands:   commands,
	}

	message := "Hello World"
	ytLive.NootMessage(message)

	return &ytLive
}

func (yt *YoutubeLive) NootMessage(message string) {
	call := yt.service.LiveChatMessages.Insert([]string{"snippet"}, &youtube.LiveChatMessage{
		Snippet: &youtube.LiveChatMessageSnippet{
			LiveChatId: yt.liveChatId,
			Type:       "textMessageEvent",
			TextMessageDetails: &youtube.LiveChatTextMessageDetails{
				MessageText: message,
			},
		},
	})
	_, err := call.Do()
	if err != nil {
		fmt.Println("Error sending message: ", message, " On Channel: ", yt.liveChatId, " Error Was: ", err)
	}
}

func (yt *YoutubeLive) Listen() {
	currentToken := ""
	allPreviousMessagesRead := false

	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovering from panic: ", r)
				}
			}()
			call := yt.service.LiveChatMessages.List(yt.liveChatId, []string{"snippet,authorDetails"})
			if currentToken != "" {
				call.PageToken(currentToken)
			}

			resp, err := call.Do()
			if err != nil {
				fmt.Println("Error with LiveChatMessages.List", err)
				time.Sleep(5 * time.Second)
				return
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
				return
			}

			for _, item := range resp.Items {
				for _, command := range yt.commands {
					prefix, postfix, found := strings.Cut(item.Snippet.DisplayMessage, command.Name)
					if !found {
						continue
					}

					message := cmd.Message{
						Author: cmd.User{
							Id:   item.AuthorDetails.ChannelId,
							Name: item.AuthorDetails.DisplayName,
						},
						Parsed: cmd.ParsedMessage{
							Command: command.Name,
							Prefix:  prefix,
							Postfix: postfix,
						},
					}

					command.Handler.Handle(yt, message)
				}
			}

			// time.Sleep(time.Duration(resp.PollingIntervalMillis) * time.Millisecond)
			time.Sleep(5 * time.Second)
		}()
	}
}
