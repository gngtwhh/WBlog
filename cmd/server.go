package main

import (
	"log"

	"github.com/gngtwhh/WBlog/internal/app"
	"github.com/gngtwhh/WBlog/internal/config"
)

func main() {
	if err := config.Load("configs/config.json"); err != nil {
		log.Fatal(err)
	}
	log.Println("config loaded...")
	server := app.NewServer()
	server.Run()
}
