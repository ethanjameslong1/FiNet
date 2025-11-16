package analysis_test

import (
	"context"
	"encoding/json"
	"github.com/ethanjameslong1/FiNet/cmd/stock/analysis"
	"testing"
	"time"
)

func TestMakeWeeklyDataSlice(t *testing.T) {
	var AlphaVantageSymbols = []string{"MSFT", "GOOG", "TSLA", "AMZN"}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Logf("Attempting to retrieve weekly stock data for symbols %v", AlphaVantageSymbols)
	dataSlice, err := analysis.MakeWeeklyDataSlice(ctx, AlphaVantageSymbols)
	if err != nil {
		t.Fatalf("MakeWeeklyDataSlice returned an error: %v", err)
	}
	if dataSlice == nil {
		t.Log("MakeWeeklyDataSlice returned a nil value")
		return
	}
	for i, stockData := range dataSlice {
		if stockData == nil {
			t.Logf("Data for symbol %v is NIL", AlphaVantageSymbols[i])
			continue
		}
		jsonData, marshalErr := json.MarshalIndent(stockData, "", "  ")
		if marshalErr != nil {
			t.Errorf("Error marsheling StockDataWeekly for symbol %q: %v", stockData.MetaData.Symbol, marshalErr)
			continue
		}
		t.Logf("\n---Retrieved Stock Data for Symbol: %s (index %d) ---", stockData.MetaData.Symbol, i)
		t.Logf("%s", jsonData)
		t.Log("-------------------------------------------------------")

		time.Sleep(1 * time.Second)
	}

}
