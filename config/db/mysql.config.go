package config

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/srikanta0427/rest_api_design/config"
)

func SetUpDB() (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = config.GetString("DB_USER", "root")
	cfg.Passwd = config.GetString("DB_PASS", "root")
	cfg.Net = "tcp"
	cfg.Addr = config.GetString("DB_HOST", "127.0.0.1")
	cfg.DBName = config.GetString("DB_NAME", "user")

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		fmt.Println("Error connecting to DB")
		panic(err.Error())
		return nil, err
	}

	errPing := db.Ping()
	if errPing != nil {
		fmt.Println("Error pinging DB")
		panic(errPing.Error())
		return nil, errPing
	}
	fmt.Println("Successfully connected to DB")
	return db, nil
}
