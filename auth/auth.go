package auth

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Auth struct {
	BotToken string `json:"bot_token"`
}

// LoadCreds loads the credentials from the .discordrc file.
func LoadCreds() (Auth, error) {
	// Get the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("Error getting home directory:", err)
		return Auth{}, err
	}

	// Build the path to the .discordrc file
	configPath := filepath.Join(homeDir, ".discordrc")

	// Open the JSON file
	file, err := os.Open(configPath)
	if err != nil {
		log.Println("Error opening file:", err)
		return Auth{}, err
	}
	defer file.Close()

	// Decode the JSON file
	var config Auth
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Println("Error decoding JSON:", err)
		return Auth{}, err
	}

	// Return the config and no error
	return config, nil
}
