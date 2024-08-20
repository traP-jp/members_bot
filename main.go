package main

import (
	"log"
	"os"

	"github.com/traP-jp/members_bot/handler"
	repoimpl "github.com/traP-jp/members_bot/repository/impl"
	"github.com/traP-jp/members_bot/repository/impl/schema"
	"github.com/traP-jp/members_bot/service/impl"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

func main() {
	botToken, ok := os.LookupEnv("TRAQ_BOT_TOKEN")
	if !ok {
		panic("TRAQ_BOT_TOKEN is not set")
	}

	orgName, ok := os.LookupEnv("GITHUB_ORG_NAME")
	if !ok {
		panic("GITHUB_ORG_NAME is not set")
	}

	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: botToken,
	})
	if err != nil {
		panic(err)
	}

	tc := impl.NewTraq(bot.API())

	gh := impl.NewGitHub(orgName)

	db, err := repoimpl.NewDB()
	if err != nil {
		panic(err)
	}

	schema.Migrate(db)

	ir := repoimpl.NewInvitation(db)

	bh, err := handler.NewBotHandler(tc, gh, ir)
	if err != nil {
		panic(err)
	}

	bot.OnError(func(message string) {
		log.Println("Received ERROR message: " + message)
	})
	bot.OnMessageCreated(bh.Invite)
	bot.OnBotMessageStampsUpdated(bh.AcceptOrReject)

	log.Fatal(bot.Start())
}
