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

// линк к БД, создается в начале и используется при всех запросах к БД
var DB *sql.DB

// Ping - функция проверки статуса подключения к BD
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

// CreateTable - функция первоначального создания таблиц, выполняется каждый раз при
// запуске приложения, если таблицы уже есть - информирует ошибкой
func CreateTable(db *sql.DB) {
	fmt.Println("DB String", flags.FlagDBString)
	ctx := context.Background()
	insertDynStmtUser := `CREATE TABLE url_users (id SERIAL PRIMARY KEY, uuid TEXT UNIQUE, token TEXT);`
	var err error
	_, err = db.ExecContext(ctx, insertDynStmtUser)
	if err != nil {
		log.Println("Database 'user' exists", err)
	}

	insertDynStmtURL := `CREATE TABLE url (id SERIAL PRIMARY KEY,
		 									correlationid TEXT,
											url_user_uuid TEXT,
											shorturl TEXT, 
											originalurl TEXT,
											deleted_flag BOOLEAN DEFAULT FALSE,
											CONSTRAINT fk_url_user_uuid FOREIGN KEY (url_user_uuid) REFERENCES url_users (uuid));`
	_, err = db.ExecContext(ctx, insertDynStmtURL)
	if err != nil {
		log.Println("Database 'url' exists", err)
	}

}

// WriteToDB - функция записи одной строки в БД
func WriteToDB(db *sql.DB, ctx context.Context, url string, originalURL string, correlationID string, uuid string) {

	insertDynStmt := `insert into "url"("shorturl", "originalurl", "correlationid", "url_user_uuid") values($1, $2, $3, $4)`
	var err error
	if correlationID == "nil" {

		_, err = db.ExecContext(ctx, insertDynStmt, url, originalURL, nil, uuid)
	} else {

		_, err = db.ExecContext(ctx, insertDynStmt, url, originalURL, correlationID, uuid)
	}
	if err != nil {
		log.Println(err)

	}
	// CreateUser(ctx)
}

// WriteToDBBatch - функция массовой записи в БД
func WriteToDBBatch(db *sql.DB, ctx context.Context, listURL []models.BatchRequest, uuid string) error {

	tx, err := db.Begin()
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

// GetFromDB - функция получения данных их БД
func GetFromDB(db *sql.DB, ctx context.Context, shortURL string) (string, bool, error) {

	insertDynStmt := `SELECT originalurl, deleted_flag FROM url where shorturl = '` + shortURL + `'`
	var err error
	row := db.QueryRowContext(ctx,
		insertDynStmt)
	var originalurl string
	var deletedFlag bool

	err = row.Scan(&originalurl, &deletedFlag)

	if err != nil {
		log.Print(err)
	}

	return originalurl, deletedFlag, err
}

// CheckForDuplicates - преверка есть ли в БД уже такие данные
func CheckForDuplicates(db *sql.DB, ctx context.Context, URL string, URLDb map[string]filesio.URLRecord, uuid string) (string, error) {
	if flags.FlagDBString != "" {

		var err error
		insertDynStmt := `SELECT shorturl FROM url where originalurl = '` + URL + `'`

		row := db.QueryRowContext(ctx,
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

// GetUserFromDB - получение данных пользователя
func GetUserFromDB(db *sql.DB, ctx context.Context, uuid string) (int, error) {
	var err error
	insertDynStmt := `SELECT id FROM url_users where uuid = '` + uuid + `'`
	row := db.QueryRowContext(ctx,
		insertDynStmt)
	var id int

	err = row.Scan(&id)

	if err != nil {
		log.Println(err)
	}

	return id, err

}

// UpdateUserToken - обновлдение токена пользователя в БД, не знаю зачем я его вообще храню, надо переделать логику.
func UpdateUserToken(db *sql.DB, ctx context.Context, uuid string, token string) error {
	var err error
	insertDynStmt := `UPDATE url_users SET token = $2 WHERE uuid = $1`
	_, err = db.ExecContext(ctx, insertDynStmt, uuid, token)
	if err != nil {
		fmt.Println("Error updating token: ", err)
		return err
	}
	return err

}

// CreateUser - функция записи данных о пользователе в БД
func CreateUser(db *sql.DB, ctx context.Context) (string, string, error) {
	// log.Println("Creating user - database.CreateUser")

	uuid := shortuuid.New()
	insertDynStmt := `insert into "url_users"("uuid", "token") values($1, $2)`
	token, err := users.BuildJWTString(uuid)

	if err != nil {
		log.Println(err)
		return "", "", err

	}

	_, err = db.ExecContext(ctx, insertDynStmt, uuid, token)
	if err != nil {
		log.Println(err)
		return "", "", err

	}
	return uuid, token, nil

}

// GetOrCreateUser - функция проверки существования пользователя, если его нет - то вызывается запись
func GetOrCreateUser(ctx context.Context, uuid string) (string, string, error) {
	_, err := GetUserFromDB(DB, ctx, uuid)
	if err != nil {
		log.Println(err)
		log.Println("Creating new user")
		newUUID, token, err := CreateUser(DB, ctx)
		if err != nil {
			return newUUID, token, err
		} else {
			return newUUID, token, err
		}

	} else {
		return uuid, "", err
	}

}

// GetUsersUrls - функция получения списка URL пользователя из БД
func GetUsersUrls(db *sql.DB, ctx context.Context, uuid string) ([]models.URLResponse, error) {

	insertDynStmt := `SELECT shorturl, originalurl FROM url WHERE url_user_uuid = $1`
	rows, err := db.QueryContext(ctx, insertDynStmt, uuid)
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

// DeleteUsersUrls - функция помечания URLов пользователя как уделенные
func DeleteUsersUrls(db *sql.DB, ctx context.Context, uuid string, ch chan []string) error {
	var err error
	var insertDynStmt = `
	UPDATE url
	SET deleted_flag = true
	WHERE shorturl = $1`
	// WHERE url_user_uuid = $1 and shorturl = $2`
	for urlList := range ch {
		for i := range urlList {
			_, err = db.Exec(insertDynStmt, urlList[i])
			// _, err = conn.Exec(insertDynStmt, uuid, urlList[i])
			if err != nil {
				log.Println(err)
				return err

			}
		}

	}
	return nil
}

// GetStats получает данные от ДБ по количеству пользователей и ссылок
func GetStats(db *sql.DB, ctx context.Context) (models.Stats, error) {
	var result models.Stats
	countURLStmt := `SELECT COUNT(*) from url`
	countUsersStmt := `SELECT COUNT(*) from url_users`
	url, err := db.QueryContext(ctx, countURLStmt)
	if err != nil {
		log.Println(err)
		return result, err

	}
	defer url.Close()

	users, err := db.QueryContext(ctx, countUsersStmt)
	if err != nil {
		log.Println(err)
		return result, err

	}
	defer users.Close()
	if url.Err() != nil {
		log.Println(url.Err())
	}
	for url.Next() {
		err = url.Scan(&result.URLs)
		if err != nil {
			log.Println(err)
			return result, err

		}
	}

	if users.Err() != nil {
		log.Println(users.Err())
	}
	for users.Next() {
		err = users.Scan(&result.Users)
		if err != nil {
			log.Println(err)
			return result, err

		}
	}
	return result, err
}
