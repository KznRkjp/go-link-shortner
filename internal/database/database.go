package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/KznRkjp/go-link-shortner.git/internal/filesio"
	"github.com/KznRkjp/go-link-shortner.git/internal/flags"
	"github.com/KznRkjp/go-link-shortner.git/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Ping(res http.ResponseWriter, req *http.Request) {
	// flags.ParseFlags()
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		log.Println(err)
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
		// fmt.Println("DB error")
		log.Panic(err)
	}
	defer conn.Close()

	insertDynStmt := "CREATE TABLE url (id SERIAL PRIMARY KEY, correlationid TEXT,shorturl TEXT, originalurl TEXT);"
	_, err = conn.Exec(insertDynStmt)
	if err != nil {
		log.Println("Database exists", err)
	}

}

func WriteToDB(url string, originalURL string, correlationID string) {
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	insertDynStmt := `insert into "url"("shorturl", "originalurl", "correlationid") values($1, $2, $3)`
	if correlationID == "nil" {
		_, err = conn.Exec(insertDynStmt, url, originalURL, nil)
	} else {
		_, err = conn.Exec(insertDynStmt, url, originalURL, correlationID)
	}
	if err != nil {
		log.Println(err)
	}
}

func WriteToDBBatch(listURL []models.BatchRequest) error {
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		return err
	}
	defer conn.Close()
	ctx := context.Background()
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	for _, v := range listURL {
		// все изменения записываются в транзакцию
		// shortURL := urlgen.GenerateShortKey()
		_, err := tx.ExecContext(ctx,
			"INSERT INTO url (shorturl, originalurl, correlationid)"+
				" VALUES($1,$2,$3)", v.ShortURL, v.URL, v.CorrelationID)
		if err != nil {
			fmt.Println("error in here")
			// если ошибка, то откатываем изменения
			tx.Rollback()
			return err
		}
	}
	// завершаем транзакцию
	return tx.Commit()

}

func GetFromDB(shortURL string) (string, error) {
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		return "", err
	}

	defer conn.Close()
	insertDynStmt := `SELECT originalurl FROM url where shorturl = '` + shortURL + `'`

	row := conn.QueryRowContext(context.Background(),
		insertDynStmt)

	// if err != nil {
	// 	return "",err
	// }
	var originalurl string

	err = row.Scan(&originalurl)

	if err != nil {
		panic(err)
	}

	return originalurl, err
}

func CheckForDuplicates(ctx context.Context, URL string, URLDb map[string]filesio.URLRecord) (string, error) {
	if flags.FlagDBString != "" {

		conn, err := sql.Open("pgx", flags.FlagDBString)
		if err != nil {
			return "", err
		}

		defer conn.Close()
		insertDynStmt := `SELECT shorturl FROM url where originalurl = '` + URL + `'`

		row := conn.QueryRowContext(ctx,
			insertDynStmt)
		fmt.Println("Check for duplicates")

		var shorturl string

		err = row.Scan(&shorturl)

		if err != nil {
			return "", err
		}
		return shorturl, err

	} else if len(flags.FlagDBFilePath) > 1 {
		for _, value := range URLDb {
			if value.OriginalURL == URL {
				return value.ShortURL, nil
			}
		}
		return "", fmt.Errorf("duplicate url (id: %s)", URL)
	}
	return "", nil
}
