package main

import (
	"chat-apps/internal/infrastructure"

	_ "github.com/swaggo/gin-swagger"
	_ "github.com/swaggo/swag"
)

func main() {
	r := infrastructure.NewRouter()
	r.Run(":8080")
}
