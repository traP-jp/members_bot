package impl

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traP-jp/members_bot/model"
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
			assert.Equal(t, invitations[i].ID(), invitation.ID)
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
	}

	invitationID1 := uuid.NewString()
	invitationID2 := uuid.NewString()

	testCases := map[string]testCase{
		"特に問題なし": {
			invitationID: invitationID1,
			fixture: []*schema.Invitation{
				{ID: invitationID1, GitHubID: "github_id", TraqID: "traq_id"},
			},
			expected: []*model.Invitation{
				model.NewInvitation(invitationID1, "traq_id", "github_id"),
			},
		},
		"他の招待があっても問題なし": {
			invitationID: invitationID2,
			fixture: []*schema.Invitation{
				{ID: invitationID2, GitHubID: "github_id2", TraqID: "traq_id2"},
				{ID: uuid.NewString(), GitHubID: "github_id3", TraqID: "traq_id3"},
			},
			expected: []*model.Invitation{
				model.NewInvitation(invitationID2, "traq_id2", "github_id2"),
			},
		},
		"招待がない": {
			invitationID: uuid.NewString(),
			fixture:      []*schema.Invitation{},
			expected:     []*model.Invitation{},
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

			assert.NoError(t, err)
			assert.Len(t, invitations, len(test.expected))
			for i, invitation := range invitations {
				assert.Equal(t, test.expected[i].ID(), invitation.ID())
				assert.Equal(t, test.expected[i].GitHubID(), invitation.GitHubID())
				assert.Equal(t, test.expected[i].TraqID(), invitation.TraqID())
			}
		})
	}
}
