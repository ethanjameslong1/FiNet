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
func AnalyzeWeeklyData(d []*StockDataWeekly) error {
	//populating map for faster access to stock data
	symbolDataMap := make(map[string]*StockDataWeekly)
	for _, data := range d {
		symbolDataMap[data.MetaData.Symbol] = data
	}

	var wg sync.WaitGroup
	ch := make(chan Relationship)

	go func(ch <-chan Relationship) {
		for r := range ch {
			fmt.Printf("Collector received relationship: %+v\n", r)
		}
		fmt.Println("Collector: Channel closed, no more relationships to process.")
	}(ch)

	for sym := range symbolDataMap {
		wg.Add(1)
		go func(ch chan<- Relationship, sym string, stockMap map[string]*StockDataWeekly, originalDate string, yearsBack int) {
			defer wg.Done()
			var OriginalFriday string
			if originalDate == "" {
				egFriday, err := findPreviousWeekString("", "", yearsBack)
				if err != nil {
					log.Fatalf("Error finding last friday's date with empty strings: %v", err)
				}
				OriginalFriday = egFriday
			} else {
				egFriday, err := findPreviousWeekString(originalDate, originalDate, yearsBack)
				if err != nil {
					log.Fatalf("Error finding last friday's date with empty strings: %v", err)
				}
				OriginalFriday = egFriday
			}

			tooFarFlag := false
			currentDate := OriginalFriday
			for !tooFarFlag { //for every week that it can this loop will reference every other symbol and try and relate it to the main predictable symbol (sym)
				predictableData := stockMap[sym].TimeSeriesWeekly[currentDate]
				fmt.Printf("Looking at date %s for sym %s, data is %v", currentDate, sym, predictableData)
				currentDate, err := findPreviousWeekString(currentDate, OriginalFriday, yearsBack)
				if err != nil {
					if errors.Is(err, ErrLookbackLimitReached) {
						tooFarFlag = true
						continue
					} else {
						log.Fatalf("Error finding week before %s: %v", currentDate, err)
					}
				}
				for s := range stockMap {
					predictorData := stockMap[s].TimeSeriesWeekly[currentDate]
					fmt.Printf("Checking %s for pattern behaviour", predictorData)
				}

			}

			var significantRelationShip Relationship
			significantRelationShip.predictableSym = "sym"
			ch <- significantRelationShip
		}(ch, sym, symbolDataMap, "", 2)
	}
	wg.Wait()
	close(ch)

	return nil
}

func findPreviousWeekString(originalDate string, recentDate string, maxLookback int) (string, error) {
	//originalDate and recentDate can both be empty to find a week from the most recent friday
	if originalDate == "" && recentDate == "" {
		today := time.Now().Truncate(24 * time.Hour)
		daysAway := (int(today.Weekday()) + 7 - int(time.Friday)) % 7
		originalDate = today.AddDate(0, 0, -daysAway).Format(dateFormat)
		recentDate = originalDate
	} else {
		originalDateTime, err := time.Parse(dateFormat, originalDate)
		if err != nil {
			return "", fmt.Errorf("Error parsing originalDate: %v: %w", originalDate, err)
		}
		daysAway := (int(originalDateTime.Weekday()) + 7 - int(time.Friday)) % 7
		originalDate = originalDateTime.AddDate(0, 0, -daysAway).Format(dateFormat)

		recentDateTime, err := time.Parse(dateFormat, recentDate)
		if err != nil {
			return "", fmt.Errorf("Error parsing recentDate: %v: %w", originalDate, err)
		}
		daysAway = (int(recentDateTime.Weekday()) + 7 - int(time.Friday)) % 7
		recentDate = recentDateTime.AddDate(0, 0, -daysAway).Format(dateFormat)
	}

	//TODO Fix redundancy issues in this code, redefining originalDateTime is redundant ect.

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
		return "", fmt.Errorf("date '%s' is older than lookback limit of %d years from '%s': %w", recentDate, maxLookback, originalDate, ErrLookbackLimitReached)
	}

	previousWeekDateTime := recentDateTime.AddDate(0, 0, -7)

	return previousWeekDateTime.Format(dateFormat), nil
}
