package main

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"github.com/ethanjameslong1/GoCloudProject.git/analysis"
// )

// // MockStockData returns mock data for 5 stocks over 9 months
// func RealStockData() []*analysis.StockDataMonthly {
//     symbols := []string{"AAPL", "MSFT", "GOOG", "AMZN", "TSLA",
//     "NVDA", "META", "JPM", "BAC", "WMT",
//     "DIS", "NFLX", "XOM", "KO", "PFE", "SPY"} // SPY needed for beta
//     ctx := context.Background()

//     stocks, err := analysis.MakeMonthlyDataSlice(ctx, symbols)
//     if err != nil {
//         log.Fatalf("Error making monthly data slice: %v", err)
//     }
//     return stocks
// }

// func TestingMain() {
//     // Extract monthly adjusted close prices
//     stocks := RealStockData()

    
//     adjClose := analysis.ExtractMonthlyAdjClosePrices(stocks)
//     returns := analysis.MonthlyStockReturns(adjClose)
//     expected := analysis.ExpectedReturn(returns)
//     stddev := analysis.StandardDeviation(returns)

//     fmt.Println("\n=== Expected Returns & StdDev ===")
//     for symbol := range expected {
//         fmt.Printf("%s: Expected=%.4f, StdDev=%.4f\n", symbol, expected[symbol], stddev[symbol])
//     }

//     analysis.CovarianceMatrixSample(returns)
//     analysis.CorrelationMatrixSample(returns)

//     portfolios, best := analysis.OptimizePortfolio(returns, 500, 0.0, 0.02, 0.20)
//     fmt.Println("\n=== Monte Carlo Portfolio Optimization ===")
//     fmt.Printf("Best Portfolio: %+v\n", best)
//     fmt.Printf("Total Portfolios: %d\n", len(portfolios))

//     marketReturns := returns["SPY"]
//     betas := analysis.BetaCoefficients(returns, marketReturns)
//     fmt.Println("\n=== Beta Coefficients ===")
//     for symbol, beta := range betas {
//         fmt.Printf("%s: Beta=%.4f\n", symbol, beta)
//     }
// }

// func main() {
//     fmt.Println("=== Running Real-Stock Test Workflow ===")
//     TestingMain()
// }
