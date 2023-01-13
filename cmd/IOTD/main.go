package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/NFortun/Astrobot-Sdk/client"
	"github.com/NFortun/Astrobot-Sdk/client/operations"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const MAX_LENGTH = 4000

func main() {
	var Token, ChannelId, host string
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&ChannelId, "channel", "", "Channel id")
	flag.Parse()

	if Token == "" {
		log.Fatal("missing token")
	}

	if ChannelId == "" {
		log.Fatal("missing ChannelId")
	}

	if host = os.Getenv("HOST"); host == "" {
		log.Fatal("missing host")
	}

	logrus.Info("Creating new bot")
	// Create a new Discordgo session
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		logrus.Error(err)
		return
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	logrus.Info("Opening connection")
	err = dg.Open()
	if err != nil {
		panic(fmt.Sprintf("error while opening connection: %s", err.Error()))
	}
	defer dg.Close()

	logrus.Info("retrieving image of the day")
	response, err := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Host: host}).Operations.GetImageOfTheDay(&operations.GetImageOfTheDayParams{
		Context:    context.Background(),
		HTTPClient: http.DefaultClient,
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Info("Sending IOTD")
	var message string
	messageMap := map[string]string{
		"Titre":       *response.Payload.Title,
		"Description": *response.Payload.Description,
		"User":        *response.Payload.User,
		"Url":         *response.Payload.URL,
	}

	for k, v := range messageMap {
		field := fmt.Sprintf("%s: %s\n", k, v)
		if len(message)+len(field) > MAX_LENGTH {
			continue
		}
		message += field
	}

	_, err = dg.ChannelMessageSend(ChannelId, message)
	if err != nil {
		logrus.Error(err)
	}

}
