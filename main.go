package main

import (
	"awesomeProject/data"
	"awesomeProject/handlers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {

	l := log.New(os.Stdout, "awsome-project ", log.LstdFlags)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	const uri = "mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir-case-study?retryWrites=true"

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	l.Println("Successfully connected and pinged to MongoDb")
	recordCollection := client.Database("getir-case-study").Collection("records")

	activeTabRepository := data.NewActiveTabsInMemoryRepository()
	recordsRepository := data.NewMongoRecordRepository(recordCollection, ctx)

	recordsHandler := handlers.NewRecord(l, recordsRepository)
	inMemoryHandler := handlers.NewInMemory(l, activeTabRepository)

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

	s.Shutdown(ctx)
}
