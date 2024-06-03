package users

import (
	"fmt"
	"log"
	"net/http"
	"time"

	// "github.com/KznRkjp/go-link-shortner.git/internal/database"
	// "github.com/KznRkjp/go-link-shortner.git/internal/database"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserUID string
}

const TokenExp = time.Hour * 10
const SecretKey = "supersecretkey"

func BuildJWTString(uuid string) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		// собственное утверждение
		UserUID: uuid,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func GetUserUID(tokenString string) (string, error) {
	fmt.Println("****** starting jwt check")
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Println("Тут что то не так")
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		fmt.Println(err)
		return claims.UserUID, err
	}

	if !token.Valid {
		fmt.Println("Token is not valid")
		return claims.UserUID, err
	}

	fmt.Println("Token is valid")
	return claims.UserUID, err
}

func Access(req *http.Request) (string, error) {
	jwt, err := req.Cookie("JWT")
	if err != nil {
		fmt.Println("cookie error")
		return "", err
	}
	uuid, err := GetUserUID(jwt.Value)
	log.Println("Access checked", uuid)
	return uuid, err

}
