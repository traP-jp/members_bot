package migrate

import (
	"context"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

var m = []func(*migrate.Migrations){
	v1,
}

func Migrate(db *bun.DB) {
	migrations := migrate.NewMigrations()

	for _, f := range m {
		f(migrations)
	}

	ctx := context.Background()

	migrator := migrate.NewMigrator(db, migrations)
	err := migrator.Init(ctx)
	if err != nil {
		panic(err)
	}

	g, err := migrator.Migrate(context.Background())
	if err != nil {
		panic(err)
	}

	if g.IsZero() {
		log.Println("database is up-to-date")
	}
}
