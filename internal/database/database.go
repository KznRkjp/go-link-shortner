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
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	_, tableCheck := conn.Query("select * from  url;")

	if tableCheck == nil {
		fmt.Println("table is there")
	} else {
		fmt.Println("table not there")
		// dynamic
		insertDynStmt := "CREATE TABLE url (id SERIAL PRIMARY KEY, shorturl TEXT, originalurl TEXT);"
		_, err = conn.Exec(insertDynStmt)
		if err != nil {
			panic(err)
		}
	}

}

func WriteToDB(url string, originalUrl string) {
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	insertDynStmt := `insert into "url"("shorturl", "originalurl") values($1, $2)`
	_, err = conn.Exec(insertDynStmt, url, originalUrl)
	if err != nil {
		panic(err)
	}
}
