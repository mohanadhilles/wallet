package main

import (
	"log"
	"net/http"
	"time"
	"wallet/route"

	"github.com/joho/godotenv"
)

func (app *Application) run() error {
	router := route.LoadRouter(&app.store)

	srv := &http.Server{
		Addr:         app.setting.addr,
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Starting server on %s", app.setting.addr)
	return srv.ListenAndServe()
}

func main() {
	godotenv.Load()
	app := bootstrap()
	if err := app.run(); err != nil {
		log.Fatal(err)
	}
}
