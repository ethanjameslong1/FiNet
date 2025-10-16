package main

import (
	"context"
	"github.com/ethanjameslong1/FiNet/cmd/stock/handler"
	"github.com/ethanjameslong1/FiNet/database"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

const (
	ctxTimeout = 24 * time.Hour
)

func main() {

	router := gin.Default()
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	servStockDB, err := database.NewDBService(ctx, database.StockDataSource)
	if err != nil {
		log.Fatal(err)
	}
	defer servStockDB.Close()
	stockHandler, err := handler.NewHandler(servStockDB)
	if err != nil {
		log.Fatal(err)
	}

	//TEST this is just to see if the 2 containers can communicate, the function is in requestHandler.go. In ../finet/main.go there is more info
	router.POST("/analysis/item", stockHandler.ReceiveAPIcall)

	router.Run(":9090")
}
