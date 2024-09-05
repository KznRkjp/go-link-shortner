package main

import (
	// ...

	"context"
	"fmt"
	"log"

	pb "github.com/KznRkjp/go-link-shortner.git/internal/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// устанавливаем соединение с сервером
	conn, err := grpc.NewClient(":8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := pb.NewHandlersClient(conn)
	// функция, в которой будем отправлять сообщения
	TestUsers(c)
}

func TestUsers(c pb.HandlersClient) {
	// набор тестовых данных
	users := []*pb.ShortenURLRequest{
		{LongURL: "mail.ru"},
	}
	for _, user := range users {
		// добавляем пользователей
		resp, err := c.ShortenURL(context.Background(), &pb.ShortenURLRequest{
			LongURL: user.LongURL,
			// ShortenURLRequest: user,
		})
		if err != nil {
			log.Fatal(err)
		}
		if resp.Error != "" {
			fmt.Println(resp.Error)
		}
		// fmt.Println("response")
		// fmt.Println(resp.ShortURL)
	}

}
