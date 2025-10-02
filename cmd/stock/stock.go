package main

import (
	"github.com/ethanjameslong1/FiNet/cmd/stock/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/item", handler.ReceiveAPIcall)

	router.Run(":6969")
}
