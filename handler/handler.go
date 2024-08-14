package handler

import (
	"context"
	"fmt"

	"github.com/traP-jp/members_bot/repository"
	"github.com/traP-jp/members_bot/service"
)

type BotHandler struct {
	traqClient   service.Traq
	githubClient service.GitHub
	ir           repository.Invitation
	botUserID    string
	*Config
}

func NewBotHandler(traqClient service.Traq, gitHubClient service.GitHub, ir repository.Invitation) (*BotHandler, error) {
	ctx := context.Background()
	botUserID, err := traqClient.GetBotUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot user: %w", err)
	}

	conf, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &BotHandler{
		traqClient:   traqClient,
		githubClient: gitHubClient,
		ir:           ir,
		botUserID:    botUserID,
		Config:       conf,
	}, nil
}
