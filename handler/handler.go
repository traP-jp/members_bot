package handler

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"text/template"

	"github.com/traP-jp/members_bot/docs"
	"github.com/traP-jp/members_bot/model"
	"github.com/traP-jp/members_bot/repository"
	"github.com/traP-jp/members_bot/service"
)

type BotHandler struct {
	traqClient   service.Traq
	githubClient service.GitHub
	ir           repository.Invitation
	botUser      *model.User
	*Config
}

var logger = log.Default()

func NewBotHandler(traqClient service.Traq, gitHubClient service.GitHub, ir repository.Invitation) (*BotHandler, error) {
	ctx := context.Background()
	botUserID, err := traqClient.GetBotUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot user: %w", err)
	}

	conf, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	logger.SetOutput(traqClient.NewWriter(conf.botChannelID))

	h := &BotHandler{
		traqClient:   traqClient,
		githubClient: gitHubClient,
		ir:           ir,
		botUser:      botUserID,
		Config:       conf,
	}

	helpDoc, err = generateHelpDoc(h)
	if err != nil {
		return nil, fmt.Errorf("failed to generate help doc: %w", err)
	}

	err = h.traqClient.UpdateUserBio(ctx, helpDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to update user bio: %w", err)
	}

	return h, nil
}

func generateHelpDoc(h *BotHandler) (string, error) {
	helpDoc := &strings.Builder{}
	err := template.Must(template.New("help").
		Parse(docs.HelpTemplate)).
		Execute(helpDoc, map[string]string{
			"ORG_NAME":               h.githubClient.OrgName(),
			"BOT_NAME":               h.botUser.Name(),
			"ACCEPT_STAMP_THRESHOLD": strconv.Itoa(h.acceptStampThreshold),
			"REJECT_STAMP_THRESHOLD": strconv.Itoa(h.rejectStampThreshold),
		})

	if err != nil {
		return "", fmt.Errorf("failed to generate help doc: %w", err)
	}

	return helpDoc.String(), nil
}
