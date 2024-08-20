package handler

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traP-jp/members_bot/model"
	repomock "github.com/traP-jp/members_bot/repository/mock"
	"github.com/traP-jp/members_bot/service/mock"
	"github.com/traPtitech/traq-ws-bot/payload"
)

func TestAcceptOrReject(t *testing.T) {

	acceptStampID := uuid.NewString()
	rejectStampID := uuid.NewString()
	inactiveStampID := uuid.NewString()
	adminIDs := []string{uuid.NewString(), uuid.NewString(), uuid.NewString()}
	botUserID := uuid.NewString()

	type testCase struct {
		addStampThreshold        int
		rejectStampThreshold     int
		stamps                   []payload.MessageStamp
		invitations              []*model.Invitation
		executeAddStamp          bool
		executePostMessage       bool
		executeSendInvitations   bool
		postMessageText          string
		executeDeleteInvitations bool
	}
	testCases := map[string]testCase{
		"承認": {
			addStampThreshold:    1,
			rejectStampThreshold: 1,
			stamps: []payload.MessageStamp{
				{StampID: acceptStampID, UserID: adminIDs[0], CreatedAt: time.Now()},
			},
			invitations:              []*model.Invitation{model.NewInvitation(uuid.NewString(), "ikura-hamu", "ikura-hamu")},
			executeAddStamp:          true,
			executePostMessage:       true,
			executeSendInvitations:   true,
			postMessageText:          "招待を送信しました。確認してください\n@ikura-hamu (ikura-hamu)\n",
			executeDeleteInvitations: true,
		},
		"却下": {
			addStampThreshold:    1,
			rejectStampThreshold: 1,
			stamps: []payload.MessageStamp{
				{StampID: rejectStampID, UserID: adminIDs[0], CreatedAt: time.Now()},
			},
			invitations:              []*model.Invitation{model.NewInvitation(uuid.NewString(), "ikura-hamu", "ikura-hamu")},
			executeAddStamp:          true,
			executeDeleteInvitations: true,
		},
		"複数人からの承認": {
			addStampThreshold:    2,
			rejectStampThreshold: 2,
			stamps: []payload.MessageStamp{
				{StampID: acceptStampID, UserID: adminIDs[0], CreatedAt: time.Now()},
				{StampID: acceptStampID, UserID: adminIDs[1], CreatedAt: time.Now()},
			},
			invitations:              []*model.Invitation{model.NewInvitation(uuid.NewString(), "ikura-hamu", "ikura-hamu")},
			executeAddStamp:          true,
			executePostMessage:       true,
			executeSendInvitations:   true,
			postMessageText:          "招待を送信しました。確認してください\n@ikura-hamu (ikura-hamu)\n",
			executeDeleteInvitations: true,
		},
		"複数人からの却下": {
			addStampThreshold:    2,
			rejectStampThreshold: 2,
			stamps: []payload.MessageStamp{
				{StampID: rejectStampID, UserID: adminIDs[0], CreatedAt: time.Now()},
				{StampID: rejectStampID, UserID: adminIDs[1], CreatedAt: time.Now()},
			},
			invitations:              []*model.Invitation{model.NewInvitation(uuid.NewString(), "ikura-hamu", "ikura-hamu")},
			executeAddStamp:          true,
			executeDeleteInvitations: true,
		},
		"承認が先だった": {
			addStampThreshold:    2,
			rejectStampThreshold: 2,
			stamps: []payload.MessageStamp{
				{StampID: acceptStampID, UserID: adminIDs[0], CreatedAt: time.Now().Add(-time.Minute * 2)},
				{StampID: rejectStampID, UserID: adminIDs[1], CreatedAt: time.Now().Add(-time.Minute * 2)},
				{StampID: acceptStampID, UserID: adminIDs[2], CreatedAt: time.Now().Add(-time.Minute)},
				{StampID: rejectStampID, UserID: adminIDs[2], CreatedAt: time.Now()},
			},
			invitations:              []*model.Invitation{model.NewInvitation(uuid.NewString(), "ikura-hamu", "ikura-hamu")},
			executeAddStamp:          true,
			executePostMessage:       true,
			executeSendInvitations:   true,
			postMessageText:          "招待を送信しました。確認してください\n@ikura-hamu (ikura-hamu)\n",
			executeDeleteInvitations: true,
		},
		"却下が先だった": {
			addStampThreshold:    2,
			rejectStampThreshold: 2,
			stamps: []payload.MessageStamp{
				{StampID: acceptStampID, UserID: adminIDs[0], CreatedAt: time.Now().Add(-time.Minute * 2)},
				{StampID: rejectStampID, UserID: adminIDs[1], CreatedAt: time.Now().Add(-time.Minute * 2)},
				{StampID: acceptStampID, UserID: adminIDs[2], CreatedAt: time.Now()},
				{StampID: rejectStampID, UserID: adminIDs[2], CreatedAt: time.Now().Add(-time.Minute)},
			},
			invitations:              []*model.Invitation{model.NewInvitation(uuid.NewString(), "ikura-hamu", "ikura-hamu")},
			executeAddStamp:          true,
			executeDeleteInvitations: true,
		},
		":kan:があるので何もしない": {
			addStampThreshold:    1,
			rejectStampThreshold: 1,
			stamps: []payload.MessageStamp{
				{StampID: inactiveStampID, UserID: botUserID, CreatedAt: time.Now()},
				{StampID: acceptStampID, UserID: adminIDs[0], CreatedAt: time.Now().Add(-time.Minute * 2)},
			},
		},
		"admin以外のスタンプは無視": {
			addStampThreshold:    2,
			rejectStampThreshold: 2,
			stamps: []payload.MessageStamp{
				{StampID: acceptStampID, UserID: adminIDs[0], CreatedAt: time.Now()},
				{StampID: acceptStampID, UserID: uuid.NewString(), CreatedAt: time.Now()},
			},
			invitations: []*model.Invitation{model.NewInvitation(uuid.NewString(), "ikura-hamu", "ikura-hamu")},
		},
		"複数人の招待でも問題なし": {
			addStampThreshold:    1,
			rejectStampThreshold: 1,
			stamps: []payload.MessageStamp{
				{StampID: acceptStampID, UserID: adminIDs[0], CreatedAt: time.Now()},
			},
			invitations: []*model.Invitation{
				model.NewInvitation(uuid.NewString(), "ikura-hamu", "ikura-hamu"),
				model.NewInvitation(uuid.NewString(), "H1rono_K", "H1rono"),
			},
			executeAddStamp:          true,
			executePostMessage:       true,
			executeSendInvitations:   true,
			postMessageText:          "招待を送信しました。確認してください\n@ikura-hamu (ikura-hamu)\n@H1rono_K (H1rono)\n",
			executeDeleteInvitations: true,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			invRepoMock := repomock.InvitationMock{}
			traqMock := mock.TraqMock{}
			gitHubMock := mock.GitHubMock{}

			payload := &payload.BotMessageStampsUpdated{
				MessageID: uuid.New().String(),
				Stamps:    test.stamps,
			}

			bh := &BotHandler{
				traqClient:   &traqMock,
				githubClient: &gitHubMock,
				ir:           &invRepoMock,
				botUserID:    botUserID,
				Config: &Config{
					acceptStampID:        acceptStampID,
					rejectStampID:        rejectStampID,
					inactiveStampID:      inactiveStampID,
					acceptStampThreshold: test.addStampThreshold,
					rejectStampThreshold: test.rejectStampThreshold,
					adminGroupID:         uuid.New().String(),
				},
			}

			traqMock.AddStampFunc = func(context.Context, string, string, int) error {
				return nil
			}
			traqMock.GetGroupMemberIDsFunc = func(context.Context, string) ([]string, error) {
				return adminIDs, nil
			}
			traqMock.PostMessageFunc = func(context.Context, string, string) (string, error) {
				return "", nil
			}

			gitHubMock.SendInvitationsFunc = func(context.Context, []*model.Invitation) error {
				return nil
			}

			invRepoMock.GetInvitationsFunc = func(context.Context, string) ([]*model.Invitation, error) {
				return test.invitations, nil
			}
			invRepoMock.DeleteInvitationsFunc = func(ctx context.Context, invitationID string) error {
				return nil
			}

			bh.AcceptOrReject(payload)

			if test.executeAddStamp {
				assert.Len(t, traqMock.AddStampCalls(), 1)
				assert.Equal(t, payload.MessageID, traqMock.AddStampCalls()[0].MessageID)
				assert.Equal(t, inactiveStampID, traqMock.AddStampCalls()[0].StampID)
				assert.Equal(t, 1, traqMock.AddStampCalls()[0].Count)
			} else {
				assert.Len(t, traqMock.AddStampCalls(), 0)
			}

			if test.executeSendInvitations {
				assert.Len(t, gitHubMock.SendInvitationsCalls(), 1)
				assert.ElementsMatch(t, test.invitations, gitHubMock.SendInvitationsCalls()[0].Invitations)
			} else {
				assert.Len(t, gitHubMock.SendInvitationsCalls(), 0)
			}

			if test.executePostMessage {
				assert.Len(t, traqMock.PostMessageCalls(), 1)
				assert.Equal(t, bh.botChannelID, traqMock.PostMessageCalls()[0].ChannelID)
				assert.Equal(t, test.postMessageText, traqMock.PostMessageCalls()[0].Text)
			} else {
				assert.Len(t, traqMock.PostMessageCalls(), 0)
			}

			if test.executeDeleteInvitations {
				assert.Len(t, invRepoMock.DeleteInvitationsCalls(), 1)
				assert.Equal(t, payload.MessageID, invRepoMock.DeleteInvitationsCalls()[0].InvitationID)
			} else {
				assert.Len(t, invRepoMock.DeleteInvitationsCalls(), 0)
			}

		})
	}
}
