package model

type Invitation struct {
	messageID string
	traqID    string
	gitHubID  string
}

func NewInvitation(id string, traqID, gitHubID string) *Invitation {
	return &Invitation{
		messageID: id,
		traqID:    traqID,
		gitHubID:  gitHubID,
	}
}

func (i *Invitation) MessageID() string {
	return i.messageID
}

func (i *Invitation) TraqID() string {
	return i.traqID
}

func (i *Invitation) GitHubID() string {
	return i.gitHubID
}
