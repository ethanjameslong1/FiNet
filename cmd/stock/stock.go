package main

import (
	"context"
	"log"
	"time"

	"github.com/ethanjameslong1/FiNet/cmd/stock/handler"
	"github.com/ethanjameslong1/FiNet/database"
	"github.com/gin-gonic/gin"
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

	router.POST("/rawDataRequest", stockHandler.RawDataAPIRequest)
	router.POST("/analysisRequest", stockHandler.AnalysisAPIRequest)

	router.Run(":8001")
}
