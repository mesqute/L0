package controllers

import "net/http"

func GetRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/data", getData)
	mux.HandleFunc("/data/create", setData)

	return mux
}
