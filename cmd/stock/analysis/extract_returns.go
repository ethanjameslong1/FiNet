package analysis

import (
	"log"
	"sort"
	"strconv"
)

func ExtractMonthlyAdjClosePrices(data []*StockDataMonthly) map[string][]float64 {
	adjClosePrices := make(map[string][]float64) // Map of stock symbol to slice of close prices

	// Iterate over each stock's data
	for _, stock := range data {
		if len(stock.TimeSeriesMonthly) == 0 {
			continue
		}

		// Extract and sort the dates to ensure chronological order
		var dates []string
		for date := range stock.TimeSeriesMonthly {
			dates = append(dates, date)
		}
		sort.Strings(dates)
		//TODO check if this is needed, for performance reasons. //sorted dates are needed to ensure chronological order for calculations

		// Collect close prices in chronological order
		var prices []float64
		for _, date := range dates {
			closeStr := stock.TimeSeriesMonthly[date].AdjClose
			closeVal, err := strconv.ParseFloat(closeStr, 64)
			if err != nil {
				log.Printf("Error parsing close price for %s on %s: %v", stock.MetaData.Symbol, date, err)
				continue
			}
			prices = append(prices, closeVal) // Append to prices slice
		}
		adjClosePrices[stock.MetaData.Symbol] = prices // Store in map
		//TODO should probably make sure date is somehow available from the map, maybe making a small struct? // dates are inherently available because they are sorted, shouldnt be needed.
	}		

	return adjClosePrices
}