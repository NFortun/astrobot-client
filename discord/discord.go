package discord

import "github.com/bwmarrin/discordgo"

type DiscordBot struct {
	session *discordgo.Session
}

func NewDiscordBot(token string) DiscordBot {
	dg, err := discordgo.New(token)
	if err != nil {
		panic(err)
	}

	return DiscordBot{dg}
}
