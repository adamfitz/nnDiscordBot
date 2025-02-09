package bot

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"main/api"
	"main/auth"
	"main/postgres"
	"os"
	"os/signal"
	"strings"
)

// CommandHandler defines the function signature for a command
type CommandHandler func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)

// CommandHandlers holds the mapping of commands to their handlers
var CommandHandlers map[string]CommandHandler

// Declare the variable at the package level
var sonarrLocalSearchUrl string

func Init() {
	// Initialize the CommandHandlers map
	CommandHandlers = map[string]CommandHandler{
		"!help":         handleHelp,
		"!bye":          handleBye,
		"!echo":         handleEcho,
		"!sonarrlookup": handleSonarrSeriesLookup,
		"!sonarrls":     handleSonarrLocalSeriesSearch, // search only the local sonarr instance
		"!dbver":        handleDatabaseVersion,
	}

	// Load the local config file
	config, err := auth.LoadConfig()
	if err != nil {
		log.Println("Error loading config file:", err)
		return
	}

	sonarrLocalSearchUrl = api.ConstructSonarrLocalSeriesURL(config.SonarrInstance, config.SonarrPort)
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

// handleSonarr responds to the !sonarrlookup command and demonstrates argument usage
func handleSonarrSeriesLookup(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Load the credentials from the auth package
	creds, err := auth.LoadCreds()
	if err != nil {
		log.Println("Error loading credentials:", err)
		s.ChannelMessageSend(m.ChannelID, "Error loading credentials. Please try again later.")
		return
	}
	sonarrApiKey := creds.SonarrApiToken

	log.Println("Sonarr Lookup arguments:", args)

	// Check if argument is provided
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !sonarrlookup <series_name>")
		return
	}

	// Call the Sonarr API
	url := sonarrLocalSearchUrl + "/lookup" // search url var declared in the init function

	result, err := api.SonarrSeriesLookupAPICall(url, "X-Api-Key", sonarrApiKey, "term", args[0])
	if err != nil {
		log.Println("Error handling Sonarr API call:", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error handling Sonarr API call: %s", err))
		return
	}

	// Process the Sonarr API response to extract series titles
	seriesTitles, err := api.ProcessSeriesLookupResponse(result)
	if err != nil {
		log.Println("Error processing Sonarr API response:", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error processing Sonarr API response: %s", err))
		return
	}

	log.Println("Series titles:", seriesTitles)

	// Prepare the response message with series titles
	if len(seriesTitles) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No series found.")
		return
	}

	var message string
	for _, title := range seriesTitles {
		message += fmt.Sprintf("- %s\n", title)
	}

	// Function to split message into chunks of 2000 characters or less
	sendMessageChunks := func(message string) {
		for len(message) > 2000 {
			// Find the last line break within 2000 characters
			truncatedMessage := message[:2000]
			lastNewlineIndex := strings.LastIndex(truncatedMessage, "\n")

			if lastNewlineIndex == -1 {
				// If there's no newline in the first 2000 characters, send the whole chunk
				s.ChannelMessageSend(m.ChannelID, truncatedMessage)
				message = message[2000:]
			} else {
				// Send the chunk up to the last complete series
				s.ChannelMessageSend(m.ChannelID, message[:lastNewlineIndex+1])
				message = message[lastNewlineIndex+1:]
			}
		}

		// Send the remaining message (less than 2000 characters)
		if len(message) > 0 {
			s.ChannelMessageSend(m.ChannelID, message)
		}
	}

	// Send the message in chunks if it's too long
	sendMessageChunks(message)
}

// handleSonarr responds to the !sonarrls command
func handleSonarrLocalSeriesSearch(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Load the credentials from the auth package
	creds, err := auth.LoadCreds()
	if err != nil {
		log.Println("Error loading credentials:", err)
		s.ChannelMessageSend(m.ChannelID, "Error loading credentials. Please try again later.")
		return
	}

	sonarrApiKey := creds.SonarrApiToken
	log.Println("Sonarr API Local search arguments:", args)

	// Check if argument is provided
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !sonarrls <series_name>\n(sonarr local search)")
		return
	}

	// Fetch all series from the local Sonarr instance
	allSeriesJSON, err := api.SonarrFetchAllLocalSeries(sonarrLocalSearchUrl, sonarrApiKey)
	if err != nil {
		log.Println("Error fetching series from Sonarr:", err)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error fetching series from Sonarr: %s", err))
		return
	}

	// Unmarshal all series data
	var allSeries []api.Series
	if err := json.Unmarshal([]byte(allSeriesJSON), &allSeries); err != nil {
		log.Println("Error unmarshalling Sonarr API response:", err)
		s.ChannelMessageSend(m.ChannelID, "Error processing Sonarr API response.")
		return
	}

	// Search for series matches
	searchQuery := strings.ToLower(strings.Join(args, " "))
	var matchingSeries []api.Series

	for _, series := range allSeries {
		title := strings.ToLower(series.Title)
		if strings.Contains(title, searchQuery) {
			matchingSeries = append(matchingSeries, series)
		}
	}

	// Prepare the response message
	if len(matchingSeries) == 0 {
		s.ChannelMessageSend(m.ChannelID, "No series found.")
		return
	}

	// Build the response message
	var message string
	for _, series := range matchingSeries {
		message += fmt.Sprintf(
			"**%s** (ID: %d)\n- Seasons: %d\n- Episodes: %d\n- Year: %d\n- Status: %s\n- Genres: %s\n\n",
			series.Title, series.ID, series.SeasonCount, series.EpisodesCount, series.Year,
			series.Status, strings.Join(series.Genres, ", "),
		)
	}

	// Function to split and send message chunks
	sendMessageChunks := func(message string) {
		for len(message) > 2000 {
			// Find the last newline before 2000 characters
			truncatedMessage := message[:2000]
			lastNewlineIndex := strings.LastIndex(truncatedMessage, "\n")

			if lastNewlineIndex == -1 {
				// No newline, send the whole chunk
				s.ChannelMessageSend(m.ChannelID, truncatedMessage)
				message = message[2000:]
			} else {
				// Send up to the last newline
				s.ChannelMessageSend(m.ChannelID, message[:lastNewlineIndex+1])
				message = message[lastNewlineIndex+1:]
			}
		}

		// Send any remaining part of the message
		if len(message) > 0 {
			s.ChannelMessageSend(m.ChannelID, message)
		}
	}

	// Send the response
	sendMessageChunks(message)
}

// database version lookup
func handleDatabaseVersion(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Check if argument is provided
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !dbver - returns the database version)")
		return
	}

	// return the database version
	dbVersion := postgres.DbVersion()
	s.ChannelMessageSend(m.ChannelID, dbVersion)
}
