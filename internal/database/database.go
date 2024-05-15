package database

import (
	"context"
	"database/sql"
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
