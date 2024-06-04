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
	"github.com/KznRkjp/go-link-shortner.git/internal/users"

	// "github.com/KznRkjp/go-link-shortner.git/internal/users"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lithammer/shortuuid"
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
	ctx := context.Background()
	insertDynStmtUser := `CREATE TABLE url_users (id SERIAL PRIMARY KEY, uuid TEXT UNIQUE, token TEXT);`
	_, err = conn.ExecContext(ctx, insertDynStmtUser)
	if err != nil {
		log.Println("Database user exists", err)
	}

	insertDynStmtURL := `CREATE TABLE url (id SERIAL PRIMARY KEY,
		 									correlationid TEXT,
											url_user_uuid TEXT,
											shorturl TEXT, 
											originalurl TEXT,
											deleted_flag BOOLEAN DEFAULT FALSE,
											CONSTRAINT fk_url_user_uuid FOREIGN KEY (url_user_uuid) REFERENCES url_users (uuid));`
	_, err = conn.ExecContext(ctx, insertDynStmtURL)
	if err != nil {
		log.Println("Database url exists", err)
	}

}

func WriteToDB(ctx context.Context, url string, originalURL string, correlationID string, uuid string) {
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	insertDynStmt := `insert into "url"("shorturl", "originalurl", "correlationid", "url_user_uuid") values($1, $2, $3, $4)`
	if correlationID == "nil" {
		_, err = conn.ExecContext(ctx, insertDynStmt, url, originalURL, nil, uuid)
	} else {
		_, err = conn.ExecContext(ctx, insertDynStmt, url, originalURL, correlationID, uuid)
	}
	if err != nil {
		log.Println(err)

	}
	CreateUser(ctx)
}

func WriteToDBBatch(ctx context.Context, listURL []models.BatchRequest, uuid string) error {
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		return err
	}
	defer conn.Close()
	// ctx := context.Background()
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	for _, v := range listURL {
		// все изменения записываются в транзакцию
		// shortURL := urlgen.GenerateShortKey()
		_, err := tx.ExecContext(ctx,
			"INSERT INTO url (shorturl, originalurl, correlationid, url_user_uuid)"+
				" VALUES($1,$2,$3,$4)", v.ShortURL, v.URL, v.CorrelationID, uuid)
		if err != nil {
			log.Println("error in WriteToDBBatch - writing to DB", err)

			// если ошибка, то откатываем изменения
			tx.Rollback()
			return err
		}
	}
	// завершаем транзакцию
	return tx.Commit()

}

func GetFromDB(ctx context.Context, shortURL string) (string, bool, error) {
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		return "", false, err
	}

	defer conn.Close()
	insertDynStmt := `SELECT originalurl, deleted_flag FROM url where shorturl = '` + shortURL + `'`

	row := conn.QueryRowContext(ctx,
		insertDynStmt)

	// if err != nil {
	// 	return "",err
	// }
	var originalurl string
	var deletedFlag bool

	err = row.Scan(&originalurl, &deletedFlag)

	if err != nil {
		panic(err)
	}

	return originalurl, deletedFlag, err
}

func CheckForDuplicates(ctx context.Context, URL string, URLDb map[string]filesio.URLRecord, uuid string) (string, error) {
	if flags.FlagDBString != "" {
		conn, err := sql.Open("pgx", flags.FlagDBString)
		if err != nil {
			return "", err
		}
		defer conn.Close()
		if uuid != "" {
			// insertDynStmt := `SELECT shorturl FROM url where originalurl = $1 and url_user_uuid = $2`
			insertDynStmt := `SELECT shorturl FROM url where originalurl = $1`
			row := conn.QueryRowContext(ctx,
				insertDynStmt, URL)
			fmt.Println("Checking for duplicates")

			var shorturl string

			err = row.Scan(&shorturl)

			if err != nil {
				log.Println("Duplicates not found")
				return "", err
			}

			fmt.Println("Duplicates found")
			return shorturl, err
		}
		insertDynStmt := `SELECT shorturl FROM url where originalurl = '` + URL + `'`

		row := conn.QueryRowContext(ctx,
			insertDynStmt)

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

func GetUserFromDB(ctx context.Context, uuid string) (int, error) {
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		return -1, err
	}
	defer conn.Close()
	insertDynStmt := `SELECT id FROM url_users where uuid = '` + uuid + `'`
	row := conn.QueryRowContext(ctx,
		insertDynStmt)
	var id int

	err = row.Scan(&id)

	if err != nil {
		panic(err)
	}

	return id, err

}

func UpdateUserToken(ctx context.Context, uuid string, token string) error {
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	insertDynStmt := `UPDATE url_users SET token = $2 WHERE uuid = $1`
	_, err = conn.ExecContext(ctx, insertDynStmt, uuid, token)
	if err != nil {
		fmt.Println("Error updating token: ", err)
		return err
	}
	return err

}

func CreateUser(ctx context.Context) (string, string, error) {
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	uuid := shortuuid.New()
	insertDynStmt := `insert into "url_users"("uuid", "token") values($1, $2)`
	token, err := users.BuildJWTString(uuid)

	if err != nil {
		log.Println(err)
		return "", "", err

	}

	_, err = conn.ExecContext(ctx, insertDynStmt, uuid, token)
	if err != nil {
		log.Println(err)
		return "", "", err

	}
	return uuid, token, nil

}

func GetOrCreateUser(ctx context.Context, uuid string) (string, string, error) {
	_, err := GetUserFromDB(ctx, uuid)
	if err != nil {
		log.Println(err)
		log.Println("Creating new user")
		newUUID, token, err := CreateUser(ctx)
		if err != nil {
			return newUUID, token, err
		} else {
			return newUUID, token, err
		}

	} else {
		return uuid, "", err
	}

}

func GetUsersUrls(ctx context.Context, uuid string) ([]models.URLResponse, error) {
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	insertDynStmt := `SELECT shorturl, originalurl FROM url WHERE url_user_uuid = $1`
	rows, err := conn.QueryContext(ctx, insertDynStmt, uuid)
	if err != nil {
		log.Println(err)
	}
	if rows.Err() != nil {
		log.Println(rows.Err())
	}

	var urls []models.URLResponse
	for rows.Next() {
		var url models.URLResponse
		if err := rows.Scan(&url.ShortURL, &url.OriginalURL); err != nil {
			return urls, err
		}
		urls = append(urls, url)

	}
	// fmt.Println(urls)
	return urls, err

}

func DeleteUsersUrls(ctx context.Context, uuid string, ch chan []string) error {
	conn, err := sql.Open("pgx", flags.FlagDBString)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	for urlList := range ch {
		var insertDynStmt = `
			UPDATE url
			SET deleted_flag = true
			WHERE url_user_uuid = $1 and shorturl = $2`
		tx, err := conn.Begin()
		for i := range urlList {
			_, err = tx.ExecContext(ctx, insertDynStmt, uuid, urlList[i])
			if err != nil {
				log.Println(err)
				tx.Rollback()
				return err
			}

		}
		return tx.Commit()
	}
	return nil

}
