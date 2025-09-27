package analysis

import (
	"gonum.org/v1/gonum/stat"
)

// MonthlyStockReturns takes in a slice of StockDataMonthly structs and returns a map of stock symbols to their monthly return percentages
func MonthlyStockReturns(adjClose map[string][]float64) map[string][]float64 {
    returns := make(map[string][]float64)

    for ticker, prices := range adjClose {
        if len(prices) < 2 {
            continue // Need at least 2 points to compute return
        }

        var monthlyReturns []float64
        for i := 1; i < len(prices); i++ {
            ret := (prices[i] - prices[i-1]) / prices[i-1] * 100
            monthlyReturns = append(monthlyReturns, ret)
        }

        returns[ticker] = monthlyReturns
    }

    return returns
}

//Will use SPY as market benchmark
func BetaCoefficients(returns map[string][]float64, marketReturns []float64) map[string]float64 {
	betas := make(map[string]float64)

	for ticker, stockReturns := range returns {
		if len(stockReturns) != len(marketReturns) || len(stockReturns) < 2 {
			continue // Ensure both slices are of equal length and have enough data points
		}

		// Calculate covariance and variance
		covariance := stat.Covariance(stockReturns, marketReturns, nil)
		variance := stat.Variance(marketReturns, nil)

		if variance != 0 {
			beta := covariance / variance
			betas[ticker] = beta
		}
	}

	return betas

}


//Using arithmetic mean for simplicity, geometric mean can be implemented if needed
func ExpectedReturn(returns map[string][]float64) map[string]float64 {
	expected := make(map[string]float64)

	for ticker, stockReturns := range returns {
		if len(stockReturns) == 0 {
			continue
		}

		// Calculate the arithmetic mean
		mean := stat.Mean(stockReturns, nil)
		expected[ticker] = mean
	}

	return expected
}

func StandardDeviation(returns map[string][]float64) map[string]float64 {
	stdDevs := make(map[string]float64)

	for ticker, stockReturns := range returns {
		if len(stockReturns) == 0 {
			continue
		}

		// Calculate the standard deviation
		stdDev := stat.StdDev(stockReturns, nil)
		stdDevs[ticker] = stdDev
	}

	return stdDevs
}