package cmd

type ApiNooter interface {
	NootMessage(msg string)
	// TODO would be good to abstract this away. Right now if you want to send a complex discord message, you should do a type-check in your commander to see if it's a DiscordChannelNooter and then send a complex message that way
	// NootComplexMessage(msg *discordgo.MessageSend)
}

type User struct {
	Id   string // This is the channel id. I think this can't change?
	Name string // This is the users current display name. I think this can change
}

type ParsedMessage struct {
	Command string // This is the command that was detected
	Prefix  string // This is string before the command
	Postfix string // This is the string after the command
}

type Message struct {
	Author User          // This is the person who sent the Message
	Parsed ParsedMessage // This is the parsed message
}

type Command struct {
	Name    string // The command string to search for
	Handler Commander
}

type Commander interface {
	Handle(ApiNooter, Message)
}
