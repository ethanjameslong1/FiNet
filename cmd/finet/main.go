// FileSearchServer.go

// testing git authentication
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ethanjameslong1/FiNet/cmd/finet/handler"
	"github.com/ethanjameslong1/FiNet/database"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	servUSDB, err := database.NewDBService(ctx, database.UserSessionDataSource)
	if err != nil {
		log.Fatal(err)
	}
	defer servUSDB.Close()
	servStockDB, err := database.NewDBService(ctx, database.StockDataSource)
	if err != nil {
		log.Fatal(err)
	}
	defer servStockDB.Close()

	go func() {
		ticker := time.NewTicker(1 * time.Hour) // Clean up sessions every hour
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Println("Session cleanup goroutine stopping.")
				return
			case <-ticker.C:
				rowsAffected, err := servUSDB.DeleteExpiredSessions(ctx)
				if err != nil {
					log.Printf("Error during session cleanup: %v", err)
				} else {
					log.Printf("Cleaned up %d expired sessions.", rowsAffected)
				}
			}
		}
	}()

	sessionLifetime := 24 * time.Hour
	appHandler, err := handler.NewHandler(servUSDB, servStockDB, sessionLifetime)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.RootHandler)))
	mux.HandleFunc("POST /login", appHandler.LoginHandler)
	mux.Handle("GET /homepage", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.HomepageHandler)))
	mux.Handle("GET /stock", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.StockRequestPageHandler)))
	mux.Handle("POST /stock", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.StockRequestHandler)))
	mux.HandleFunc("GET /register", appHandler.ShowRegistration)
	mux.HandleFunc("POST /register", appHandler.RegistrationHandler)
	mux.Handle("GET /predictions", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.ShowPredictionsHandler))) // PROD: requires sessionManagement, not removing middleware
	mux.Handle("GET /rawdata", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.RawDataRequest)))
	mux.Handle("POST /rawdata", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.RawDataHandler)))
	mux.HandleFunc("GET /logout", appHandler.LogoutHandler)
	mux.HandleFunc("POST /middleware", appHandler.Middleware)
	mux.Handle("GET /clearpredictions", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.ClearPredictions)))

	fmt.Printf("PROD: finet port running on app_network: finet:8000/\n")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatal(err)
	}
}
