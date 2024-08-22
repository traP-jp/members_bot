package migrate

import (
	"context"
	"fmt"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

var m = []func(*migrate.Migrations){
	v1,
	v2,
}

func Migrate(db *bun.DB) error {
	migrations := migrate.NewMigrations()

	for _, f := range m {
		f(migrations)
	}

	ctx := context.Background()

	migrator := migrate.NewMigrator(db, migrations)
	err := migrator.Init(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize migrator: %w", err)
	}

	g, err := migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	if g.IsZero() {
		log.Println("database is up-to-date")
	}

	return nil
}
