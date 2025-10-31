package analysis

import (
	"gonum.org/v1/gonum/stat"
	"math"
)

// CovarianceMatrixSample computes sample covariance (N-1)
// consider refactoring double map
func CovarianceMatrixSample(returns map[string][]float64) map[string]map[string]float64 {
	covMatrix := make(map[string]map[string]float64)
	tickers := make([]string, 0, len(returns))
	for t := range returns {
		tickers = append(tickers, t)
	}

	for i, a := range tickers {
		if covMatrix[a] == nil {
			covMatrix[a] = make(map[string]float64)
		}
		for j, b := range tickers {
			if i > j {
				continue
			}
			if a == b {
				n := float64(len(returns[a]))
				covMatrix[a][b] = stat.Variance(returns[a], nil) * n / (n - 1)
			} else {
				n := float64(len(returns[a]))
				cov := stat.Covariance(returns[a], returns[b], nil) * n / (n - 1)
				covMatrix[a][b] = cov
				if covMatrix[b] == nil {
					covMatrix[b] = make(map[string]float64)
				}
				covMatrix[b][a] = cov
			}
		}
	}
	return covMatrix
}

// CorrelationMatrixSample computes correlation using sample covariance
func CorrelationMatrixSample(returns map[string][]float64) map[string]map[string]float64 {
	corrMatrix := make(map[string]map[string]float64)
	covMatrix := CovarianceMatrixSample(returns)

	// Precompute sample variances
	variances := make(map[string]float64)
	for t, r := range returns {
		n := float64(len(r))
		variances[t] = stat.Variance(r, nil) * n / (n - 1)
	}

	for a, covRow := range covMatrix {
		if corrMatrix[a] == nil {
			corrMatrix[a] = make(map[string]float64)
		}
		for b, cov := range covRow {
			if a == b {
				corrMatrix[a][b] = 1.0
				continue
			}
			va, vb := variances[a], variances[b]
			if va == 0 || vb == 0 {
				corrMatrix[a][b] = 0.0
			} else {
				corrMatrix[a][b] = cov / (math.Sqrt(va) * math.Sqrt(vb))
			}
		}
	}
	return corrMatrix
}