package main

import "github.com/gngtwhh/WBlog/internal/app"

func main() {
	server := app.NewServer()
	server.Run()
}
