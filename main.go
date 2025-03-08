package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jhenriquem/gongo-wabot/internal/bot"
	"github.com/jhenriquem/gongo-wabot/internal/services"
	"github.com/joho/godotenv"
	"go.mau.fi/whatsmeow"
)

func finish(wac *whatsmeow.Client) {
	wac.Disconnect()
	services.UploadSQLite()
	services.SVC.Database.Close()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Init services
	services.InitDatabase()
	services.InitStorage()

	wac, err := bot.WAConnect()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ONLINE BOT")

	wac.AddEventHandler(bot.GetEventHandler(wac))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	defer finish(wac)
}
