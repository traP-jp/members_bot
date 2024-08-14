package impl

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v63/github"
	"github.com/traP-jp/members_bot/model"
	"github.com/traP-jp/members_bot/service"
)

var _ service.GitHub = &GitHub{}

type GitHub struct {
	cl      *github.Client
	orgName string
}

func NewGitHub(orgName string) *GitHub {
	token := os.Getenv("GITHUB_TOKEN")
	cl := github.NewClient(nil).WithAuthToken(token)
	return &GitHub{
		cl:      cl,
		orgName: orgName,
	}
}

func (g *GitHub) SendInvitations(ctx context.Context, invitations []*model.Invitation) error {
	gitHubUserIDs := make([]int64, 0, len(invitations))

	for _, invitation := range invitations {
		user, _, err := g.cl.Users.Get(ctx, invitation.GitHubID())
		if err != nil {
			return fmt.Errorf("failed to get GitHub user: %w", err)
		}

		gitHubUserIDs = append(gitHubUserIDs, user.GetID())
	}

	for _, gitHubUserID := range gitHubUserIDs {
		_, _, err := g.cl.Organizations.CreateOrgInvitation(ctx, g.orgName, &github.CreateOrgInvitationOptions{
			InviteeID: &gitHubUserID,
		})
		if err != nil {
			return fmt.Errorf("failed to create GitHub invitation: %w", err)
		}
	}

	return nil
}
