package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	// "fmt"

	"github.com/KznRkjp/go-link-shortner.git/internal/flags"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Ping(res http.ResponseWriter, req *http.Request) {
	// flags.ParseFlags()
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = conn.PingContext(ctx); err != nil {
		http.Error(res, "Can't connect to DB", http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func CreateTable() {
	fmt.Println("DB String", flags.FlagDBString)
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		fmt.Println("DB error")
		panic(err)
	}
	defer conn.Close()

	insertDynStmt := "CREATE TABLE url (id SERIAL PRIMARY KEY, correlationid TEXT,shorturl TEXT, originalurl TEXT);"
	_, err = conn.Exec(insertDynStmt)
	if err != nil {
		fmt.Println("Database exists")
	}

}

func WriteToDB(url string, originalURL string, correlationID string) {
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	insertDynStmt := `insert into "url"("shorturl", "originalurl", "correlationid") values($1, $2, $3)`
	if correlationID == "nil" {
		_, err = conn.Exec(insertDynStmt, url, originalURL, nil)
	} else {
		_, err = conn.Exec(insertDynStmt, url, originalURL, correlationID)
	}
	if err != nil {
		panic(err)
	}
}

func GetFromDB(shortURL string) (string, error) {
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		return "", err
	}

	defer conn.Close()
	insertDynStmt := `SELECT originalurl, correlationid, shorturl FROM url where shorturl =` + shortURL

	row := conn.QueryRowContext(context.Background(),
		insertDynStmt)

	// if err != nil {
	// 	return "",err
	// }
	var originalURL string
	err = row.Scan(&originalURL)
	if err != nil {
		panic(err)
	}
	fmt.Println(originalURL)
	return originalURL, err
}
