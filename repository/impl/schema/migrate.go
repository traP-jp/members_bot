package schema

import (
	"github.com/traP-jp/members_bot/repository/impl/schema/internal/migrate"
	"github.com/uptrace/bun"
)

func Migrate(db *bun.DB) error {
	return migrate.Migrate(db)
}
