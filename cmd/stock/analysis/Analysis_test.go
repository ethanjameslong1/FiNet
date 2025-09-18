package analysis_test

import (
	"bytes"
	"github.com/ethanjameslong1/GoCloudProject.git/analysis"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"testing"
)

// TestStoreWeeklyData tests the main data collection and processing logic.
// It verifies that relationships are calculated correctly and that the function
// can gracefully handle gaps in the time-series data (e.g., market holidays).
func TestStoreWeeklyDataV1(t *testing.T) {
	// --- Mock Data Setup ---
	// This simulates the data you would get from an API call.
	// Note the gap: data for "2025-06-13" is intentionally missing for TESTA
	// to ensure the error handling for data gaps is working correctly.
	mockWeeklyData := []*analysis.StockDataWeekly{
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
	relationshipMap, err := analysis.StoreWeeklyDataV1(mockWeeklyData, "2025-06-27", analysis.StockWeights{})
	if err != nil {
		t.Fatalf("StoreWeeklyData returned an unexpected error: %v", err)
	}

	// --- Assertions ---
	if relationshipMap == nil {
		t.Fatal("Expected a valid map, but got nil")
	}

	// Check the relationship where TESTA is predictable and TESTB is the predictor.
	// This is the first valid data point the function should find.
	key := analysis.RelationshipKey{
		PredictorSym:   "TESTB",
		PredictableSym: "TESTA",
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
	if math.Abs(foundData.PredictableCloseDelta-expectedPredictableDelta) > tolerance {
		t.Errorf("For key %+v, expected predictable delta %f, but got %f", key, expectedPredictableDelta, foundData.PredictableCloseDelta)
	}

	if math.Abs(foundData.PredictorCloseDelta-expectedPredictorDelta) > tolerance {
		t.Errorf("For key %+v, expected predictor delta %f, but got %f", key, expectedPredictorDelta, foundData.PredictorCloseDelta)
	}
}

func TestAnalyzeStoredDataV1(t *testing.T) {
	// --- Test Cases ---
	testCases := []struct {
		name           string
		inputData      map[analysis.RelationshipKey][]analysis.RelationshipData
		expectedOutput string // A substring we expect to find in the printed output
		expectError    bool
	}{
		{
			name: "Perfect Positive Correlation",
			inputData: map[analysis.RelationshipKey][]analysis.RelationshipData{
				{PredictorSym: "P_POS", PredictableSym: "S_POS"}: {
					{PredictorCloseDelta: 0.1, PredictableCloseDelta: 0.2},
					{PredictorCloseDelta: 0.2, PredictableCloseDelta: 0.4},
					{PredictorCloseDelta: 0.3, PredictableCloseDelta: 0.6},
				},
			},
			expectedOutput: "Correlation between P_POS (predictor) and S_POS (predictable): 1.0000",
			expectError:    false,
		},
		{
			name: "Perfect Negative Correlation",
			inputData: map[analysis.RelationshipKey][]analysis.RelationshipData{
				{PredictorSym: "P_NEG", PredictableSym: "S_NEG"}: {
					{PredictorCloseDelta: 0.1, PredictableCloseDelta: -0.1},
					{PredictorCloseDelta: 0.2, PredictableCloseDelta: -0.2},
					{PredictorCloseDelta: 0.3, PredictableCloseDelta: -0.3},
				},
			},
			expectedOutput: "Correlation between P_NEG (predictor) and S_NEG (predictable): -1.0000",
			expectError:    false,
		},
		{
			name: "Not Enough Data For Correlation",
			inputData: map[analysis.RelationshipKey][]analysis.RelationshipData{
				{PredictorSym: "P_FEW", PredictableSym: "S_FEW"}: {
					{PredictorCloseDelta: 0.1, PredictableCloseDelta: 0.2},
				},
			},
			expectedOutput: "Not enough data to calculate correlation for P_FEW and S_FEW.",
			expectError:    false,
		},
		{
			name:           "Empty Input Map",
			inputData:      map[analysis.RelationshipKey][]analysis.RelationshipData{},
			expectedOutput: "--- Correlation Analysis ---", // Should just print the header and footer
			expectError:    false,
		},
	}

	// --- Test Execution ---
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Redirect stdout to capture the function's printed output
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			// This also redirects any `log` statements that print to stdout
			log.SetOutput(w)

			// Execute the function
			_, err := analysis.AnalyzeStoredDataV1(tc.inputData)

			// Restore stdout and close the pipe
			w.Close()
			os.Stdout = originalStdout
			log.SetOutput(os.Stderr) // Restore default logger output

			var buf bytes.Buffer
			io.Copy(&buf, r)
			capturedOutput := buf.String()

			// --- Assertions ---
			if (err != nil) != tc.expectError {
				t.Fatalf("Expected error state '%v', but got: %v", tc.expectError, err)
			}

			if !strings.Contains(capturedOutput, tc.expectedOutput) {
				t.Errorf("Expected output to contain:\n'%s'\n\nGot:\n'%s'", tc.expectedOutput, capturedOutput)
			}
		})
	}
}
