package analysis

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
	// "github.com/google/uuid"
)

// Test Values
var AlphaVantageSymbols = []string{"MSFT", "GOOG", "TSLA", "AMZN"}
var ErrLookbackLimitReached = errors.New("lookback limit reached")

const (
	lookBackTimeYear = 10
	dateFormat       = "2006-01-02"
)

type RelationshipData struct {
	predictableSym        string
	predictorSym          string
	predictableCloseDelta float64
	predictorCloseDelta   float64
}
type RelationshipKey struct {
	predictorSym   string
	predictableSym string
}
type RelationshipPackage struct {
	Key  RelationshipKey
	Data RelationshipData
}

// Creates a hashmap of relationalData, returns the map and an error (if any)
func StoreWeeklyData(d []*StockDataWeekly, startDate string, weights StockWeights) (map[RelationshipKey][]RelationshipData, error) {
	//populating map for faster access to stock data
	symbolDataMap := make(map[string]*StockDataWeekly)
	for _, data := range d {
		symbolDataMap[data.MetaData.Symbol] = data
	}
	relationshipMap := make(map[RelationshipKey][]RelationshipData)
	var mapMutex sync.Mutex

	var wg sync.WaitGroup
	ch := make(chan RelationshipPackage)

	//COLLECTOR FUNC *********************************************************************
	go func(ch <-chan RelationshipPackage) {
		for pkg := range ch {
			mapMutex.Lock()
			relationshipMap[pkg.Key] = append(relationshipMap[pkg.Key], pkg.Data)
			mapMutex.Unlock()
			fmt.Printf("Collector recieved relationship: %+v", pkg)
		}
		fmt.Println("Collector: Channel closed, no more relationships to process.")
	}(ch)
	//COLLECTOR FUNC END *****************************************************************

	//DATE SETUP **********************************************************************
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
	//DATE SETUP **********************************************************************

	for sym := range symbolDataMap {
		if sym == "" {
			continue
		}
		wg.Add(1) //for every symbol spawn a seperate thread to analysis it for predictors
		go func(ch chan<- RelationshipPackage, sym string, stockMap map[string]*StockDataWeekly, originalDate2 string, yearsBack int) {
			defer wg.Done()

			//Helper variables to be accessed during the main loop*****************************
			tooFarFlag := false
			currentDate := originalDate2 //establish currentDate as the original date to go back from
			//Helper variables to be accessed during the main loop*****************************

			for !tooFarFlag { //for every week that it can this loop will reference every other symbol and try and relate it to the main predictable symbol (sym), will be multithread as well
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
					log.Printf("Missing later data for %s on %s", sym, laterDate)
					currentDate = olderDate
					continue
				}

				olderPredictableData, ok := stockMap[sym].TimeSeriesWeekly[olderDate]
				if !ok {
					log.Printf("Missing older data for %s on %s", sym, currentDate)
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
					laterPredictorData := stockMap[s].TimeSeriesWeekly[olderDate]
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
						log.Printf("Missing older data for %s on %s", s, tempOldDate)
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
						predictorSym:   s,
						predictableSym: sym,
					}
					data := RelationshipData{
						predictorSym:          s,
						predictableSym:        sym,
						predictableCloseDelta: deltaPredictable,
						predictorCloseDelta:   deltaPredictorClose,
					}

					newPackage := RelationshipPackage{
						Key:  key,
						Data: data,
					}
					currentDate = olderDate

					ch <- newPackage
				}

			}
		}(ch, sym, symbolDataMap, OriginalFriday, 1)
	}

	wg.Wait()
	close(ch)

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

func analyzeStoredData(map[RelationshipKey][]RelationshipData, w StockWeights)
