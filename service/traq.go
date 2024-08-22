package service

//go:generate go run github.com/matryer/moq -pkg mock -out mock/${GOFILE} . Traq

import (
	"context"

	"github.com/traP-jp/members_bot/model"
)

type Traq interface {
	GetBotUser(context.Context) (*model.User, error)
	PostMessage(ctx context.Context, channelID, text string) (string, error)
	AddStamp(ctx context.Context, messageID, stampID string, count int) error
	GetGroupMemberIDs(ctx context.Context, groupID string) ([]string, error)
	UpdateUserBio(ctx context.Context, bio string) error
}
