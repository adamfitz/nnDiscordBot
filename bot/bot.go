package bot

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"main/auth"
	"os"
	"os/signal"
	"strings"
)

//var BotToken = make(map[string]string)

func RunBot() {
	// Load credentials
	creds, _ := auth.LoadCreds()

	// create a session
	discordBot, err := discordgo.New("Bot " + creds.BotToken)
	if err != nil {
		log.Fatal("Error message")
	}

	// add a event handler
	discordBot.AddHandler(newMessage)

	// open session
	discordBot.Open()
	defer discordBot.Close() // close session, after function termination

	// exectuion until os signal interruption (ctrl + C)
	log.Println("nnDiscordBot running....")
	botChannel := make(chan os.Signal, 1)
	signal.Notify(botChannel, os.Interrupt)
	<-botChannel

}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	/* prevent bot responding to its own message
	this is achived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// respond to user message if it contains `!help` or `!bye`
	switch {
	case strings.Contains(message.Content, "!help"):
		discord.ChannelMessageSend(message.ChannelID, "Hello World😃")
	case strings.Contains(message.Content, "!bye"):
		discord.ChannelMessageSend(message.ChannelID, "Good Bye👋")
		// add more cases if required
	}
}
