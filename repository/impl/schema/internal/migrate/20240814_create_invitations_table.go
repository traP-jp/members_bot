package migrate

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

type InvitationV1 struct {
	bun.BaseModel `bun:"table:invitations"`
	ID            string `bun:",pk"`
	TraqID        string
	GitHubID      string
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func v1(m *migrate.Migrations) {
	m.MustRegister(
		func(ctx context.Context, db *bun.DB) error {
			_, err := db.NewCreateTable().
				Model(&InvitationV1{}).
				Exec(ctx)
			return err
		},
		func(ctx context.Context, db *bun.DB) error {
			_, err := db.NewDropTable().
				Model(&InvitationV1{}).
				IfExists().
				Exec(ctx)
			return err
		},
	)
}
