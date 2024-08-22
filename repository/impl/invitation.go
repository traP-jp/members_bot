package impl

import (
	"context"
	"fmt"

	"github.com/traP-jp/members_bot/model"
	"github.com/traP-jp/members_bot/repository"
	"github.com/traP-jp/members_bot/repository/impl/schema"
	"github.com/uptrace/bun"
)

var _ repository.Invitation = &Invitation{}

type Invitation struct {
	db *bun.DB
}

func NewInvitation(db *bun.DB) *Invitation {
	return &Invitation{db: db}
}

func (i *Invitation) CreateInvitation(ctx context.Context, invitations []*model.Invitation) error {
	invitationSchemes := make([]schema.Invitation, 0, len(invitations))
	for _, invitation := range invitations {
		invitationSchemes = append(invitationSchemes,
			schema.Invitation{
				MessageID: invitation.MessageID(),
				GitHubID:  invitation.GitHubID(),
				TraqID:    invitation.TraqID(),
			})
	}

	_, err := i.db.NewInsert().Model(&invitationSchemes).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create invitation: %w", err)
	}

	return nil
}

func (i *Invitation) GetInvitations(ctx context.Context, id string) ([]*model.Invitation, error) {
	invitationSchemes := make([]*schema.Invitation, 0)
	err := i.db.NewSelect().Model(&invitationSchemes).Where("message_id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get invitations: %w", err)
	}

	if len(invitationSchemes) == 0 {
		return nil, repository.ErrRecordNotFound
	}

	invitations := make([]*model.Invitation, 0, len(invitationSchemes))
	for _, invitationScheme := range invitationSchemes {
		invitations = append(invitations,
			model.NewInvitation(invitationScheme.MessageID, invitationScheme.TraqID, invitationScheme.GitHubID))
	}

	return invitations, nil
}

func (i *Invitation) DeleteInvitations(ctx context.Context, id string) error {
	_, err := i.db.NewDelete().Model(&schema.Invitation{}).Where("message_id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete invitations: %w", err)
	}

	return nil
}

func (i *Invitation) GetAllInvitations(ctx context.Context) ([]*model.Invitation, error) {
	var invitations []schema.Invitation
	err := i.db.NewSelect().Model(&invitations).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get invitations: %w", err)
	}

	invitationsModel := make([]*model.Invitation, 0, len(invitations))
	for _, invitation := range invitations {
		invitationsModel = append(invitationsModel,
			model.NewInvitation(invitation.MessageID, invitation.TraqID, invitation.GitHubID))
	}

	return invitationsModel, nil
}
