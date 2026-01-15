package analysis

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"gonum.org/v1/gonum/stat"
)

// Test Values
var (
	AlphaVantageSymbols     = []string{"MSFT", "GOOG", "TSLA", "AMZN"} // TODO: probably no longer needed
	ErrLookbackLimitReached = errors.New("lookback limit reached")
)

const (
	lookBackTimeYear = 10 // TODO: probably no longer needed
	dateFormat       = "2006-01-02"
)

type RelationshipData struct {
	PredictableSym        string
	PredictorSym          string
	PredictableCloseDelta float64
	PredictorCloseDelta   float64
}
type RelationshipKey struct {
	PredictorSym   string
	PredictableSym string
}
type RelationshipPackage struct {
	Key  RelationshipKey
	Data RelationshipData
}
type Prediction struct {
	PredictableSym string
	PredictorSym   string
	Correlation    float64
}

// Creates a hashmap of relationalData, returns the map and an error (if any)
func StoreWeeklyDataV1(d []*StockDataWeekly, startDate string, lookBackTime int) (map[RelationshipKey][]RelationshipData, error) {
	symbolDataMap := make(map[string]*StockDataWeekly)
	for _, data := range d {
		symbolDataMap[data.MetaData.Symbol] = data
	}
	relationshipMap := make(map[RelationshipKey][]RelationshipData)
	var mapMutex sync.Mutex

	var wg sync.WaitGroup
	var collectorwg sync.WaitGroup
	ch := make(chan RelationshipPackage)

	collectorwg.Add(1)
	// COLLECTOR FUNC *********************************************************************
	go func(ch <-chan RelationshipPackage) {
		defer collectorwg.Done()

		for pkg := range ch {
			mapMutex.Lock()
			relationshipMap[pkg.Key] = append(relationshipMap[pkg.Key], pkg.Data)
			mapMutex.Unlock()
			// fmt.Printf("Collector recieved relationship: %+v\n", pkg)
		}
		fmt.Println("Collector: Channel closed, no more relationships to process.")
	}(ch)
	// COLLECTOR FUNC END *****************************************************************

	// DATE SETUP **********************************************************************
	var OriginalFriday string
	if startDate == "" {
		egFriday := findMostRecentFriday()
		OriginalFriday = egFriday
	} else {
		egFriday, err := findPriorFriday(startDate)
		if err != nil {
			log.Printf("Error finding last friday's date: %v", err)
			return nil, fmt.Errorf("Error finding last friday's date: %w", err)
		}
		OriginalFriday = egFriday
	}
	// DATE SETUP **********************************************************************

	for sym := range symbolDataMap {
		if sym == "" {
			continue
		}
		wg.Add(1)
		go func(ch chan<- RelationshipPackage, sym string, stockMap map[string]*StockDataWeekly, originalDate2 string, yearsBack int) {
			defer wg.Done()

			tooFarFlag := false
			currentDate := originalDate2 // establish currentDate as the original date to go back from

			for !tooFarFlag {
				laterDate := currentDate
				olderDate, err := findPreviousWeekString(originalDate2, currentDate, yearsBack)
				if err != nil {
					if errors.Is(err, ErrLookbackLimitReached) {
						fmt.Printf("Lookback limit reached for %s\n", sym)
						tooFarFlag = true
						continue
					} else {
						log.Printf("Error finding week before %s: %v", currentDate, err)
						break
					}
				}

				laterPredictableData, ok := stockMap[sym].TimeSeriesWeekly[laterDate]
				if !ok {
					currentDate = olderDate
					continue
				}

				olderPredictableData, ok := stockMap[sym].TimeSeriesWeekly[olderDate]
				if !ok {
					currentDate = olderDate
					continue
				}
				laterClose, err := strconv.ParseFloat(laterPredictableData.Close, 64)
				if err != nil {
					log.Printf("Error parsing later close price (%s) for %s on %s: %v", olderPredictableData.Close, sym, laterDate, err)
					currentDate = olderDate
					continue
				}

				olderClose, err := strconv.ParseFloat(olderPredictableData.Close, 64)
				if err != nil {
					log.Printf("Error parsing older close price (%s) for %s on %s: %v", olderPredictableData.Close, sym, currentDate, err)
					currentDate = olderDate
					continue
				}

				var deltaPredictable float64
				if olderClose != 0 {
					deltaPredictable = (laterClose - olderClose) / olderClose
				} else {
					log.Printf("Error: olderClose is 0, unable to calculate relationships")
					currentDate = olderDate
					continue
				}

				for s := range stockMap {
					laterPredictorData, ok := stockMap[s].TimeSeriesWeekly[olderDate]
					if !ok {
						currentDate = olderDate
						continue
					}
					laterPredictorClose, err := strconv.ParseFloat(laterPredictorData.Close, 64)
					if err != nil {
						log.Printf("Error parsing later close price (%s) for %s on %s: %v", laterPredictorData.Close, s, laterDate, err)
						currentDate = olderDate
						continue
					}
					tempOldDate, err := findPreviousWeekString(originalDate2, olderDate, yearsBack)
					if err != nil {
						if errors.Is(err, ErrLookbackLimitReached) {
							tooFarFlag = true
							break
						} else {
							log.Printf("Error finding Week Before %s: %v", currentDate, err)
							break
						}
					}

					olderPredictorData, ok := stockMap[s].TimeSeriesWeekly[tempOldDate]
					if !ok {
						currentDate = olderDate
						continue
					}
					olderPredictorClose, err := strconv.ParseFloat(olderPredictorData.Close, 64)
					if err != nil {
						log.Printf("Error parsing older close price (%s) for %s on %s: %v", olderPredictorData.Close, s, tempOldDate, err)
						currentDate = olderDate
						continue
					}
					var deltaPredictorClose float64
					if olderPredictorClose != 0 {
						deltaPredictorClose = (laterPredictorClose - olderPredictorClose) / olderPredictorClose
					} else {
						log.Printf("Error: olderClose is 0, unable to calculate predictor delta")
						currentDate = olderDate
						continue
					}

					key := RelationshipKey{
						PredictorSym:   s,
						PredictableSym: sym,
					}
					data := RelationshipData{
						PredictorSym:          s,
						PredictableSym:        sym,
						PredictableCloseDelta: deltaPredictable,
						PredictorCloseDelta:   deltaPredictorClose,
					}

					newPackage := RelationshipPackage{
						Key:  key,
						Data: data,
					}
					currentDate = olderDate
					ch <- newPackage
				}

			}
		}(ch, sym, symbolDataMap, OriginalFriday, lookBackTimeYear)
	}

	wg.Wait()
	close(ch)
	collectorwg.Wait()

	return relationshipMap, nil
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

func AnalyzeStoredDataV1(data map[RelationshipKey][]RelationshipData) ([]Prediction, error) {
	var pSym1, pSym2 string
	var c float64
	var PredictionSlice []Prediction
	for key, relationshipDataSlice := range data {
		var predictorDeltas []float64
		var predictableDeltas []float64

		for _, dataPoint := range relationshipDataSlice {
			predictorDeltas = append(predictorDeltas, dataPoint.PredictorCloseDelta)
			predictableDeltas = append(predictableDeltas, dataPoint.PredictableCloseDelta)
		}

		if len(predictorDeltas) < 2 {
			fmt.Printf("Not enough data to calculate correlation for %s and %s.\n", key.PredictorSym, key.PredictableSym)
			continue
		}

		correlation := stat.Correlation(predictorDeltas, predictableDeltas, nil)

		pSym1, pSym2, c = key.PredictableSym, key.PredictorSym, correlation
		PredictionSlice = append(PredictionSlice, Prediction{PredictableSym: pSym1, PredictorSym: pSym2, Correlation: c})

	}
	return PredictionSlice, nil
}
