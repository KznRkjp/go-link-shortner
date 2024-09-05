package grpc

import (
	context "context"
	"fmt"
	"log"

	"github.com/KznRkjp/go-link-shortner.git/internal/app"
	"github.com/KznRkjp/go-link-shortner.git/internal/database"
	"github.com/KznRkjp/go-link-shortner.git/internal/flags"
	"github.com/KznRkjp/go-link-shortner.git/internal/users"
	"google.golang.org/grpc/metadata"
)

// GrpcHandlers - структура для gRPC сервера
type GrpcHandlers struct {
	UnimplementedHandlersServer

	// service *service.Service
}

// GetUserFromContext возвращает значение токена пользователя из контекста
func GetUserFromContext(ctx context.Context) string {
	var user string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("User")
		if len(values) > 0 {
			user = values[0]
		}
	}
	return user
}

func (s *GrpcHandlers) ShortenURL(ctx context.Context, in *ShortenURLRequest) (*ShortenURLResponse, error) {
	var response ShortenURLResponse
	url := in.LongURL
	fmt.Println(url)
	uuid, token := ManageCookieGRPC(ctx)
	// fmt.Println(uuid, token)
	shortURL, err := database.CheckForDuplicates(database.DB, context.Background(), url, app.URLDb, uuid)
	var resultURL string
	if err != nil {
		// log.Print(err)
		body := []byte(url)
		fmt.Println(uuid)
		resultURL = app.SaveData(context.Background(), body, uuid)
	} else {

		resultURL = shortURL
		log.Println(err)
	}
	response.Token = token
	response.ShortURL = resultURL
	// fmt.Println(resultURL)
	return &response, nil

}

func ManageCookieGRPC(ctx context.Context) (uuid string, token string) {
	if flags.FlagDBString != "" {
		uuid, token, err := database.CreateUser(database.DB, ctx)
		if err != nil {
			fmt.Println("MGGRPC error")
			log.Println(err)
			return uuid, token
		} else {
			// uuid := shortuuid.New()
			token, err := users.BuildJWTString(uuid)
			if err != nil {
				log.Println(err)
			}
			return uuid, token
		}
	}

	return uuid, token
}
