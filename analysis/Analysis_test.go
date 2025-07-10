package analysis

import (
	"math"
	"testing"
)

// TestStoreWeeklyData tests the main data collection and processing logic.
// It verifies that relationships are calculated correctly and that the function
// can gracefully handle gaps in the time-series data (e.g., market holidays).
func TestStoreWeeklyData(t *testing.T) {
	// --- Mock Data Setup ---
	// This simulates the data you would get from an API call.
	// Note the gap: data for "2025-06-13" is intentionally missing for TESTA
	// to ensure the error handling for data gaps is working correctly.
	mockWeeklyData := []*StockDataWeekly{
		{
			MetaData: struct {
				Information   string `json:"1. Information"`
				Symbol        string `json:"2. Symbol"`
				LastRefreshed string `json:"3. Last Refreshed"`
				TimeZone      string `json:"4. Time Zone"`
			}{Symbol: "TESTA"},
			TimeSeriesWeekly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-06-27": {Close: "100.00"},
				"2025-06-20": {Close: "90.00"},
				// "2025-06-13" is missing for TESTA
				"2025-06-06": {Close: "85.00"},
			},
		},
		{
			MetaData: struct {
				Information   string `json:"1. Information"`
				Symbol        string `json:"2. Symbol"`
				LastRefreshed string `json:"3. Last Refreshed"`
				TimeZone      string `json:"4. Time Zone"`
			}{Symbol: "TESTB"},
			TimeSeriesWeekly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-06-27": {Close: "500.00"},
				"2025-06-20": {Close: "525.00"},
				"2025-06-13": {Close: "510.00"},
				"2025-06-06": {Close: "505.00"},
			},
		},
	}

	// --- Test Execution ---
	relationshipMap, err := StoreWeeklyData(mockWeeklyData, "2025-06-27", StockWeights{})
	if err != nil {
		t.Fatalf("StoreWeeklyData returned an unexpected error: %v", err)
	}

	// --- Assertions ---
	if relationshipMap == nil {
		t.Fatal("Expected a valid map, but got nil")
	}

	// Check the relationship where TESTA is predictable and TESTB is the predictor.
	// This is the first valid data point the function should find.
	key := RelationshipKey{
		predictorSym:   "TESTB",
		predictableSym: "TESTA",
	}

	dataSlice, ok := relationshipMap[key]
	if !ok {
		t.Fatalf("Expected to find key %+v in the relationship map, but it was missing.", key)
	}

	if len(dataSlice) == 0 {
		t.Fatalf("Data slice for key %+v is empty.", key)
	}

	// The first valid relationship found should be this one.
	foundData := dataSlice[0]

	// Define the expected deltas for this specific relationship.
	const tolerance = 1e-9

	// Expected Predictable (TESTA) Delta: (100.00 - 90.00) / 90.00
	// This is for the week of 2025-06-27 vs 2025-06-20.
	expectedPredictableDelta := (100.00 - 90.00) / 90.00

	// Expected Predictor (TESTB) Delta: (525.00 - 510.00) / 510.00
	// This is for the week of 2025-06-20 vs 2025-06-13.
	expectedPredictorDelta := (525.00 - 510.00) / 510.00

	// Assert that the calculated deltas match the expected values.
	if math.Abs(foundData.predictableCloseDelta-expectedPredictableDelta) > tolerance {
		t.Errorf("For key %+v, expected predictable delta %f, but got %f", key, expectedPredictableDelta, foundData.predictableCloseDelta)
	}

	if math.Abs(foundData.predictorCloseDelta-expectedPredictorDelta) > tolerance {
		t.Errorf("For key %+v, expected predictor delta %f, but got %f", key, expectedPredictorDelta, foundData.predictorCloseDelta)
	}
}
