package handler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traP-jp/members_bot/model"
	repomock "github.com/traP-jp/members_bot/repository/mock"
	"github.com/traP-jp/members_bot/service/mock"
	"github.com/traPtitech/traq-ws-bot/payload"
)

func TestInvite(t *testing.T) {
	t.Parallel()

	botUserID := uuid.New().String()
	messageID := uuid.New().String()
	botPostMessageID := uuid.New().String()

	type test struct {
		plainText        string
		messageID        string
		embedded         []payload.EmbeddedInfo
		postTextFunc     func(test) string
		postToBotChannel bool
		invitations      []*model.Invitation
	}

	testCases := map[string]test{
		"特に問題なし": {
			plainText: "@BOT_traP-jp /invite @ikura-hamu ikura-hamu",
			messageID: messageID,
			embedded: []payload.EmbeddedInfo{
				{Type: "user", Raw: "@BOT_traP-jp", ID: botUserID},
				{Type: "user", Raw: "@ikura-hamu", ID: uuid.New().String()},
			},
			postTextFunc: func(t test) string {
				return fmt.Sprintf(`@GitHub_org_Admin
@ikura-hamu https://github.com/ikura-hamu
https://q.trap.jp/messages/%s`, t.messageID)
			},
			postToBotChannel: true,
			invitations:      []*model.Invitation{model.NewInvitation(botPostMessageID, "@ikura-hamu", "ikura-hamu")},
		},
		"「招待」でも問題なし": {
			plainText: "@BOT_traP-jp /招待 @ikura-hamu ikura-hamu",
			messageID: messageID,
			embedded: []payload.EmbeddedInfo{
				{Type: "user", Raw: "@BOT_traP-jp", ID: botUserID},
				{Type: "user", Raw: "@ikura-hamu", ID: uuid.New().String()},
			},
			postTextFunc: func(t test) string {
				return fmt.Sprintf(`@GitHub_org_Admin
@ikura-hamu https://github.com/ikura-hamu
https://q.trap.jp/messages/%s`, t.messageID)
			},
			postToBotChannel: true,
			invitations:      []*model.Invitation{model.NewInvitation(botPostMessageID, "@ikura-hamu", "ikura-hamu")},
		},
		"複数人でも問題なし": {
			plainText: "@BOT_traP-jp /invite @ikura-hamu ikura-hamu @H1rono_K H1rono",
			messageID: messageID,
			embedded: []payload.EmbeddedInfo{
				{Type: "user", Raw: "@BOT_traP-jp", ID: botUserID},
				{Type: "user", Raw: "@ikura-hamu", ID: uuid.New().String()},
				{Type: "user", Raw: "@H1rono_K", ID: uuid.New().String()},
			},
			postTextFunc: func(t test) string {
				return fmt.Sprintf(`@GitHub_org_Admin
@ikura-hamu https://github.com/ikura-hamu
@H1rono_K https://github.com/H1rono
https://q.trap.jp/messages/%s`, t.messageID)
			},
			postToBotChannel: true,
			invitations: []*model.Invitation{
				model.NewInvitation(botPostMessageID, "@ikura-hamu", "ikura-hamu"),
				model.NewInvitation(botPostMessageID, "@H1rono_K", "H1rono"),
			},
		},
		"引数が足りないのでエラー": {
			plainText:    "@BOT_traP-jp /invite",
			messageID:    uuid.New().String(),
			embedded:     []payload.EmbeddedInfo{{Type: "user", Raw: "@BOT_traP-jp", ID: botUserID}},
			postTextFunc: func(test) string { return inviteCommandMessage("引数が足りません") },
		},
		"引数が奇数なのでエラー": {
			plainText: "@BOT_traP-jp /invite @ikura-hamu",
			messageID: uuid.New().String(),
			embedded: []payload.EmbeddedInfo{
				{Type: "user", Raw: "@BOT_traP-jp", ID: botUserID},
				{Type: "user", Raw: "@ikura-hamu", ID: uuid.New().String()},
			},
			postTextFunc: func(test) string { return inviteCommandMessage("引数の数が合いません") },
		},
		"ヘルプ": {
			plainText: "@BOT_traP-jp /invite -h",
			messageID: uuid.New().String(),
			embedded: []payload.EmbeddedInfo{
				{Type: "user", Raw: "@BOT_traP-jp", ID: botUserID},
			},
			postTextFunc: func(test) string {
				return inviteCommandMessage("/invite は、GitHubのOrganizationに招待するためのコマンドです。")
			},
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			traqMock := &mock.TraqMock{
				GetBotUserIDFunc: func(context.Context) (string, error) {
					return botUserID, nil
				},
				PostMessageFunc: func(ctx context.Context, channelID string, text string) (string, error) {
					return botPostMessageID, nil
				},
			}
			repositoryMock := &repomock.InvitationMock{
				CreateInvitationFunc: func(ctx context.Context, invitations []*model.Invitation) error {
					return nil
				},
			}

			bh := &BotHandler{
				traqClient: traqMock,
				ir:         repositoryMock,
				botUserID:  botUserID,
				Config:     &Config{botChannelID: "botChannelID", adminGroupName: "GitHub_org_Admin"},
			}

			payload := &payload.MessageCreated{
				Message: payload.Message{
					PlainText: test.plainText,
					ID:        test.messageID,
					ChannelID: uuid.New().String(),
					Text:      "現時点の実装では使われない",
					Embedded:  test.embedded,
					User:      payload.User{},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Base: payload.Base{EventTime: time.Now()},
			}
			bh.Invite(payload)

			assert.Len(t, traqMock.PostMessageCalls(), 1)
			assert.Equal(t, test.postTextFunc(test), traqMock.PostMessageCalls()[0].Text)
			if test.postToBotChannel {
				assert.Equal(t, "botChannelID", traqMock.PostMessageCalls()[0].ChannelID)
			} else {
				assert.Equal(t, payload.Message.ChannelID, traqMock.PostMessageCalls()[0].ChannelID)
			}

			if test.invitations != nil {
				assert.Len(t, repositoryMock.CreateInvitationCalls(), 1)
				for i, inv := range test.invitations {
					assert.Equal(t, inv, repositoryMock.CreateInvitationCalls()[0].Invitations[i])
				}
			}
		})
	}

}
