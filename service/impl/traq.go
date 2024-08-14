package impl

import (
	"context"
	"fmt"

	"github.com/traP-jp/members_bot/service"
	"github.com/traPtitech/go-traq"
)

type Traq struct {
	traqClient *traq.APIClient
}

var _ service.Traq = &Traq{}

func NewTraq(traqClient *traq.APIClient) *Traq {
	return &Traq{traqClient: traqClient}
}

func (t *Traq) GetBotUserID(ctx context.Context) (string, error) {
	me, _, err := t.traqClient.MeApi.GetMe(ctx).Execute()
	if err != nil {
		return "", err
	}

	return me.Id, nil
}

func (t *Traq) PostMessage(ctx context.Context, channelID, text string) (string, error) {
	tr := true
	mes, _, err := t.traqClient.
		MessageApi.PostMessage(ctx, channelID).
		PostMessageRequest(traq.PostMessageRequest{Content: text, Embed: &tr}).
		Execute()
	return mes.Id, err
}

func (t *Traq) AddStamp(ctx context.Context, messageID, stampID string, count int) error {
	for range count {
		_, err := t.traqClient.MessageApi.AddMessageStamp(ctx, messageID, stampID).Execute()
		if err != nil {
			return fmt.Errorf("failed to add stamp: %w", err)
		}
	}
	return nil
}

func (t *Traq) GetGroupMemberIDs(ctx context.Context, groupID string) ([]string, error) {
	members, _, err := t.traqClient.GroupApi.GetUserGroupMembers(ctx, groupID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get group member IDs: %w", err)
	}

	memberIDs := make([]string, 0, len(members))
	for _, member := range members {
		memberIDs = append(memberIDs, member.Id)
	}

	return memberIDs, nil
}
