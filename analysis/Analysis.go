package analysis

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Test Values
var AlphaVantageSymbols = []string{"MSFT", "GOOG", "TSLA", "AMZN"}
var ErrLookbackLimitReached = errors.New("lookback limit reached")

const (
	lookBackTimeYear = 10
	dateFormat       = "2006-01-02"
)

type Relationship struct {
	predictableSym       string
	predictableUUID      uuid.UUID
	predictedSlope       float32
	predictorSyms        []string
	predictorUUIDs       []uuid.UUID
	predictorParam       []string
	predictorParamDeltas []float32
}

// Makes a custom Personalized table in a user database, returns nil when successful
func AnalyzeWeeklyData(d []*StockDataWeekly, startDate string) error {
	//populating map for faster access to stock data
	symbolDataMap := make(map[string]*StockDataWeekly)
	for _, data := range d {
		symbolDataMap[data.MetaData.Symbol] = data
	}

	var wg sync.WaitGroup
	ch := make(chan Relationship)

	//COLLECTOR FUNC *********************************************************************
	go func(ch <-chan Relationship) { // *
		for r := range ch { // *
			fmt.Printf("Collector received relationship: %+v\n", r) // *
		} // *
		fmt.Println("Collector: Channel closed, no more relationships to process.") // *
	}(ch) // *
	//COLLECTOR FUNC END *****************************************************************

	for sym := range symbolDataMap {
		wg.Add(1) //for evecy symbol spawn a seperate thread to analysis it for predictors
		go func(ch chan<- Relationship, sym string, stockMap map[string]*StockDataWeekly, startDate string, yearsBack int) {
			defer wg.Done()
			//DATE SETUP **********************************************************************
			var OriginalFriday string
			if startDate == "" {
				egFriday := findMostRecentFriday()
				OriginalFriday = egFriday
			} else {
				egFriday, err := findPriorFriday(startDate)
				if err != nil {
					log.Printf("Error finding last friday's date: %v", err)
					return
				}
				OriginalFriday = egFriday
			}
			//DATE SETUP **********************************************************************

			tooFarFlag := false
			currentDate := OriginalFriday //establish currentDate as the original date to go back from
			var err error
			for !tooFarFlag { //for every week that it can this loop will reference every other symbol and try and relate it to the main predictable symbol (sym), will be multithread as well

				predictableData := stockMap[sym].TimeSeriesWeekly[currentDate] //predictable data should be 1 week ahead of predictor data
				fmt.Printf("Looking at date %s for sym %s, data is %+v", currentDate, sym, predictableData)
				currentDate, err = findPreviousWeekString(OriginalFriday, currentDate, yearsBack)

				if err != nil {
					if errors.Is(err, ErrLookbackLimitReached) {
						fmt.Println("Lookback Limit Reached")
						tooFarFlag = true
						continue
					} else {
						log.Printf("Error finding Week Before %s: %v", currentDate, err)
						return
					}
				}

				for s := range stockMap {
					predictorData := stockMap[s].TimeSeriesWeekly[currentDate]
					log.Printf("Checking %s for pattern behaviour", s)
					log.Printf("PredictorData: %+v", predictorData)
				}

			}

			var significantRelationShip Relationship
			significantRelationShip.predictableSym = "sym"
			ch <- significantRelationShip
		}(ch, sym, symbolDataMap, startDate, 1)
	}

	wg.Wait()
	close(ch)

	return nil
}

func findPreviousWeekString(originalDate string, recentDate string, maxLookback int) (string, error) {
	originalDateTime, err := time.Parse(dateFormat, originalDate)
	if err != nil {
		return "", fmt.Errorf("Error parsing originalDate: %v: %w", originalDate, err)
	}

	recentDateTime, err := time.Parse(dateFormat, recentDate)
	if err != nil {
		return "", fmt.Errorf("Error parsing recentDate: %v: %w", originalDate, err)
	}

	maxLookbackThreshold := originalDateTime.AddDate(-maxLookback, 0, 0)
	if recentDateTime.Before(maxLookbackThreshold) {
		return "", fmt.Errorf("date '%s' is older than lookback limit of %d years from '%s': %w", recentDate, maxLookback, originalDate, ErrLookbackLimitReached)
	}

	previousWeekDateTime := recentDateTime.AddDate(0, 0, -7)

	return previousWeekDateTime.Format(dateFormat), nil
}

func findMostRecentFriday() string {
	today := time.Now().Truncate(24 * time.Hour)
	daysAway := (int(today.Weekday()) + 7 - int(time.Friday)) % 7
	recentFriday := today.AddDate(0, 0, -daysAway).Format(dateFormat)
	return recentFriday
}

func findPriorFriday(startDate string) (string, error) {
	today, err := time.Parse(dateFormat, startDate)
	if err != nil {
		return "", fmt.Errorf("Error parsing startDate %s: %w", startDate, err)
	}
	today = today.Truncate(24 * time.Hour)
	daysAway := (int(today.Weekday()) + 7 - int(time.Friday)) % 7
	recentFriday := today.AddDate(0, 0, -daysAway).Format(dateFormat)
	return recentFriday, nil
}
