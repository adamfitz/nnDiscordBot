package main

import (
	"log"
	"main/bot"
	"os"
)

func init() {
	logFile, err := os.OpenFile("nndiscordbot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
	bot.Init()
}

func main() {

	bot.RunBot()

}
