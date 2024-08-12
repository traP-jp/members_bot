package model

type Invitation struct {
	id       string
	traqID   string
	gitHubID string
}

func NewInvitation(id string, traqID, gitHubID string) *Invitation {
	return &Invitation{
		id:       id,
		traqID:   traqID,
		gitHubID: gitHubID,
	}
}

func (i *Invitation) ID() string {
	return i.id
}

func (i *Invitation) TraqID() string {
	return i.traqID
}

func (i *Invitation) GitHubID() string {
	return i.gitHubID
}
