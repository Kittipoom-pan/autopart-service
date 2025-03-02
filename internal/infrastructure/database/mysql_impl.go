package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/Kittipoom-pan/autopart-service/config"
	_ "github.com/go-sql-driver/mysql"
)

type mysqlDatabase struct {
	Db *sql.DB
}

var (
	once       sync.Once
	dbInstance *mysqlDatabase
)

func NewMySQLDatabase(conf *config.Config) Database {
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
			log.Fatalf("failed to connect database: %v", err)
		}

		fmt.Println("Connected MySQL successfully.")
		dbInstance = &mysqlDatabase{Db: db}
	})

	return dbInstance
}

func (p *mysqlDatabase) GetDb() *sql.DB {
	return dbInstance.Db
}
