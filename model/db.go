package model

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDb(dbname string) {
	var err error

	db, err = sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatalln("Initialize DataBase ", err.Error())
	}

	CreateTables()
}

func CreateTables() {
	indirect_cfg_tbl := `CREATE TABLE IF NOT EXISTS "indirect_cfg"(
						"id" INTEGER PRIMARY KEY AUTOINCREMENT,
						"protocol" CHAR(3) NULL,
						"listen-addr" VARCHAR(64) NULL,
						"listen-port" VARCHAR(16) NULL,
						"dest-addr" VARCHAR(64) NULL,
						"dest-port" VARCHAR(16) NULL,
						"acl" VARCHAR(16) NULL,
						"deny-addr" VARCHAR(1024) NULL,
						"admit-addr" VARCHAR(1024) NULL,
						"max-conns" VARCHAR(16) NULL,
						"memo" TEXT NULL
					)`

	log_tbl := `CREATE TABLE IF NOT EXISTS "log"(
					"id" INTEGER PRIMARY KEY AUTOINCREMENT,
					"time" TIMESTAMP NOT NULL DEFAULT (datetime('now','localtime')),
					"content" TEXT
				)`

	if _, err := db.Exec(indirect_cfg_tbl); err != nil {
		log.Fatalln(err.Error())
	}

	if _, err := db.Exec(log_tbl); err != nil {
		log.Fatalln(err.Error())
	}
}
