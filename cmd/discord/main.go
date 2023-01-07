package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NFortun/Astrobot-Sdk/client"
	"github.com/NFortun/Astrobot-Sdk/client/operations"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func main() {
	var Token, channelId, host string
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&channelId, "channel", "", "Channel Id")
	flag.Parse()

	if Token == "" {
		log.Fatal("missing token")
	}

	if channelId == "" {
		log.Fatal("missing ChannelId")
	}

	if host = os.Getenv("HOST"); host == "" {
		log.Fatal("missing host")
	}

	// Create a new Discordgo session
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Println(err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	err = dg.Open()
	if err != nil {
		panic(fmt.Sprintf("error while opening connection: %s", err.Error()))
	}
	defer dg.Close()

	go func() {
		for {
			user := "test"
			response, err := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Host: host}).Operations.GetImages(&operations.GetImagesParams{
				User:       &user,
				Context:    context.Background(),
				HTTPClient: &http.Client{Timeout: 5 * time.Second},
			})
			if err != nil {
				logrus.Error(err)
				continue
			}

			dg.ChannelMessageSend(channelId, *response.Payload[0].URL)
			time.Sleep(5 * time.Second)
		}
	}()

	fmt.Println("Bot is now running, press CTRL C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	s.ChannelMessageSend(m.ChannelID, "test")
}
