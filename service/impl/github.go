package impl

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v63/github"
	"github.com/traP-jp/members_bot/model"
	"github.com/traP-jp/members_bot/service"
)

var _ service.GitHub = &GitHub{}

type GitHub struct {
	cl      *github.Client
	orgName string
}

func NewGitHub(orgName string) (*GitHub, error) {
	gitHubAppIDStr := os.Getenv("GITHUB_APP_ID")
	gitHubAppID, err := strconv.ParseInt(gitHubAppIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GITHUB_APP_ID: %w", err)
	}

	gitHubAppInstallIDStr := os.Getenv("GITHUB_APP_INSTALLATION_ID")
	gitHubAppInstallationID, err := strconv.ParseInt(gitHubAppInstallIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GITHUB_APP_INSTALL_ID: %w", err)
	}

	privateKeyStr, ok := os.LookupEnv("GITHUB_APP_PRIVATE_KEY")
	if !ok {
		return nil, errors.New("GITHUB_PRIVATE_KEY is not set")
	}

	log.Println(privateKeyStr)

	log.Println(`aa
bb`)

	irt, err := ghinstallation.New(http.DefaultTransport, gitHubAppID, gitHubAppInstallationID, []byte(privateKeyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub installation round tripper: %w", err)
	}

	// token := os.Getenv("GITHUB_TOKEN")
	cl := github.NewClient(&http.Client{Transport: irt})

	return &GitHub{
		cl:      cl,
		orgName: orgName,
	}, nil
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

func (g *GitHub) CheckUserExist(ctx context.Context, userID string) (bool, error) {
	user, _, err := g.cl.Users.Get(ctx, userID)
	var gitHubErr *github.ErrorResponse
	if errors.As(err, &gitHubErr) && gitHubErr.Response.StatusCode == 404 {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get GitHub user: %w", err)
	}

	if user.GetType() != "User" {
		return false, nil
	}

	return true, nil
}

func (g *GitHub) OrgName() string {
	return g.orgName
}
