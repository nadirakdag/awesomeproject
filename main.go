package main

import (
	"awesomeProject/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	l := log.New(os.Stdout, "awsome-project ", log.LstdFlags)

	recordsHandler := handlers.NewRecord(l)
	inMemoryHandler := handlers.NewInMemory(l)

	sm := http.NewServeMux()
	sm.Handle("/api/records", recordsHandler)
	sm.Handle("/api/in-memory", inMemoryHandler)

	serverPortEnv := os.Getenv("PORT")
	var serverPort string

	if serverPortEnv != "" {
		l.Printf("PORT found in ENV " + serverPortEnv)
		serverPort = serverPortEnv
	} else {
		l.Printf("PORT not found in ENV, gonna use default 8080")
		serverPort = "8080"
	}

	// create a new server
	s := http.Server{
		Addr:         ":" + serverPort,
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Println("Starting server on port " + serverPort)

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal:", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
