package server

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "\\/")
	path = "./" + h.staticPath + "/" + path
	http.ServeFile(w, r, path)
}

func Start(addr string) {
	router := mux.NewRouter()
	router.HandleFunc("/api", getResultHandler())

	spa := spaHandler{staticPath: "web", indexPath: "status_page.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
