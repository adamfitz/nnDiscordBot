package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Auth struct {
	BotToken            string `json:"bot_token"`
	SonarrApiToken      string `json:"sonarr_api_token"`
	Opnsense_api_key    string `json:"opnsense_api_key"`
	Opnsense_api_secret string `json:"opnsense_api_secret"`
}

type Config struct {
	SonarrInstance string `json:"sonarr_instance"`
	SonarrPort     string `json:"sonarr_port"`
	DbServer       string `json:"db_server"`
	DbPort         string `json:"db_port"`
	DbUser         string `json:"db_user"`
	DbPassword     string `json:"db_user_pass"`
	DbName         string `json:"db_name"`
	OpnsenseWanInt string `json:"opnsense_wan_int"`
	OpnsenseFwIp   string `json:"opnsense_fw_ip"`
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

// load config
func LoadConfig() (Config, error) {
	// Get the home directory
	homeDir, err := os.UserHomeDir()
	//construct path to nnDiscoBotConfig
	configPath := fmt.Sprintf("%s/.config/nnDiscordBot", homeDir)
	if err != nil {
		log.Println("Error getting home directory:", err)
		return Config{}, err
	}

	// Build the path to the .nnDiscordCBotConfig file
	configFile := filepath.Join(configPath, "nnDiscordCBot.config")

	// Open the JSON file
	file, err := os.Open(configFile)
	if err != nil {
		log.Println("Error opening file:", err)
		return Config{}, err
	}
	defer file.Close()

	// Decode the JSON file
	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Println("Error decoding JSON:", err)
		return Config{}, err
	}

	// Return the config and no error
	return config, nil
}
