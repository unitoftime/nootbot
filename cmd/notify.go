package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type NotificationCommander struct {
	NotifierId            string
	NotificationChannelId string
	NotificationEmoji     string
}

func NewNotificationCommander(configPath string) NotificationCommander {
	cfg, err := ioutil.ReadFile(configPath)
	if err != nil {

		log.Println("Cannot read notification config file " + err.Error())
	}
	var commander NotificationCommander
	parseError := json.Unmarshal(cfg, &commander)
	if parseError != nil {
		log.Println("Cannot parse notification config " + parseError.Error())
	}
	return commander
}
func (c NotificationCommander) Handle(s ApiNooter, m Message) {
	d, isDiscordNooter := s.(*DiscordNooter)
	// We only allow this command on discord
	// Additionaly only the no(o)tifier can access this command and a notification has to be set
	if !isDiscordNooter || len(m.Parsed.Postfix) == 0 || m.Author.Id != c.NotifierId {
		return
	}
	if c.NotificationChannelId == "" || c.NotificationEmoji == "" || c.NotifierId == "" {
		s.NootMessage("One of the notify configuration fields was not set, please make sure they are set!")
		return
	}
	// Using this we can use notify and still type text afterwards
	notificationType := strings.Split(m.Parsed.Postfix, " ")[1]
	// Retrieve channelmessages of deticated notification channel
	notificationSubscribers, err := d.session.ChannelMessages(c.NotificationChannelId, 100, "", "", "")
	if err != nil {
		panic("Cannot find channel id")
	}

	subscribers := make(map[string]bool)

	for i := range notificationSubscribers {

		notification := notificationSubscribers[i]
		_, postfix, found := strings.Cut(notification.Content, fmt.Sprintf("[%s]", notificationType))
		if !found || postfix == "" {
			continue
		} else {
			// Get all the discord users that reacted with a certain emoji
			// Currently only unicode emojis work using this call
			users, err := d.session.MessageReactions(c.NotificationChannelId, notification.ID, c.NotificationEmoji, 100, "", "")
			if err != nil {
				log.Println("Could not find reaction to message")
			}
			for j := range users {
				user := users[j].ID
				if _, ok := subscribers[user]; !ok {
					subscribers[user] = true
				}
			}
			_, msgText, _ := strings.Cut(strings.Trim(m.Parsed.Postfix, " "), " ")
			s.NootMessage(formatNootMessage(fmt.Sprintf("[%s] %s", notificationType, msgText), subscribers))
			return
		}
	}

}
func formatNootMessage(notification string, subscribers map[string]bool) string {
	for name := range subscribers {
		notification += fmt.Sprintf(" <@%s>", name)
	}
	return notification
}
