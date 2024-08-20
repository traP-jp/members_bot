package repository

//go:generate go run github.com/matryer/moq -pkg mock -out mock/${GOFILE} . Invitation

import (
	"context"

	"github.com/traP-jp/members_bot/model"
)

type Invitation interface {
	CreateInvitation(ctx context.Context, invitations []*model.Invitation) error
	GetInvitations(ctx context.Context, invitationID string) ([]*model.Invitation, error)
	DeleteInvitations(ctx context.Context, invitationID string) error
}
