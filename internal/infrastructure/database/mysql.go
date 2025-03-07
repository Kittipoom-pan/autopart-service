package database

import (
	"database/sql"

	db "github.com/Kittipoom-pan/autopart-service/internal/infrastructure/database/sqlc"

	"fmt"
	"sync"

	"github.com/Kittipoom-pan/autopart-service/config"
	_ "github.com/go-sql-driver/mysql"
)

var (
	once       sync.Once
	dbInstance *db.Queries
	dbErr      error
)

func NewMySQLDatabase(conf *config.Config) (*db.Queries, error) {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
			conf.Db.MySqlUser,
			conf.Db.MySqlPassword,
			conf.Db.MySqlHost,
			conf.Db.MySqlPort,
			conf.Db.MySqlDatabase,
		)

		sqlDB, err := sql.Open("mysql", dsn)
		if err != nil {
			dbErr = fmt.Errorf("can't open database connection: %w", err)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			sqlDB.Close()
			dbErr = fmt.Errorf("can't connect to database: %w", err)
			return
		}

		fmt.Println("MySQL has been connected.")
		dbInstance = db.New(sqlDB)
	})

	return dbInstance, dbErr
}
