package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ReceiveAPIcall(c *gin.Context) {
	var symbolList []string
	if err := c.BindJSON(&symbolList); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, str := range symbolList {
		println(str)
	}
	c.IndentedJSON(http.StatusAccepted, symbolList)
}

/*
func receiveAPIcall()

dataSlice, err := analysis.MakeWeeklyDataSlice(r.Context(), symbolList)
if err != nil {
log.Printf("Error creating data slice for analysis: %v", err)
http.Error(w, "Bad Request", http.StatusBadRequest)
return
}

DataMap, err := analysis.StoreWeeklyDataV1(dataSlice, "", 1)
if err != nil {
log.Printf("Error collecting weekly stock data: %v", err)
http.Error(w, "Internal server error", http.StatusInternalServerError)
}
Pred, err := analysis.AnalyzeStoredDataV1(DataMap)
if err != nil {
log.Printf("Error analyzing stored stock data: %v", err)
http.Error(w, "Internal server error", http.StatusInternalServerError)
}

for _, prediction := range Pred {
log.Printf("AddPrediction begin called with %s (predictable), %s (predictor) and %f (correlation)", prediction.PredictableSym, prediction.PredictorSym, prediction.Correlation)
h.StockDBService.AddPrediction(r.Context(), prediction.PredictableSym, prediction.PredictorSym, prediction.Correlation, "First Draft", user.ID)
}


*/
