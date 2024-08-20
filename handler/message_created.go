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

func (h *BotHandler) MessageCreated(p *payload.MessageCreated) {
	mentionRawText, isMention := checkIfBotMentioned(p, h.botUser.ID())
	if !isMention {
		return
	}

	splitText := regexp.MustCompile(`\s+`).Split(strings.TrimSpace(strings.Replace(p.Message.PlainText, mentionRawText, "", 1)), -1)

	m := []struct {
		filter func(p *payload.MessageCreated) bool
		fn     func(p *payload.MessageCreated)
	}{
		{
			filter: func(p *payload.MessageCreated) bool {
				ok, _ := regexp.MatchString(`^/(invite|招待)$`, splitText[0])
				return ok
			},
			fn: h.invite,
		},
		{
			filter: func(p *payload.MessageCreated) bool {
				ok, _ := regexp.MatchString(`^/(list|確認)$`, splitText[0])
				return ok
			},
			fn: h.list,
		},
		{
			filter: func(p *payload.MessageCreated) bool {
				ok, _ := regexp.MatchString(`^/(help|ヘルプ|助けて)$`, splitText[0])
				return ok
			},
			fn: h.help,
		},
		{
			filter: func(p *payload.MessageCreated) bool {
				ok, _ := regexp.MatchString(`^/ping$`, splitText[0])
				return ok
			},
			fn: func(p *payload.MessageCreated) {
				ctx := context.Background()
				_, err := h.traqClient.PostMessage(ctx, p.Message.ChannelID, "pong")
				if err != nil {
					log.Println("failed to post message: ", err)
				}
			},
		},
		{
			filter: func(p *payload.MessageCreated) bool {
				ok, _ := regexp.MatchString(`^/`, splitText[0])
				return ok
			},
			fn: func(p *payload.MessageCreated) {
				ctx := context.Background()
				_, err := h.traqClient.PostMessage(ctx, p.Message.ChannelID, ":shiran_zubora.ex-large:")
				if err != nil {
					log.Println("failed to post message: ", err)
				}
			},
		},
	}

	for _, v := range m {
		if v.filter(p) {
			v.fn(p)
			return
		}
	}
}

const (
	inviteCommandUsage = "`@BOT_traP-jp /(invite|招待) <traQID> <GitHubID> ...`"
	listCommandUsage   = "`@BOT_traP-jp /(list|確認)`"
)

func inviteCommandMessage(message string) string {
	return fmt.Sprintf("%s\n%s", message, inviteCommandUsage)
}

func listCommandMessage(message string) string {
	return fmt.Sprintf("%s\n%s", message, listCommandUsage)
}

func (h *BotHandler) invite(p *payload.MessageCreated) {
	ctx := context.Background()

	mentionRawText, _ := checkIfBotMentioned(p, h.botUser.ID())
	splitText := regexp.MustCompile(`\s+`).Split(strings.TrimSpace(strings.Replace(p.Message.PlainText, mentionRawText, "", 1)), -1)

	if len(splitText) < 2 {
		_, err := h.traqClient.PostMessage(ctx, p.Message.ChannelID, inviteCommandMessage("引数が足りません"))
		if err != nil {
			log.Println("failed to post message: ", err)
		}
		return
	}

	splitText = splitText[1:]

	if slices.Contains([]string{"-h", "-help", "--help"}, splitText[0]) {
		_, err := h.traqClient.PostMessage(ctx, p.Message.ChannelID, inviteCommandMessage("/invite は、GitHubのOrganizationに招待するためのコマンドです。"))
		if err != nil {
			log.Println("failed to post message: ", err)
		}
		return
	}

	if len(splitText)%2 != 0 || len(splitText) == 0 {
		_, err := h.traqClient.PostMessage(ctx, p.Message.ChannelID, inviteCommandMessage("引数の数が合いません"))
		if err != nil {
			log.Println("failed to post message: ", err)
		}
		return
	}

	traQIDs := make([]string, 0, len(splitText)/2)
	gitHubIDs := make([]string, 0, len(splitText)/2)
	for i := 0; i < len(splitText); i += 2 {
		traQID := splitText[i]
		gitHubID := splitText[i+1]

		exist, err := h.githubClient.CheckUserExist(ctx, gitHubID)
		if err != nil {
			log.Println("failed to check user exist: ", err)
			return
		}
		if !exist {
			_, err := h.traqClient.PostMessage(ctx, p.Message.ChannelID,
				fmt.Sprintf("GitHubユーザー %s は存在しません", gitHubID))
			if err != nil {
				log.Println("failed to post message: ", err)
			}

			return
		}

		traQIDs = append(traQIDs, traQID)
		gitHubIDs = append(gitHubIDs, gitHubID)
	}

	invitationMessage := fmt.Sprintf("@%s\n", h.adminGroupName)
	for i := range traQIDs {
		invitationMessage += fmt.Sprintf("%s https://github.com/%s\n", traQIDs[i], gitHubIDs[i])
	}
	invitationMessage += fmt.Sprintf("https://q.trap.jp/messages/%s", p.Message.ID)

	messageID, err := h.traqClient.PostMessage(ctx, h.botChannelID, invitationMessage)
	if err != nil {
		log.Printf("failed to post message: %v", err)
	}

	invitations := make([]*model.Invitation, 0, len(splitText)/2)
	for i := range traQIDs {
		invitations = append(invitations, model.NewInvitation(messageID, traQIDs[i], gitHubIDs[i]))
	}

	err = h.ir.CreateInvitation(ctx, invitations)
	if err != nil {
		log.Println("failed to create invitation: ", err)
		return
	}

	err = h.traqClient.AddStamp(ctx, p.Message.ID, h.acceptStampID, 1)
	if err != nil {
		log.Println("failed to add stamp: ", err)
		return
	}
	err = h.traqClient.AddStamp(ctx, p.Message.ID, h.rejectStampID, 1)
	if err != nil {
		log.Println("failed to add stamp: ", err)
		return
	}

}

func (h *BotHandler) list(p *payload.MessageCreated) {
	ctx := context.Background()

	mentionRawText, _ := checkIfBotMentioned(p, h.botUser.ID())
	splitText := regexp.MustCompile(`\s+`).Split(strings.TrimSpace(strings.Replace(p.Message.PlainText, mentionRawText, "", 1)), -1)

	if len(splitText) > 1 && slices.Contains([]string{"-h", "-help", "--help"}, splitText[1]) {
		_, err := h.traqClient.PostMessage(ctx, p.Message.ChannelID,
			listCommandMessage("/list は、招待一覧を表示するためのコマンドです。"))
		if err != nil {
			log.Println("failed to post message: ", err)
		}
		return
	}

	invitations, err := h.ir.GetAllInvitations(ctx)
	if err != nil {
		log.Println("failed to get invitations: ", err)
		return
	}

	message := "招待一覧\n"
	for _, inv := range invitations {
		message += fmt.Sprintf("@%s (%s)\n", inv.TraqID(), inv.GitHubID())
	}

	_, err = h.traqClient.PostMessage(ctx, p.Message.ChannelID, message)
	if err != nil {
		log.Println("failed to post message: ", err)
	}
}

var helpDoc string

func (h *BotHandler) help(p *payload.MessageCreated) {
	ctx := context.Background()

	_, err := h.traqClient.PostMessage(ctx, p.Message.ChannelID, helpDoc)
	if err != nil {
		log.Println("failed to post message: ", err)
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
