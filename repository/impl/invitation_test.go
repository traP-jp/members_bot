package impl

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traP-jp/members_bot/model"
	"github.com/traP-jp/members_bot/repository"
	"github.com/traP-jp/members_bot/repository/impl/schema"
)

func TestCreateInvitation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		t.Cleanup(func() {
			_, err := testDB.NewTruncateTable().Model(&schema.Invitation{}).Exec(ctx)
			require.NoError(t, err)
		})

		ir := NewInvitation(testDB)

		invitations := []*model.Invitation{
			model.NewInvitation(uuid.NewString(), "github_id", "traq_id"),
		}

		err := ir.CreateInvitation(ctx, invitations)
		assert.NoError(t, err)

		var invitationsTable []schema.Invitation
		err = ir.db.NewSelect().Model(&invitationsTable).Scan(ctx)
		require.NoError(t, err)

		assert.Len(t, invitationsTable, len(invitations))
		for i, invitation := range invitationsTable {
			assert.Equal(t, invitations[i].MessageID(), invitation.MessageID)
			assert.Equal(t, invitations[i].GitHubID(), invitation.GitHubID)
			assert.Equal(t, invitations[i].TraqID(), invitation.TraqID)
			assert.WithinDuration(t, time.Now(), invitation.CreatedAt, time.Second)
		}
	})
}

func TestGetInvitations(t *testing.T) {
	type testCase struct {
		invitationID string
		fixture      []*schema.Invitation
		expected     []*model.Invitation
		expectedErr  error
	}

	t.Cleanup(func() {
		_, err := testDB.NewTruncateTable().Model(&schema.Invitation{}).Exec(context.Background())
		require.NoError(t, err)
	})

	invitationID1 := uuid.NewString()
	invitationID2 := uuid.NewString()

	testCases := map[string]testCase{
		"特に問題なし": {
			invitationID: invitationID1,
			fixture: []*schema.Invitation{
				{MessageID: invitationID1, GitHubID: "github_id", TraqID: "traq_id"},
			},
			expected: []*model.Invitation{
				model.NewInvitation(invitationID1, "traq_id", "github_id"),
			},
		},
		"他の招待があっても問題なし": {
			invitationID: invitationID2,
			fixture: []*schema.Invitation{
				{MessageID: invitationID2, GitHubID: "github_id2", TraqID: "traq_id2"},
				{MessageID: uuid.NewString(), GitHubID: "github_id3", TraqID: "traq_id3"},
			},
			expected: []*model.Invitation{
				model.NewInvitation(invitationID2, "traq_id2", "github_id2"),
			},
		},
		"招待がない": {
			invitationID: uuid.NewString(),
			fixture:      []*schema.Invitation{},
			expected:     []*model.Invitation{},
			expectedErr:  repository.ErrRecordNotFound,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			ir := NewInvitation(testDB)

			{
				if len(test.fixture) != 0 {
					_, err := ir.db.NewInsert().Model(&test.fixture).Exec(ctx)
					require.NoError(t, err)
				}
			}

			invitations, err := ir.GetInvitations(ctx, test.invitationID)

			assert.Len(t, invitations, len(test.expected))
			for i, invitation := range invitations {
				assert.Equal(t, test.expected[i].MessageID(), invitation.MessageID())
				assert.Equal(t, test.expected[i].GitHubID(), invitation.GitHubID())
				assert.Equal(t, test.expected[i].TraqID(), invitation.TraqID())
			}

			assert.ErrorIs(t, err, test.expectedErr)
		})
	}
}

func TestDeleteInvitations(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		t.Cleanup(func() {
			_, err := testDB.NewTruncateTable().Model(&schema.Invitation{}).Exec(ctx)
			require.NoError(t, err)
		})

		ir := NewInvitation(testDB)

		invitationID := uuid.NewString()
		{
			_, err := ir.db.NewInsert().Model(&schema.Invitation{MessageID: invitationID}).Exec(ctx)
			require.NoError(t, err)
		}

		err := ir.DeleteInvitations(ctx, invitationID)
		assert.NoError(t, err)

		var invitations []schema.Invitation
		err = ir.db.NewSelect().Model(&invitations).Scan(ctx)
		require.NoError(t, err)

		assert.Len(t, invitations, 0)
	})
}

func TestGetAllInvitations(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		t.Cleanup(func() {
			_, err := testDB.NewTruncateTable().Model(&schema.Invitation{}).Exec(ctx)
			require.NoError(t, err)
		})

		ir := NewInvitation(testDB)

		invitations := []*model.Invitation{
			model.NewInvitation(uuid.NewString(), "github_id", "traq_id"),
			model.NewInvitation(uuid.NewString(), "github_id2", "traq_id2"),
		}
		{
			invitationsTable := make([]schema.Invitation, 0, len(invitations))
			for _, invitation := range invitations {
				invitationsTable = append(invitationsTable, schema.Invitation{
					MessageID: invitation.MessageID(),
					GitHubID:  invitation.GitHubID(),
					TraqID:    invitation.TraqID(),
				})
			}

			_, err := ir.db.NewInsert().Model(&invitationsTable).Exec(ctx)
			require.NoError(t, err)
		}

		result, err := ir.GetAllInvitations(ctx)
		require.NoError(t, err)

		assert.Len(t, result, len(invitations))

		for i, invitation := range result {
			assert.Equal(t, invitations[i].MessageID(), invitation.MessageID())
			assert.Equal(t, invitations[i].GitHubID(), invitation.GitHubID())
			assert.Equal(t, invitations[i].TraqID(), invitation.TraqID())
		}
	})
}
