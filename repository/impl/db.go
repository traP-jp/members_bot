package impl

import (
	"cmp"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

func NewDB() (*bun.DB, error) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, fmt.Errorf("failed to load location: %w", err)
	}

	conf := mysql.Config{
		User:                 cmp.Or(os.Getenv("NS_MARIADB_USER"), os.Getenv("MYSQL_USER"), "root"),
		Passwd:               cmp.Or(os.Getenv("NS_MARIADB_PASSWORD"), os.Getenv("MYSQL_PASSWORD"), "pass"),
		Net:                  "tcp",
		Addr:                 cmp.Or(os.Getenv("NS_MARIADB_HOST"), os.Getenv("MYSQL_HOST"), "db") + ":" + cmp.Or(os.Getenv("NS_MARIADB_PORT"), os.Getenv("MYSQL_PORT"), "3306"),
		DBName:               cmp.Or(os.Getenv("NS_MARIADB_DATABASE"), os.Getenv("MYSQL_DATABASE"), "members_bot"),
		Loc:                  jst,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	sqldb, err := sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, mysqldialect.New())

	return db, nil
}
