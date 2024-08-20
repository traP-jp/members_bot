package handler

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	botChannelID         string
	acceptStampID        string
	acceptStampThreshold int
	rejectStampID        string
	rejectStampThreshold int
	inactiveStampID      string
	adminGroupID         string
	adminGroupName       string
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

	acceptStampThresholdStr, ok := os.LookupEnv("ACCEPT_STAMP_THRESHOLD")
	if !ok {
		return nil, errors.New("ACCEPT_STAMP_THRESHOLD is not set")
	}
	acceptStampThreshold, err := strconv.Atoi(acceptStampThresholdStr)
	if err != nil {
		return nil, errors.New("ACCEPT_STAMP_THRESHOLD is not a number")
	}

	rejectStampID, ok := os.LookupEnv("REJECT_STAMP_ID")
	if !ok {
		return nil, errors.New("REJECT_STAMP_ID is not set")
	}

	rejectStampThresholdStr, ok := os.LookupEnv("REJECT_STAMP_THRESHOLD")
	if !ok {
		return nil, errors.New("REJECT_STAMP_THRESHOLD is not set")
	}
	rejectStampThreshold, err := strconv.Atoi(rejectStampThresholdStr)
	if err != nil {
		return nil, errors.New("REJECT_STAMP_THRESHOLD is not a number")
	}

	inactiveStampID, ok := os.LookupEnv("INACTIVE_STAMP_ID")
	if !ok {
		return nil, errors.New("INACTIVE_STAMP_ID is not set")
	}

	adminGroupID, ok := os.LookupEnv("ADMIN_GROUP_ID")
	if !ok {
		return nil, errors.New("ADMIN_IDS is not set")
	}

	adminGroupName, ok := os.LookupEnv("ADMIN_GROUP_NAME")
	if !ok {
		return nil, errors.New("ADMIN_GROUP_NAME is not set")
	}

	return &Config{
		botChannelID:         channelID,
		acceptStampID:        acceptStampID,
		acceptStampThreshold: acceptStampThreshold,
		rejectStampID:        rejectStampID,
		rejectStampThreshold: rejectStampThreshold,
		inactiveStampID:      inactiveStampID,
		adminGroupID:         adminGroupID,
		adminGroupName:       adminGroupName,
	}, nil
}
