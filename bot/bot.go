package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"main/auth"
	"os"
	"os/signal"
	"strings"
)

// CommandHandler defines the function signature for a command
type CommandHandler func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)

// CommandHandlers holds the mapping of commands to their handlers
var CommandHandlers map[string]CommandHandler

func Init() {
	// Initialize the CommandHandlers map
	CommandHandlers = map[string]CommandHandler{
		"!help": handleHelp,
		"!bye":  handleBye,
		"!echo": handleEcho,
	}
}

func RunBot() {
	// Load credentials
	creds, _ := auth.LoadCreds()

	// create a session
	discordBot, err := discordgo.New("Bot " + creds.BotToken)
	if err != nil {
		log.Fatal("Error message")
	}

	// add a event handler
	discordBot.AddHandler(messageHandler)

	// open session
	discordBot.Open()
	defer discordBot.Close() // close session, after function termination

	// exectuion until os signal interruption (ctrl + C)
	log.Println("nnDiscordBot started....")
	botChannel := make(chan os.Signal, 1)
	signal.Notify(botChannel, os.Interrupt)
	<-botChannel
	log.Println("nnDiscordBot stopped....")

}

// messageHandler processes incoming messages
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if the message starts with a recognized command
	for cmd, handler := range CommandHandlers {
		if strings.HasPrefix(m.Content, cmd) {
			// Split the message to get the command arguments
			args := strings.Fields(m.Content[len(cmd):])
			handler(s, m, args)
			return
		}
	}
}

// handleHelp responds to the !help command
func handleHelp(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	response := "Hello! Here are the available commands:\n"
	for cmd := range CommandHandlers {
		response += fmt.Sprintf("- %s\n", cmd)
	}
	s.ChannelMessageSend(m.ChannelID, response)
}

// handleBye responds to the !bye command
func handleBye(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	s.ChannelMessageSend(m.ChannelID, "Goodbye! 👋")
}

// handleEcho responds to the !echo command and demonstrates argument usage
func handleEcho(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !echo <message>")
		return
	}
	// Join the arguments into a single string
	message := strings.Join(args, " ")
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You said: %s", message))
}
