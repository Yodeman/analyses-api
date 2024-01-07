package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/yodeman/analyses-api/api"
	db "github.com/yodeman/analyses-api/dbase/sqlc"
	"github.com/yodeman/analyses-api/util"
)

func main() {
	dbPass, found := os.LookupEnv("DBASE_PASSWORD")
	if !found {
		log.Fatal("Cannot find database password!!!")
	}
	tokenSymmKey, found := os.LookupEnv("TOKEN_SYMMETRIC_KEY")
	if !found {
		log.Fatal("Cannot find token symmetric key!!!")
	}
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalln(err)
	}
	config.TokenSymmetricKey = tokenSymmKey

	connStr := fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=disable",
		config.DBDriver, config.DBUser, dbPass, config.DBAddr,
		config.DBName)
	conn, err := sql.Open(config.DBDriver, connStr)
	if err != nil {
		log.Fatalln(err)
	}

	server, err := api.NewServer(config, db.New(conn))
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalf("%v\n", server.Start(config.ServerAddr))
}
