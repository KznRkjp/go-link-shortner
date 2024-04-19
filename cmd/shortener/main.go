package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/KznRkjp/go-link-shortner.git/internal/flags"
	"github.com/KznRkjp/go-link-shortner.git/internal/router"
)

func main() {
	flags.ParseFlags()
	dd := router.Main()
	fmt.Println("Server is listening @", flags.FlagRunAddr)
	fmt.Println("Press Ctrl+C to stop")
	log.Fatal(http.ListenAndServe(flags.FlagRunAddr, dd))
}
