package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/NFortun/Astrobot-Sdk/client"
	"github.com/NFortun/Astrobot-Sdk/client/operations"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := run(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
func run() error {
	var Token, ChannelId, host string
	var maxLength int
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&ChannelId, "channel", "", "Channel id")
	flag.IntVar(&maxLength, "max-length", 2000, "max length of body")
	flag.Parse()

	if Token == "" {
		return fmt.Errorf("missing token")
	}

	if ChannelId == "" {
		return fmt.Errorf("missing ChannelId")
	}

	if host = os.Getenv("HOST"); host == "" {
		return fmt.Errorf("missing host")
	}

	if maxLength == 0 {
		return fmt.Errorf("missing max length")
	}

	logrus.Info("Creating new bot")

	// Create a new Discordgo session
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		return err
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	logrus.Info("Opening connection")
	err = dg.Open()
	if err != nil {
		return fmt.Errorf("error while opening connection: %s", err.Error())
	}
	defer dg.Close()

	logrus.Info("retrieving image of the day")
	response, err := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Host: host}).Operations.GetImageOfTheDay(&operations.GetImageOfTheDayParams{
		Context:    context.Background(),
		HTTPClient: http.DefaultClient,
	})
	if err != nil {
		return err
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
		if len(message)+len(field) > maxLength {
			continue
		}
		message += field
	}

	_, err = dg.ChannelMessageSend(ChannelId, message)
	if err != nil {
		return err
	}

	return nil
}
