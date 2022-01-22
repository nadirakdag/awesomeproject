package main

import (
	"awesomeProject/data"
	"awesomeProject/db"
	"awesomeProject/handlers"
	"awesomeProject/helpers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// applicaation start point
func main() {

	l := log.New(os.Stdout, "awsome-project ", log.LstdFlags)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	const uri = "mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir-case-study?retryWrites=true"
	dbConnection, err := db.NewMongoConnection(uri, helpers.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = dbConnection.Db.Client().Disconnect(context.Background())
	}()

	keyValuePairsRepository, recordsRepository := initRepositories(l, dbConnection)

	recordsHandler := handlers.NewRecord(l, recordsRepository)
	inMemoryHandler := handlers.NewInMemory(l, keyValuePairsRepository)

	serverPort := getServerPort(l)
	server := createdAndStartHttpServer(l, recordsHandler, inMemoryHandler, serverPort)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal:", sig)

	server.Shutdown(ctx)
}

// Creates http server and starts it
func createdAndStartHttpServer(l *log.Logger, recordsHandler *handlers.Records, inMemoryHandler *handlers.InMemory, serverPort string) http.Server {
	sm := http.NewServeMux()
	sm.Handle("/api/records", recordsHandler)
	sm.Handle("/api/in-memory", inMemoryHandler)

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

	return s
}

// get Server port from Enviroment Variables
func getServerPort(l *log.Logger) string {
	serverPortEnv := os.Getenv("PORT")
	var serverPort string

	if serverPortEnv != "" {
		l.Printf("PORT found in ENV " + serverPortEnv)
		serverPort = serverPortEnv
	} else {
		l.Printf("PORT not found in ENV, gonna use default 8080")
		serverPort = "8080"
	}
	return serverPort
}

// inits and returns KeyValuePairRepository and RecordRepository
func initRepositories(l *log.Logger, db *db.MongoConnection) (data.KeyValueRepository, data.RecordRepository) {

	recordCollection := db.Db.Collection(helpers.CollectionName)

	keyValuePairRepository := data.NewKeyValueInMemoryRepository()
	recordsRepository := data.NewMongoRecordRepository(recordCollection)

	return keyValuePairRepository, recordsRepository
}
