package components

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"push_app/configs"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type AdHocListener struct {
	QueryChan chan AdHocQuery
}

type AdHocQuery struct {
	Query string
	Keys  []string
	Topic string
}

func (l *AdHocListener) handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %v request for %s\n", r.Method)
	decoder := json.NewDecoder(r.Body)
	var t AdHocQuery
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	go func() { l.QueryChan <- t }()
	/*
		// Send req using http Client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
		}
		//defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while reading the response bytes:", err)
		}
		//log.Println(string([]byte(body)))
		w.Write([]byte(body))

		w.Write([]byte(fmt.Sprintf("Hello, %s\n", "Jingyu Su")))
	*/
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (l *AdHocListener) run() {
	// Create Server and Route Handlers
	r := mux.NewRouter()

	r.HandleFunc("/", l.handler)
	r.HandleFunc("/health", healthHandler)
	r.HandleFunc("/readiness", readinessHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}

func MakeListener(conf configs.OperatorConf) *AdHocListener {
	return &AdHocListener{QueryChan: make(chan AdHocQuery, conf.ChanBufSize)}
}

func (l *AdHocListener) GetQueryChan() chan AdHocQuery {
	return l.QueryChan
}
