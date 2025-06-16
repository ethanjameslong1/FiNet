// FileSearchServer.go
package main

import (
	"context"
	"fmt"
	"github.com/ethanjameslong1/GoCloudProject.git/database"
	"github.com/ethanjameslong1/GoCloudProject.git/handler"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Uservice, err := database.NewDBService(ctx, database.UserDataSource)
	if err != nil {
		log.Fatal(err)
	}
	defer Uservice.Close()
	Sservice, err := database.NewDBService(ctx, database.SessionDataSource)
	if err != nil {
		log.Fatal(err)
	}
	defer Sservice.Close()

	go func() {
		ticker := time.NewTicker(1 * time.Hour) // Clean up every hour
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done(): // Listen for app shutdown signal
				log.Println("Session cleanup goroutine stopping.")
				return
			case <-ticker.C:
				rowsAffected, err := Sservice.DeleteExpiredSessions(ctx)
				if err != nil {
					log.Printf("Error during session cleanup: %v", err)
				} else {
					log.Printf("Cleaned up %d expired sessions.", rowsAffected)
				}
			}
		}
	}()
	sessionLifetime := 24 * time.Hour
	appHandler, err := handler.NewHandler(Uservice, Sservice, sessionLifetime)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.RootHandler)))
	mux.Handle("GET /login", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.ShowLogin)))
	mux.Handle("GET /adduser", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.ShowAddUser)))
	mux.Handle("GET /stock", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.StockPageHandler)))
	mux.HandleFunc("POST /login", appHandler.LoginHandler)
	mux.HandleFunc("POST /adduser", appHandler.AddUser)
	mux.Handle("POST /stock", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.StockRequestHandler)))

	fmt.Printf("port running on localhost:8080/\n")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}
