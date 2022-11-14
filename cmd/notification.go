package cmd

import (
	"fmt"
	"strings"
)

type NotificationCommander struct {
	notifierId string
	channelId  string
	emoji      string
}

func NewNotificationCommander(notifierId string, channelId string, emoji string) NotificationCommander {
	return NotificationCommander{
		notifierId: notifierId,
		channelId:  channelId,
		emoji:      emoji,
	}
}
func (c NotificationCommander) Handle(s ApiNooter, m Message) {

	d, isDiscordNooter := s.(*DiscordNooter)
	// We only allow this command on discord
	// Additionaly only the notifier can access this command
	if !isDiscordNooter || len(m.Parsed.Postfix) == 0 || m.Author.Id != c.notifierId {
		return
	}

	notificationType := strings.TrimLeft(m.Parsed.Postfix, " ")

	notificationSubscribers, err := d.session.ChannelMessages(c.channelId, 100, "", "", "")
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
			// Get all the people that reacted with a certain emoji
			users, err := d.session.MessageReactions(c.channelId, notification.ID, c.emoji, 100, "", "")
			if err != nil {
				panic("cannot find message reactions")
			}
			for j := range users {
				user := users[j].ID
				if _, ok := subscribers[user]; !ok {
					subscribers[user] = true
				}
			}
			s.NootMessage(formatNootMessage(postfix, subscribers))
			return
		}
	}

}
func formatNootMessage(notification string, subscribers map[string]bool) string {
	message := notification + ""
	for name := range subscribers {
		message += fmt.Sprintf(" <@%s>", name)
	}
	return message
}
