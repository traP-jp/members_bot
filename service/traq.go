package service

//go:generate go run github.com/matryer/moq -pkg mock -out mock/${GOFILE} . Traq

import "context"

type Traq interface {
	GetBotUserID(context.Context) (string, error)
	PostMessage(ctx context.Context, channelID, text string) error
}
