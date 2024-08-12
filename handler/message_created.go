package handler

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"

	"github.com/traP-jp/members_bot/model"
	"github.com/traPtitech/traq-ws-bot/payload"
)

const (
	inviteCommandUsage = `@BOT_traP-jp /(invite|招待) <traQID> <GitHubID> ...`
)

func inviteCommandMessage(message string) string {
	return fmt.Sprintf("%s\n%s", message, inviteCommandUsage)
}

func (h *BotHandler) Invite(p *payload.MessageCreated) {
	ctx := context.Background()

	mentionRawText, ok := checkIfBotMentioned(p, h.botUserID)
	if !ok {
		return
	}

	splitText := regexp.MustCompile(`\s+`).Split(strings.TrimSpace(strings.Replace(p.Message.PlainText, mentionRawText, "", 1)), -1)

	if len(splitText) < 2 {
		err := h.traqClient.PostMessage(ctx, p.Message.ChannelID, inviteCommandMessage("引数が足りません"))
		if err != nil {
			log.Println("failed to post message: ", err)
		}
		return
	}

	if ok, _ := regexp.MatchString(`^/(invite|招待)$`, splitText[0]); !ok {
		return
	}

	splitText = splitText[1:]

	if slices.Contains([]string{"-h", "-help", "--help"}, splitText[0]) {
		err := h.traqClient.PostMessage(ctx, p.Message.ChannelID, inviteCommandMessage("/invite は、GitHubのOrganizationに招待するためのコマンドです。"))
		if err != nil {
			log.Println("failed to post message: ", err)
		}
		return
	}

	if len(splitText)%2 != 0 || len(splitText) == 0 {
		err := h.traqClient.PostMessage(ctx, p.Message.ChannelID, inviteCommandMessage("引数の数が合いません"))
		if err != nil {
			log.Println("failed to post message: ", err)
		}
		return
	}

	invitations := make([]*model.Invitation, 0, len(splitText)/2)
	invitationMessage := "@GitHub_org_Admin\n"
	for i := 0; i < len(splitText); i += 2 {
		traQID := splitText[i]
		gitHubID := splitText[i+1]
		invitations = append(invitations, model.NewInvitation(p.Message.ID, traQID, gitHubID))
		invitationMessage += fmt.Sprintf("%s https://github.com/%s\n", traQID, gitHubID)
	}

	invitationMessage += fmt.Sprintf("https://q.trap.jp/messages/%s", p.Message.ID)

	err := h.ir.CreateInvitation(ctx, invitations)
	if err != nil {
		log.Println("failed to create invitation: ", err)
		return
	}

	err = h.traqClient.PostMessage(ctx, h.botChannelID, invitationMessage)
	if err != nil {
		log.Printf("failed to post message: %v", err)
	}
}

func checkIfBotMentioned(p *payload.MessageCreated, botUserID string) (string, bool) {
	for _, embed := range p.Message.Embedded {
		if embed.Type == "user" && embed.ID == botUserID {

			return embed.Raw, true
		}
	}
	return "", false
}
