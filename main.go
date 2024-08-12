package main

import (
	"log"
	"os"

	"github.com/traP-jp/members_bot/handler"
	"github.com/traP-jp/members_bot/service/impl"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
)

func main() {
	botToken, ok := os.LookupEnv("TRAQ_BOT_TOKEN")
	if !ok {
		panic("TRAQ_BOT_TOKEN is not set")
	}

	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: botToken,
	})
	if err != nil {
		panic(err)
	}

	tc := impl.NewTraq(bot.API())

	botChannelID, ok := os.LookupEnv("BOT_CHANNEL_ID")
	if !ok {
		panic("BOT_CHANNEL_ID is not set")
	}

	bh, err := handler.NewBotHandler(tc, nil, botChannelID)
	if err != nil {
		panic(err)
	}

	bot.OnError(func(message string) {
		log.Println("Received ERROR message: " + message)
	})
	bot.OnMessageCreated(bh.Invite)
	bot.OnMessageCreated(func(p *payload.MessageCreated) {
		log.Println("Received MESSAGE_CREATED event 2: " + p.Message.Text)
	})

	log.Fatal(bot.Start())
}
