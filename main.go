package main

import (
	"L0/controllers"
	"log"
	"net/http"
)

func main() {

	server := &http.Server{
		Addr:    ":8080",
		Handler: controllers.GetRoutes(),
	}

	// запуск сервера
	log.Printf("Запуск сервера на %s", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("[main] ListenAndServe: %s", err)
	}
}
