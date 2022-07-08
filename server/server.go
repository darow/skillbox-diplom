package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/heroku/go-getting-started/simulator"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.staticPath, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func Start(addr string) {
	router := mux.NewRouter()
	router.HandleFunc("/api", getResultHandler())
	router.HandleFunc("/simulator_switch_on", func(w http.ResponseWriter, r *http.Request) {
		go simulator.Start()
	})

	router.HandleFunc("/", serveFiles)
	router.HandleFunc("/status_page.html", serveFiles)
	router.HandleFunc("/chart.min.js", serveFiles)
	router.HandleFunc("/main.js", serveFiles)
	router.HandleFunc("/main.css", serveFiles)
	router.HandleFunc("/static/true.png", serveFiles)
	router.HandleFunc("/static/false.png", serveFiles)

	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func serveFiles(w http.ResponseWriter, r *http.Request) {
	p := "./front" + r.URL.Path
	if p == "./front/" {
		p = "./front/status_page.html"
	}
	http.ServeFile(w, r, p)
}
