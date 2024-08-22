package impl

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/traP-jp/members_bot/repository/impl/schema"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

var testDB *bun.DB

const (
	testDBRootPassword = "password"
	testDBDatabase     = "members_bot"
	testDBPort         = "3306"
	testDBHostName     = "testDB"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "mariadb:lts",
			Env: map[string]string{
				"MYSQL_DATABASE":      testDBDatabase,
				"MYSQL_ROOT_PASSWORD": testDBRootPassword,
			},
			ExposedPorts: []string{testDBPort + "/tcp"},
			WaitingFor:   wait.ForListeningPort(testDBPort + "/tcp"),
		},
		Started: true,
	})
	if err != nil {
		log.Fatalf("failed to create mysql container: %s", err)
	}
	defer func() {
		if err := mysqlC.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate mysql container: %s", err)
		}
	}()

	host, err := mysqlC.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get host: %s", err)
	}
	port, err := mysqlC.MappedPort(ctx, "3306/tcp")
	if err != nil {
		log.Fatalf("failed to get externally mapped port: %s", err)
	}

	db, err := setUpTestDB(host, port.Port())
	if err != nil {
		log.Println("failed to set up test db:", err)
		os.Exit(1)
	}

	testDB = db

	if err := schema.Migrate(testDB); err != nil {
		log.Println("failed to migrate schema:", err)
		os.Exit(1)
	}

	m.Run()
}

func setUpTestDB(host string, port string) (*bun.DB, error) {
	conf := mysql.Config{
		User:                 "root",
		Passwd:               testDBRootPassword,
		Net:                  "tcp",
		Addr:                 host + ":" + port,
		DBName:               testDBDatabase,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	sqldb, err := sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open sql: %w", err)
	}

	db := bun.NewDB(sqldb, mysqldialect.New())

	return db, nil
}
