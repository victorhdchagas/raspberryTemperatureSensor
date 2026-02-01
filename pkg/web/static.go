package web

import (
	"net/http"
)

func RegisterStaticRoutes(mux *http.ServeMux) {
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/", fileServer)
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
}
