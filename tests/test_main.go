package main

import (
	"fmt"

	"github.com/ethanjameslong1/GoCloudProject.git/analysis"
)

// MockStockData returns mock data for 5 stocks over 9 months
func MockStockData() []*analysis.StockDataMonthly {
	return []*analysis.StockDataMonthly{
		{
			MetaData: struct {
				Information   string "json:\"1. Information\""
				Symbol        string "json:\"2. Symbol\""
				LastRefreshed string "json:\"3. Last Refreshed\""
				TimeZone      string "json:\"4. Time Zone\""
			}{
				Information:   "Monthly Adjusted Prices and Volumes",
				Symbol:        "IBM",
				LastRefreshed: "2025-09-25",
				TimeZone:      "US/Eastern",
			},
			TimeSeriesMonthly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-09-25": {Open: "240.90", High: "284.23", Low: "238.25", Close: "281.44", AdjClose: "281.44", Volume: "89264876", DivAmount: "0.00"},
				"2025-08-29": {Open: "251.40", High: "255.00", Low: "233.36", Close: "243.49", AdjClose: "243.49", Volume: "104957357", DivAmount: "1.68"},
				"2025-07-31": {Open: "294.55", High: "295.61", Low: "252.22", Close: "253.15", AdjClose: "251.41", Volume: "109055173", DivAmount: "0.00"},
				"2025-06-30": {Open: "257.85", High: "296.16", Low: "257.22", Close: "294.78", AdjClose: "292.75", Volume: "74395935", DivAmount: "0.00"},
				"2025-05-30": {Open: "241.44", High: "269.28", Low: "237.95", Close: "259.06", AdjClose: "257.28", Volume: "78164014", DivAmount: "1.68"},
			},
		},
		{
			MetaData: struct {
				Information   string "json:\"1. Information\""
				Symbol        string "json:\"2. Symbol\""
				LastRefreshed string "json:\"3. Last Refreshed\""
				TimeZone      string "json:\"4. Time Zone\""
			}{
				Information:   "Monthly Adjusted Prices and Volumes",
				Symbol:        "AAPL",
				LastRefreshed: "2025-09-25",
				TimeZone:      "US/Eastern",
			},
			TimeSeriesMonthly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-09-25": {Open: "168.00", High: "171.00", Low: "166.50", Close: "170.50", AdjClose: "170.50", Volume: "95000000", DivAmount: "0.22"},
				"2025-08-29": {Open: "165.50", High: "167.80", Low: "163.20", Close: "165.80", AdjClose: "165.80", Volume: "89000000", DivAmount: "0.22"},
				"2025-07-31": {Open: "159.00", High: "161.50", Low: "158.00", Close: "160.40", AdjClose: "160.40", Volume: "88000000", DivAmount: "0.20"},
				"2025-06-30": {Open: "171.00", High: "173.00", Low: "170.00", Close: "172.10", AdjClose: "172.10", Volume: "97000000", DivAmount: "0.24"},
				"2025-05-30": {Open: "167.50", High: "169.50", Low: "166.50", Close: "168.00", AdjClose: "168.00", Volume: "91000000", DivAmount: "0.22"},
			},
		},
		{
			MetaData: struct {
				Information   string "json:\"1. Information\""
				Symbol        string "json:\"2. Symbol\""
				LastRefreshed string "json:\"3. Last Refreshed\""
				TimeZone      string "json:\"4. Time Zone\""
			}{
				Information:   "Monthly Adjusted Prices and Volumes",
				Symbol:        "GOOG",
				LastRefreshed: "2025-09-25",
				TimeZone:      "US/Eastern",
			},
			TimeSeriesMonthly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-09-25": {Open: "1320.00", High: "1340.00", Low: "1310.00", Close: "1335.00", AdjClose: "1335.00", Volume: "1200000", DivAmount: "0.00"},
				"2025-08-29": {Open: "1290.00", High: "1305.00", Low: "1280.00", Close: "1295.00", AdjClose: "1295.00", Volume: "1100000", DivAmount: "0.00"},
				"2025-07-31": {Open: "1270.00", High: "1285.00", Low: "1260.00", Close: "1275.00", AdjClose: "1275.00", Volume: "1150000", DivAmount: "0.00"},
				"2025-06-30": {Open: "1300.00", High: "1325.00", Low: "1295.00", Close: "1315.00", AdjClose: "1315.00", Volume: "1250000", DivAmount: "0.00"},
				"2025-05-30": {Open: "1285.00", High: "1300.00", Low: "1275.00", Close: "1290.00", AdjClose: "1290.00", Volume: "1180000", DivAmount: "0.00"},
			},
		},
		{
			MetaData: struct {
				Information   string "json:\"1. Information\""
				Symbol        string "json:\"2. Symbol\""
				LastRefreshed string "json:\"3. Last Refreshed\""
				TimeZone      string "json:\"4. Time Zone\""
			}{
				Information:   "Monthly Adjusted Prices and Volumes",
				Symbol:        "MSFT",
				LastRefreshed: "2025-09-25",
				TimeZone:      "US/Eastern",
			},
			TimeSeriesMonthly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-09-25": {Open: "210.00", High: "215.00", Low: "208.00", Close: "212.50", AdjClose: "212.50", Volume: "78000000", DivAmount: "0.56"},
				"2025-08-29": {Open: "208.00", High: "212.00", Low: "205.00", Close: "209.00", AdjClose: "209.00", Volume: "76000000", DivAmount: "0.56"},
				"2025-07-31": {Open: "205.00", High: "210.00", Low: "203.00", Close: "207.50", AdjClose: "207.50", Volume: "74000000", DivAmount: "0.55"},
				"2025-06-30": {Open: "212.00", High: "216.00", Low: "210.00", Close: "214.00", AdjClose: "214.00", Volume: "77000000", DivAmount: "0.57"},
				"2025-05-30": {Open: "208.00", High: "212.00", Low: "205.00", Close: "209.50", AdjClose: "209.50", Volume: "75000000", DivAmount: "0.56"},
			},
		},
		{
			MetaData: struct {
				Information   string "json:\"1. Information\""
				Symbol        string "json:\"2. Symbol\""
				LastRefreshed string "json:\"3. Last Refreshed\""
				TimeZone      string "json:\"4. Time Zone\""
			}{
				Information:   "Monthly Adjusted Prices and Volumes",
				Symbol:        "AMZN",
				LastRefreshed: "2025-09-25",
				TimeZone:      "US/Eastern",
			},
			TimeSeriesMonthly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-09-25": {Open: "3200.00", High: "3300.00", Low: "3180.00", Close: "3250.00", AdjClose: "3250.00", Volume: "4500000", DivAmount: "0.00"},
				"2025-08-29": {Open: "3150.00", High: "3220.00", Low: "3120.00", Close: "3180.00", AdjClose: "3180.00", Volume: "4600000", DivAmount: "0.00"},
				"2025-07-31": {Open: "3100.00", High: "3170.00", Low: "3080.00", Close: "3125.00", AdjClose: "3125.00", Volume: "4700000", DivAmount: "0.00"},
				"2025-06-30": {Open: "3250.00", High: "3350.00", Low: "3220.00", Close: "3300.00", AdjClose: "3300.00", Volume: "4800000", DivAmount: "0.00"},
				"2025-05-30": {Open: "3180.00", High: "3250.00", Low: "3150.00", Close: "3200.00", AdjClose: "3200.00", Volume: "4550000", DivAmount: "0.00"},
			},
		},		// NFLX
		{
			MetaData: struct {
				Information   string "json:\"1. Information\""
				Symbol        string "json:\"2. Symbol\""
				LastRefreshed string "json:\"3. Last Refreshed\""
				TimeZone      string "json:\"4. Time Zone\""
			}{
				Information:   "Monthly Adjusted Prices and Volumes",
				Symbol:        "NFLX",
				LastRefreshed: "2025-09-25",
				TimeZone:      "US/Eastern",
			},
			TimeSeriesMonthly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-05-30": {Open: "450.00", High: "470.00", Low: "440.00", Close: "460.00", AdjClose: "460.00", Volume: "7500000", DivAmount: "0.00"},
				"2025-04-30": {Open: "440.00", High: "460.00", Low: "430.00", Close: "450.00", AdjClose: "450.00", Volume: "7400000", DivAmount: "0.00"},
				"2025-03-31": {Open: "430.00", High: "450.00", Low: "420.00", Close: "440.00", AdjClose: "440.00", Volume: "7300000", DivAmount: "0.00"},
				"2025-02-28": {Open: "420.00", High: "440.00", Low: "410.00", Close: "430.00", AdjClose: "430.00", Volume: "7200000", DivAmount: "0.00"},
				"2025-01-31": {Open: "410.00", High: "430.00", Low: "400.00", Close: "420.00", AdjClose: "420.00", Volume: "7100000", DivAmount: "0.00"},
			},
		},
		// TSLA
		{
			MetaData: struct {
				Information   string "json:\"1. Information\""
				Symbol        string "json:\"2. Symbol\""
				LastRefreshed string "json:\"3. Last Refreshed\""
				TimeZone      string "json:\"4. Time Zone\""
			}{
				Information:   "Monthly Adjusted Prices and Volumes",
				Symbol:        "TSLA",
				LastRefreshed: "2025-09-25",
				TimeZone:      "US/Eastern",
			},
			TimeSeriesMonthly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-05-30": {Open: "700.00", High: "750.00", Low: "680.00", Close: "730.00", AdjClose: "730.00", Volume: "30000000", DivAmount: "0.00"},
				"2025-04-30": {Open: "680.00", High: "720.00", Low: "660.00", Close: "700.00", AdjClose: "700.00", Volume: "29000000", DivAmount: "0.00"},
				"2025-03-31": {Open: "650.00", High: "690.00", Low: "630.00", Close: "680.00", AdjClose: "680.00", Volume: "28000000", DivAmount: "0.00"},
				"2025-02-28": {Open: "620.00", High: "660.00", Low: "600.00", Close: "650.00", AdjClose: "650.00", Volume: "27000000", DivAmount: "0.00"},
				"2025-01-31": {Open: "600.00", High: "640.00", Low: "580.00", Close: "620.00", AdjClose: "620.00", Volume: "26000000", DivAmount: "0.00"},
			},
		},
		// NVDA
		{
			MetaData: struct {
				Information   string "json:\"1. Information\""
				Symbol        string "json:\"2. Symbol\""
				LastRefreshed string "json:\"3. Last Refreshed\""
				TimeZone      string "json:\"4. Time Zone\""
			}{
				Information:   "Monthly Adjusted Prices and Volumes",
				Symbol:        "NVDA",
				LastRefreshed: "2025-09-25",
				TimeZone:      "US/Eastern",
			},
			TimeSeriesMonthly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-05-30": {Open: "500.00", High: "520.00", Low: "480.00", Close: "510.00", AdjClose: "510.00", Volume: "12000000", DivAmount: "0.00"},
				"2025-04-30": {Open: "480.00", High: "500.00", Low: "460.00", Close: "490.00", AdjClose: "490.00", Volume: "11500000", DivAmount: "0.00"},
				"2025-03-31": {Open: "460.00", High: "480.00", Low: "440.00", Close: "470.00", AdjClose: "470.00", Volume: "11000000", DivAmount: "0.00"},
				"2025-02-28": {Open: "440.00", High: "460.00", Low: "420.00", Close: "450.00", AdjClose: "450.00", Volume: "10500000", DivAmount: "0.00"},
				"2025-01-31": {Open: "420.00", High: "440.00", Low: "400.00", Close: "430.00", AdjClose: "430.00", Volume: "10000000", DivAmount: "0.00"},
			},
		},
		// BABA
		{
			MetaData: struct {
				Information   string "json:\"1. Information\""
				Symbol        string "json:\"2. Symbol\""
				LastRefreshed string "json:\"3. Last Refreshed\""
				TimeZone      string "json:\"4. Time Zone\""
			}{
				Information:   "Monthly Adjusted Prices and Volumes",
				Symbol:        "BABA",
				LastRefreshed: "2025-09-25",
				TimeZone:      "US/Eastern",
			},
			TimeSeriesMonthly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-05-30": {Open: "180.00", High: "190.00", Low: "175.00", Close: "185.00", AdjClose: "185.00", Volume: "8000000", DivAmount: "0.00"},
				"2025-04-30": {Open: "175.00", High: "185.00", Low: "170.00", Close: "180.00", AdjClose: "180.00", Volume: "7800000", DivAmount: "0.00"},
				"2025-03-31": {Open: "170.00", High: "180.00", Low: "165.00", Close: "175.00", AdjClose: "175.00", Volume: "7600000", DivAmount: "0.00"},
				"2025-02-28": {Open: "165.00", High: "175.00", Low: "160.00", Close: "170.00", AdjClose: "170.00", Volume: "7400000", DivAmount: "0.00"},
				"2025-01-31": {Open: "160.00", High: "170.00", Low: "155.00", Close: "165.00", AdjClose: "165.00", Volume: "7200000", DivAmount: "0.00"},
			},
		},
		// JPM
		{
			MetaData: struct {
				Information   string "json:\"1. Information\""
				Symbol        string "json:\"2. Symbol\""
				LastRefreshed string "json:\"3. Last Refreshed\""
				TimeZone      string "json:\"4. Time Zone\""
			}{
				Information:   "Monthly Adjusted Prices and Volumes",
				Symbol:        "JPM",
				LastRefreshed: "2025-09-25",
				TimeZone:      "US/Eastern",
			},
			TimeSeriesMonthly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-05-30": {Open: "145.00", High: "150.00", Low: "140.00", Close: "148.00", AdjClose: "148.00", Volume: "7000000", DivAmount: "0.90"},
				"2025-04-30": {Open: "140.00", High: "145.00", Low: "135.00", Close: "142.00", AdjClose: "142.00", Volume: "6800000", DivAmount: "0.88"},
				"2025-03-31": {Open: "135.00", High: "140.00", Low: "130.00", Close: "138.00", AdjClose: "138.00", Volume: "6600000", DivAmount: "0.85"},
				"2025-02-28": {Open: "130.00", High: "135.00", Low: "125.00", Close: "132.00", AdjClose: "132.00", Volume: "6400000", DivAmount: "0.83"},
				"2025-01-31": {Open: "125.00", High: "130.00", Low: "120.00", Close: "128.00", AdjClose: "128.00", Volume: "6200000", DivAmount: "0.80"},
			},
		},		// WMT
		{
			MetaData: struct {
				Information   string "json:\"1. Information\""
				Symbol        string "json:\"2. Symbol\""
				LastRefreshed string "json:\"3. Last Refreshed\""
				TimeZone      string "json:\"4. Time Zone\""
			}{
				Information:   "Monthly Adjusted Prices and Volumes",
				Symbol:        "WMT",
				LastRefreshed: "2025-09-25",
				TimeZone:      "US/Eastern",
			},
			TimeSeriesMonthly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-05-30": {Open: "160.00", High: "165.00", Low: "155.00", Close: "162.50", AdjClose: "162.50", Volume: "9000000", DivAmount: "0.56"},
				"2025-04-30": {Open: "155.00", High: "160.00", Low: "150.00", Close: "157.50", AdjClose: "157.50", Volume: "8800000", DivAmount: "0.55"},
				"2025-03-31": {Open: "150.00", High: "155.00", Low: "145.00", Close: "152.50", AdjClose: "152.50", Volume: "8600000", DivAmount: "0.55"},
				"2025-02-28": {Open: "145.00", High: "150.00", Low: "140.00", Close: "147.50", AdjClose: "147.50", Volume: "8400000", DivAmount: "0.54"},
				"2025-01-31": {Open: "140.00", High: "145.00", Low: "135.00", Close: "142.50", AdjClose: "142.50", Volume: "8200000", DivAmount: "0.53"},
			},
		},
		// DIS
		{
			MetaData: struct {
				Information   string "json:\"1. Information\""
				Symbol        string "json:\"2. Symbol\""
				LastRefreshed string "json:\"3. Last Refreshed\""
				TimeZone      string "json:\"4. Time Zone\""
			}{
				Information:   "Monthly Adjusted Prices and Volumes",
				Symbol:        "DIS",
				LastRefreshed: "2025-09-25",
				TimeZone:      "US/Eastern",
			},
			TimeSeriesMonthly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-05-30": {Open: "130.00", High: "135.00", Low: "125.00", Close: "132.50", AdjClose: "132.50", Volume: "7000000", DivAmount: "0.88"},
				"2025-04-30": {Open: "125.00", High: "130.00", Low: "120.00", Close: "127.50", AdjClose: "127.50", Volume: "6800000", DivAmount: "0.87"},
				"2025-03-31": {Open: "120.00", High: "125.00", Low: "115.00", Close: "122.50", AdjClose: "122.50", Volume: "6600000", DivAmount: "0.85"},
				"2025-02-28": {Open: "115.00", High: "120.00", Low: "110.00", Close: "117.50", AdjClose: "117.50", Volume: "6400000", DivAmount: "0.84"},
				"2025-01-31": {Open: "110.00", High: "115.00", Low: "105.00", Close: "112.50", AdjClose: "112.50", Volume: "6200000", DivAmount: "0.83"},
			},
		},
		// INTC
		{
			MetaData: struct {
				Information   string "json:\"1. Information\""
				Symbol        string "json:\"2. Symbol\""
				LastRefreshed string "json:\"3. Last Refreshed\""
				TimeZone      string "json:\"4. Time Zone\""
			}{
				Information:   "Monthly Adjusted Prices and Volumes",
				Symbol:        "INTC",
				LastRefreshed: "2025-09-25",
				TimeZone:      "US/Eastern",
			},
			TimeSeriesMonthly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-05-30": {Open: "55.00", High: "58.00", Low: "53.00", Close: "56.50", AdjClose: "56.50", Volume: "12000000", DivAmount: "0.38"},
				"2025-04-30": {Open: "53.00", High: "56.00", Low: "51.00", Close: "54.50", AdjClose: "54.50", Volume: "11800000", DivAmount: "0.37"},
				"2025-03-31": {Open: "51.00", High: "54.00", Low: "49.00", Close: "52.50", AdjClose: "52.50", Volume: "11600000", DivAmount: "0.36"},
				"2025-02-28": {Open: "49.00", High: "52.00", Low: "47.00", Close: "50.50", AdjClose: "50.50", Volume: "11400000", DivAmount: "0.35"},
				"2025-01-31": {Open: "47.00", High: "50.00", Low: "45.00", Close: "48.50", AdjClose: "48.50", Volume: "11200000", DivAmount: "0.34"},
			},
		},
		// PFE
		{
			MetaData: struct {
				Information   string "json:\"1. Information\""
				Symbol        string "json:\"2. Symbol\""
				LastRefreshed string "json:\"3. Last Refreshed\""
				TimeZone      string "json:\"4. Time Zone\""
			}{
				Information:   "Monthly Adjusted Prices and Volumes",
				Symbol:        "PFE",
				LastRefreshed: "2025-09-25",
				TimeZone:      "US/Eastern",
			},
			TimeSeriesMonthly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-05-30": {Open: "45.00", High: "47.00", Low: "43.50", Close: "46.00", AdjClose: "46.00", Volume: "10000000", DivAmount: "0.40"},
				"2025-04-30": {Open: "44.00", High: "46.00", Low: "42.50", Close: "45.00", AdjClose: "45.00", Volume: "9800000", DivAmount: "0.39"},
				"2025-03-31": {Open: "43.50", High: "45.50", Low: "42.00", Close: "44.00", AdjClose: "44.00", Volume: "9600000", DivAmount: "0.38"},
				"2025-02-28": {Open: "42.50", High: "44.50", Low: "41.00", Close: "43.00", AdjClose: "43.00", Volume: "9400000", DivAmount: "0.37"},
				"2025-01-31": {Open: "41.50", High: "43.50", Low: "40.00", Close: "42.00", AdjClose: "42.00", Volume: "9200000", DivAmount: "0.36"},
			},
		},
		// SPY
		{
			MetaData: struct {
				Information   string "json:\"1. Information\""
				Symbol        string "json:\"2. Symbol\""
				LastRefreshed string "json:\"3. Last Refreshed\""
				TimeZone      string "json:\"4. Time Zone\""
			}{
				Information:   "Monthly Adjusted Prices and Volumes",
				Symbol:        "SPY",
				LastRefreshed: "2025-09-25",
				TimeZone:      "US/Eastern",
			},
			TimeSeriesMonthly: map[string]struct {
				Open      string `json:"1. open"`
				High      string `json:"2. high"`
				Low       string `json:"3. low"`
				Close     string `json:"4. close"`
				AdjClose  string `json:"5. adjusted close"`
				Volume    string `json:"6. volume"`
				DivAmount string `json:"7. dividend amount"`
			}{
				"2025-05-30": {Open: "430.00", High: "440.00", Low: "425.00", Close: "435.00", AdjClose: "435.00", Volume: "150000000", DivAmount: "1.55"},
				"2025-04-30": {Open: "425.00", High: "435.00", Low: "420.00", Close: "430.00", AdjClose: "430.00", Volume: "148000000", DivAmount: "1.50"},
				"2025-03-31": {Open: "420.00", High: "430.00", Low: "415.00", Close: "425.00", AdjClose: "425.00", Volume: "146000000", DivAmount: "1.48"},
				"2025-02-28": {Open: "415.00", High: "425.00", Low: "410.00", Close: "420.00", AdjClose: "420.00", Volume: "144000000", DivAmount: "1.45"},
				"2025-01-31": {Open: "410.00", High: "420.00", Low: "405.00", Close: "415.00", AdjClose: "415.00", Volume: "142000000", DivAmount: "1.43"},
			},
		},
	}
}

// TestMain runs your 5-stock test workflow
func TestMain() {
	stocks := MockStockData()

	// Step 2: Extract adjusted close prices
	fmt.Println("=== Step 2: Extract Adjusted Close Prices ===")
	adjClose := analysis.ExtractMonthlyAdjClosePrices(stocks)
	for symbol, prices := range adjClose {
		fmt.Printf("%s: %v\n", symbol, prices)
	}
	fmt.Println()

	// Step 3: Compute monthly returns
	fmt.Println("=== Step 3: Compute Monthly Returns ===")
	returns := analysis.MonthlyStockReturns(adjClose)
	for symbol, r := range returns {
		fmt.Printf("%s: %v\n", symbol, r)
	}
	fmt.Println()

	// Step 4: Expected return and standard deviation
	fmt.Println("=== Step 4: Expected Returns & Standard Deviations ===")
	expected := analysis.ExpectedReturn(returns)
	stddev := analysis.StandardDeviation(returns)
	for symbol := range expected {
		fmt.Printf("%s: Expected=%.4f, StdDev=%.4f\n", symbol, expected[symbol], stddev[symbol])
	}
	fmt.Println()

	// Step 5: Covariance and Correlation
	fmt.Println("=== Step 5: Covariance & Correlation Matrices ===")
	cov := analysis.CovarianceMatrixSample(returns)
	corr := analysis.CorrelationMatrixSample(returns)
	fmt.Printf("Covariance: %+v\n", cov)
	fmt.Printf("Correlation: %+v\n\n", corr)

	// Step 6: Monte Carlo Portfolio Optimization
	fmt.Println("=== Step 6: Monte Carlo Portfolio Optimization ===")
	portfolios, best := analysis.OptimizePortfolio(returns, 100, 0.0, 0.02, 0.20)
	fmt.Printf("Best Portfolio (Sharpe): %+v\n", best)
	fmt.Printf("Total Portfolios Generated: %d\n", len(portfolios))
	fmt.Println()

	// Step 7: Beta Coefficients
	fmt.Println("=== Step 7: Beta Coefficients ===")
	// For demonstration, using SPY returns as market returns (not realistic)
	marketReturns := returns["SPY"]
	betas := analysis.BetaCoefficients(returns, marketReturns)
	for symbol, beta := range betas {
		fmt.Printf("%s: Beta=%.4f\n", symbol, beta)
	}
	fmt.Println()
}

func main() {
	fmt.Println("=== Running 20-Stock Test Workflow ===")
	TestMain()
}
