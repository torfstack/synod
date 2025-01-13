package main

import "main/internal/http"

func main() {
	server := http.NewServer()
	server.Start()
}
