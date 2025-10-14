package handler

import (
	"errors"
	"github.com/ethanjameslong1/FiNet/cmd/stock/analysis"
	"github.com/ethanjameslong1/FiNet/database"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

const (
	HandlerTimeout = 2 * 24 * time.Hour
)

type apiCall struct {
	symbolList []string
	userID     int
}

type StockHandler struct {
	StockDBService  *database.DBService
	SessionDuration time.Duration
}

func NewHandler(StockDB *database.DBService) (*StockHandler, error) {
	if StockDB == nil {
		return nil, errors.New("database services cannot be nil")
	}
	return &StockHandler{
		StockDBService:  StockDB,
		SessionDuration: HandlerTimeout,
	}, nil
}

func (h *StockHandler) ReceiveAPIcall(c *gin.Context) {
	var apiInfo apiCall
	if err := c.BindJSON(&apiInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"JSON Request Binding Error": err.Error()})
		return
	}
	for _, str := range apiInfo.symbolList {
		println(str)
	}

	dataSlice, err := analysis.MakeWeeklyDataSlice(c, apiInfo.symbolList)
	if err != nil {
		log.Printf("Error creating data slice for analysis: %v", err)
		return
	}

	DataMap, err := analysis.StoreWeeklyDataV1(dataSlice, "", 1)
	if err != nil {
		log.Printf("Error collecting weekly stock data: %v", err)
	}
	Pred, err := analysis.AnalyzeStoredDataV1(DataMap)
	if err != nil {
		log.Printf("Error analyzing stored stock data: %v", err)
	}

	for _, prediction := range Pred {
		log.Printf("AddPrediction begin called with %s (predictable), %s (predictor) and %f (correlation)", prediction.PredictableSym, prediction.PredictorSym, prediction.Correlation)
		h.StockDBService.AddPrediction(c, prediction.PredictableSym, prediction.PredictorSym, prediction.Correlation, "First Draft", apiInfo.userID)
	}

	c.IndentedJSON(http.StatusAccepted, apiInfo.symbolList)
}
