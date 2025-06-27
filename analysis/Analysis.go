package analysis

import (
	"fmt"
	"sync"
	"time"
)

// Test Values
var AlphaVantageSymbols = []string{"MSFT", "GOOG", "TSLA", "AMZN"}

const (
	lookBackTimeYear = 10
	dateFormat       = "2006-01-02"
)

type Relationship struct {
	predictableSym string
	predictors     []string
}

// Makes a custom Personalized table in a user database, returns nil when successful
func AnalyzeWeeklyData(d []*StockDataWeekly) error {
	//populating map for faster access to stock data
	symbolDataMap := make(map[string]*StockDataWeekly)
	for _, data := range d {
		symbolDataMap[data.MetaData.Symbol] = data
	}

	var wg sync.WaitGroup
	ch := make(chan Relationship)

	for sym := range symbolDataMap {
		wg.Add(1)
		go func(c chan Relationship, sym string) {
			var significantRelationShip Relationship
			//is this structure going to work, I plan on looking for relationships and if I find them I will pipe them to ch
			//after the funcs all finish I will use the channel to update an external database that wil hold important information

			significantRelationShip.predictableSym = "sym" //this is just so it doesn't throw errors
			ch <- significantRelationShip                  //this is just so it doesn't throw errors

		}(ch, sym)

	}
	wg.Wait()

	return nil

}

func findPreviousWeekString(originalDate string, recentDate string, maxLookback int) (string, error) {
	originalDateTime, err := time.Parse(dateFormat, originalDate)
	if err != nil {
		return "", fmt.Errorf("Error Parsing originalDate: %v: %w", originalDate, err)
	}
	recentDateTime, err := time.Parse(dateFormat, recentDate)
	if err != nil {
		return "", fmt.Errorf("Error Parsing recentDate: %v: %w", recentDate, err)
	}
	maxLookbackThreshold := originalDateTime.AddDate(-maxLookback, 0, 0)
	if recentDateTime.Before(maxLookbackThreshold) {
		return "", fmt.Errorf("recent date '%s' is more than %d years before original date '%s'", recentDate, maxLookback, originalDate)
	}

	previousWeekDateTime := recentDateTime.AddDate(0, 0, -7)

	return previousWeekDateTime.Format(dateFormat), nil
}
