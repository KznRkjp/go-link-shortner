package grpc

import (
	context "context"
	"fmt"

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
	response.Token = "fff"
	response.LongURL = "dfdf"
	return &response, nil

}
