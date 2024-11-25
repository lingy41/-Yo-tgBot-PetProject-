package environment

import (
	"flag"
	"log"
)

func MustToken() string {
	token := flag.String(
		"token-bot",
		"",
		"token for acces to telegram bot",
	)

	flag.Parse()
	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}

func MustPath() string {
	return "postgres://postgres:psql123@localhost:54321/postgres"
	// pathToDb := flag.String(
	// 	"database-path",
	// 	"",
	// 	"path with user ind password to connect to database",
	// )

	// flag.Parse()

	// if *pathToDb == "" {
	// 	log.Fatal("path to databse is not specified")
	// }
	// return *pathToDb
}

func NewBasePath(token string) string {
	return "bot" + token
}
