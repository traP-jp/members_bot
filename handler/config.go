package handler

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	botChannelID  string
	acceptStampID string
	rejectStampID string
	adminIDs      []string
}

func loadConfig() (*Config, error) {
	channelID, ok := os.LookupEnv("BOT_CHANNEL_ID")
	if !ok {
		return nil, errors.New("BOT_CHANNEL_ID is not set")
	}

	acceptStampID, ok := os.LookupEnv("ACCEPT_STAMP_ID")
	if !ok {
		return nil, errors.New("ACCEPT_STAMP_ID is not set")
	}

	rejectStampID, ok := os.LookupEnv("REJECT_STAMP_ID")
	if !ok {
		return nil, errors.New("REJECT_STAMP_ID is not set")
	}

	adminIDs, ok := os.LookupEnv("ADMIN_IDS")
	if !ok {
		return nil, errors.New("ADMIN_IDS is not set")
	}

	return &Config{
		botChannelID:  channelID,
		acceptStampID: acceptStampID,
		rejectStampID: rejectStampID,
		adminIDs:      strings.Split(adminIDs, ","),
	}, nil
}
