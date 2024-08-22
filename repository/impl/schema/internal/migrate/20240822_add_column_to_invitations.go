package migrate

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

type InvitationV2 struct {
	bun.BaseModel `bun:"table:invitations"`
	ID            int `bun:",pk,autoincrement"`
	MessageID     string
	TraqID        string
	GitHubID      string
	CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}

func v2(m *migrate.Migrations) {
	m.MustRegister(
		func(ctx context.Context, db *bun.DB) (err error) {
			return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {

				_, err = tx.NewRaw("ALTER TABLE invitations RENAME COLUMN id TO message_id").Exec(ctx)
				if err != nil {
					return fmt.Errorf("failed to rename column: %w", err)
				}

				_, err = tx.NewRaw("ALTER TABLE invitations DROP PRIMARY KEY").Exec(ctx)
				if err != nil {
					return fmt.Errorf("failed to drop primary key: %w", err)
				}

				_, err = tx.NewRaw("ALTER TABLE invitations ADD COLUMN id INT AUTO_INCREMENT PRIMARY KEY").Exec(ctx)
				if err != nil {
					return fmt.Errorf("failed to add column: %w", err)
				}

				return nil
			})
		},
		func(ctx context.Context, db *bun.DB) (err error) {
			return db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
				_, err = tx.NewRaw("ALTER TABLE invitations DROP COLUMN id").Exec(ctx)
				if err != nil {
					return fmt.Errorf("failed to drop column: %w", err)
				}

				_, err = tx.NewRaw("ALTER TABLE invitations RENAME COLUMN message_id TO id").Exec(ctx)
				if err != nil {
					return fmt.Errorf("failed to rename column: %w", err)
				}

				_, err = tx.NewRaw("ALTER TABLE invitations ADD PRIMARY KEY (id)").Exec(ctx)
				if err != nil {
					return fmt.Errorf("failed to add primary key: %w", err)
				}

				return nil
			})
		},
	)
}
