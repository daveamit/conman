package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func handler(name string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		json.NewEncoder(w).Encode(map[string]string{"message": "Not Implemented", "endpoint": name})
	}
}

func makeRouter() *mux.Router {
	r := mux.NewRouter()

	// Routes

	// Auth
	r.HandleFunc("/login", handler("POST /login")).
		Methods("POST")

	// User management
	// Retrive users
	r.HandleFunc("/users", handler("GET /users")).
		Methods("GET")
	r.HandleFunc("/users/{username}", handler("users/{username}")).
		Methods("GET")

	// Add users
	r.HandleFunc("/users", handler("POST /users")).
		Methods("POST")
	// Update user
	r.HandleFunc("/users/{username}", handler("PATCH /users/{username}")).
		Methods("PATCH")

	// Remove user
	r.HandleFunc("/users/{username}", handler("DELETE /users/{username}")).
		Methods("DELETE")

	// Setting management
	// Get all settings
	r.HandleFunc("/settings", handler("GET /settings?filter={filter}")).
		Queries("filter", "{filter}").
		Methods("GET")
	// Get setting by pattern
	r.HandleFunc("/settings/{key}", handler("GET /settings/{key}")).
		Methods("GET")
	// Delete a setting
	r.HandleFunc("/settings/{username}", handler("DELETE /settings/{username}")).
		Methods("DELETE")

	return r
}

// Serve serves http server
func Serve(ctx context.Context) {

	writeTimeout := 15 * time.Second
	readTimeout := 15 * time.Second
	address := fmt.Sprintf("%s:%d", inet, port)
	log := log.With(zap.String("address", address))
	srv := &http.Server{
		Handler: makeRouter(),
		Addr:    address,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}

	log.Info("Started")
	exited := make(chan bool)

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.With(zap.Error(err)).Fatal("Exited")
		}
		close(exited)
	}()

	<-ctx.Done()
	log.Info("Context cancled")
	srv.Close()
	log.Info("Server closed, listenAndServe will exit soon")

	<-exited
	log.Info("Server exited")
}
func main() {
	defineConfigFlags()
	flag.Parse()

	setLog(stringToLogLevel(verbosity), "json")

	ctx, cancel := context.WithCancel(context.Background())

	go Serve(ctx)

	// Wait for os Interrupt signal (usually CTRL+C)
	done := make(chan os.Signal)
	go signal.Notify(done, os.Interrupt, os.Kill)

	<-done
	cancel()
}
