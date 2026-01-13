// Package handler comment to make linter happy
package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ethanjameslong1/FiNet/cmd/stock/analysis"
	"github.com/ethanjameslong1/FiNet/database"
	"github.com/gin-gonic/gin"
)

const (
	HandlerTimeout = 2 * 24 * time.Hour
)

type apiCall struct { // any changes made here need to be reflected in finet Go web server
	SymbolList []string `json:"symbolList"`
	TimePeriod string   `json:"time"`
	UserID     int      `json:"Id"`
}
type anlysisCall struct { // any changes made here need to be reflected in finet Go web server
	SymbolList []string `json:"symbolList"`
	UserID     int      `json:"Id"`
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

func (h *StockHandler) RawDataAPIRequest(c *gin.Context) {
	var apiInfo apiCall
	if err := c.BindJSON(&apiInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"JSON Request Binding Error": err.Error()})
		return
	}
	fmt.Printf("ApiInfo: %v\n", apiInfo)
	stockData, err := analysis.MakeWeeklyDataSlice(c, apiInfo.SymbolList, apiInfo.TimePeriod)
	if err != nil {
		log.Printf("Error creating data slice for analysis: %v", err)
		return
	}

	c.IndentedJSON(http.StatusCreated, stockData)
}

func (h *StockHandler) AnalysisAPIRequest(c *gin.Context) {
	var apiInfo anlysisCall
	if err := c.BindJSON(&apiInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"JSON Request Binding Error": err.Error()})
		return
	}
	fmt.Printf("ApiInfo: %v\n", apiInfo.SymbolList)

	dataSlice, err := analysis.AnalysisStoreWeeklyDataSlice(c, apiInfo.SymbolList)
	if err != nil {
		log.Printf("Error creating data slice for analysis: %v", err)
		c.IndentedJSON(http.StatusBadRequest, apiInfo.SymbolList)
		return
	}
	for i := range dataSlice {
		fmt.Printf("\ndataSlice: %+v\n", dataSlice[i].MetaData.Symbol)
	}

	DataMap, err := analysis.StoreWeeklyDataV1(dataSlice, "", 1)
	if err != nil {
		log.Printf("Error collecting weekly stock data: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, apiInfo.SymbolList)
	}
	Pred, err := analysis.AnalyzeStoredDataV1(DataMap)
	if err != nil {
		log.Printf("Error analyzing stored stock data: %v", err)
		c.IndentedJSON(http.StatusBadRequest, apiInfo.SymbolList)
	}

	for _, prediction := range Pred {
		log.Printf("AddPrediction begin called with %s (predictable), %s (predictor) and %f (correlation)", prediction.PredictableSym, prediction.PredictorSym, prediction.Correlation)
		h.StockDBService.AddPrediction(c, prediction.PredictableSym, prediction.PredictorSym, prediction.Correlation, "First Draft", apiInfo.UserID)
	}

	c.IndentedJSON(http.StatusCreated, apiInfo.SymbolList)
}
