// FileSearchServer.go

//testing git authentication

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
		ticker := time.NewTicker(1 * time.Hour) // Clean up every hour
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done(): // Listen for app shutdown signal
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
	mux.HandleFunc("GET /login", appHandler.ShowLogin)
	mux.Handle("GET /homepage", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.HomepageHandler)))
	mux.Handle("GET /stock", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.StockRequestPageHandler)))
	mux.Handle("POST /stock", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.StockRequestHandler)))
	mux.HandleFunc("GET /register", appHandler.ShowRegistration)
	mux.HandleFunc("POST /register", appHandler.RegistrationHandler)
	mux.Handle("GET /predictions", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.ShowPredictionsHandler)))
	mux.Handle("GET /rawdatarequest", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.RawDataRequest)))
	mux.Handle("POST /rawdataview", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.RawDataHandler)))
	mux.HandleFunc("GET /logout", appHandler.LogoutHandler)

	fmt.Printf("port running on localhost:8080/\n")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}
