package database

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/Kittipoom-pan/autopart-service/config"
	_ "github.com/go-sql-driver/mysql"
)

var (
	once       sync.Once
	dbInstance *sql.DB
	dbErr      error
)

func NewMySQLDatabase(conf *config.Config) (*sql.DB, error) {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
			conf.Db.MySqlUser,
			conf.Db.MySqlPassword,
			conf.Db.MySqlHost,
			conf.Db.MySqlPort,
			conf.Db.MySqlDatabase,
		)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			dbErr = fmt.Errorf("can't open database connection: %w", err)
			return
		}

		if err := db.Ping(); err != nil {
			db.Close()
			dbErr = fmt.Errorf("can't connect to database: %w", err)
			return
		}

		fmt.Println("MySQL has been connected.")
		dbInstance = db
	})

	return dbInstance, dbErr
}
