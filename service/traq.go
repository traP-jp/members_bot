package service

//go:generate go run github.com/matryer/moq -pkg mock -out mock/${GOFILE} . Traq

import "context"

type Traq interface {
	GetBotUserID(context.Context) (string, error)
	PostMessage(ctx context.Context, channelID, text string) (string, error)
	AddStamp(ctx context.Context, messageID, stampID string, count int) error
	GetGroupMemberIDs(ctx context.Context, groupID string) ([]string, error)
}
