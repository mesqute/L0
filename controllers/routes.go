package controllers

import "net/http"

func GetRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// обработчик запросов статических файлов (css, js, image и т.д.)
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/data", getData)
	mux.HandleFunc("/data/create", setData)

	return mux
}
