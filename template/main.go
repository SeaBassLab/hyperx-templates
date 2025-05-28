package main

import (
	"log"
	"os"

	server "github.com/SeaBassLab/hyperx-server"
)

func main() {
	env := os.Getenv("HYPERX_ENV") // puede ser "dev" o "prod"
	if env == "" {
		env = "dev" // fallback por defecto
	}

	err := server.StartServer(env, "3000")
	if err != nil {
		log.Fatal(err)
	}
}
