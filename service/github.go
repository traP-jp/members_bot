package service

//go:generate go run github.com/matryer/moq -pkg mock -out mock/${GOFILE} . GitHub

import (
	"context"

	"github.com/traP-jp/members_bot/model"
)

type GitHub interface {
	SendInvitations(ctx context.Context, invitations []*model.Invitation) error
	CheckUserExist(ctx context.Context, userID string) (bool, error)
}
