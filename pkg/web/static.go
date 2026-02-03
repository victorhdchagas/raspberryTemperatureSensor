package web

import (
	"net/http"
)

func RegisterStaticRoutes(mux *http.ServeMux) {
	fileServer := http.FileServer(http.Dir("static"))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=300")
		fileServer.ServeHTTP(w, r)
	})

	mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/static/js/app.js" {
			w.Header().Set("Cache-Control", "public, max-age=3600")
		} else if r.URL.Path == "/static/css/htmx.css" {
			w.Header().Set("Cache-Control", "public, max-age=86400")
		}
		fileServer.ServeHTTP(w, r)
	})))
}

func SetNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}
