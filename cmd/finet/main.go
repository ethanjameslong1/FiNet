// FileSearchServer.go

//testing git authentication

package main

import (
	"context"
	"fmt"
	"github.com/ethanjameslong1/FiNet/cmd/finet/handler"
	"github.com/ethanjameslong1/FiNet/database"
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
	appHandler.UserSessionDBService.AddUser(ctx, "ethan", "test")

	mux := http.NewServeMux()
	mux.Handle("/", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.RootHandler)))
	mux.HandleFunc("POST /login", appHandler.LoginHandler)
	mux.Handle("GET /homepage", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.HomepageHandler)))
	mux.HandleFunc("GET /stock", http.HandlerFunc(appHandler.StockRequestPageHandler))
	mux.HandleFunc("POST /stock", http.HandlerFunc(appHandler.StockRequestHandler))
	mux.HandleFunc("GET /register", appHandler.ShowRegistration)
	mux.HandleFunc("POST /register", appHandler.RegistrationHandler)
	mux.Handle("GET /predictions", appHandler.AuthMiddleware(http.HandlerFunc(appHandler.ShowPredictionsHandler))) // PROD: requires sessionManagement, not removing middleware
	mux.HandleFunc("GET /rawdata", appHandler.RawDataRequest)
	mux.HandleFunc("POST /rawdata", appHandler.RawDataHandler)
	mux.HandleFunc("GET /logout", appHandler.LogoutHandler)

	//TEST this is just to ensure that once nginx is in place we can properly move analysis logic over. Functions are in loginhandler.go, just change the struct there to test whatever
	mux.HandleFunc("GET /testapi", appHandler.TESTAPISTOCK)
	mux.HandleFunc("POST /itemtest", appHandler.TESTAPISTOCKhandle)

	fmt.Printf("PROD: finet port running on app_network: finet:8000/\n")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatal(err)
	}

}
