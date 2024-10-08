package impl

import (
	"context"
	"fmt"
	"io"

	"github.com/traP-jp/members_bot/model"
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

func (t *Traq) GetBotUser(ctx context.Context) (*model.User, error) {
	me, _, err := t.traqClient.MeApi.GetOIDCUserInfo(ctx).Execute()
	if err != nil {
		return nil, err
	}

	return model.NewUser(me.Sub, me.Name), nil
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

func (t *Traq) UpdateUserBio(ctx context.Context, bio string) error {
	_, err := t.traqClient.MeApi.
		EditMe(ctx).PatchMeRequest(traq.PatchMeRequest{Bio: &bio}).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to update user bio: %w", err)
	}

	return nil
}

func (t *Traq) NewWriter(channelID string) io.Writer {
	return &TraqWriter{channelID: channelID, t: t}
}

var _ io.Writer = (*TraqWriter)(nil)

type TraqWriter struct {
	channelID string
	t         *Traq
}

func (tw *TraqWriter) Write(p []byte) (n int, err error) {
	_, err = tw.t.PostMessage(context.Background(), tw.channelID, string(p))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
