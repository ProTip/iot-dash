package main

import (
	"crypto/tls"
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	migrationBytes, ioErr := ioutil.ReadFile("./migrate.sql")
	if ioErr != nil {
		log.Fatal("SQL schema file [migrate.sql] missing!")
	}

	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(string(migrationBytes))
	if err != nil {
		log.Fatal(err)
	}

	repo := AppRepo{db}
	app := App{AppRepo: repo, sessions: &sync.Map{}}

	mux := http.NewServeMux()

	mux.Handle(
		"/security/login",
		HandleCsrf(true,
			HandleSecurityLogin(app)))

	mux.Handle(
		"/security/logout",
		HandleCsrf(true,
			HandleSecurityLogout(app)))

	mux.HandleFunc(
		"/account/upgrade",
		HandleAuth(app,
			HandleAccountUpgrade(app)))

	mux.HandleFunc(
		"/metrics",
		HandleAuth(
			app,
			HandleMetrics(app)))

	/* Serve static SPA assets from filesystem */
	mux.Handle(
		"/static/",
		http.FileServer(http.Dir("ui/")))

	/*
	  Serve the SPA index for any other path
	  TODO API routes would be better placed under /api
	*/
	mux.HandleFunc(
		"/",
		HandleCsrf(false,
			http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Add("Strict-Transport-Security", "max-age=63072000")
				http.ServeFile(w, req, "ui/index.html")
			})))

	cfg := &tls.Config{
		MinVersion: tls.VersionTLS13,
	}

	srv := &http.Server{
		Addr:      ":8000",
		Handler:   mux,
		TLSConfig: cfg,
		/*
			TODO consider the below Mozilla recommendation to mitigate "low and slow" attacks
			as well as limiting request sizes
		*/
		// Consider setting ReadTimeout, WriteTimeout, and IdleTimeout
		// to prevent connections from taking resources indefinitely.
	}

	log.Fatal(srv.ListenAndServeTLS(
		"ssl/cert.pem",
		"ssl/key.pem",
	))

}
