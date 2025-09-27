package analysis

//monte carlo portfolio optimization for now, implement QP solver later
import (
	"errors"
	"math"
	"math/rand"
	"time"
)

type Portfolio struct {
	Weights map[string]float64
	Return  float64
	Risk    float64
	Sharpe  float64
}

func newRand() *rand.Rand {
	// Seed with current time for randomness
	seed := time.Now().UnixNano()
	src := rand.NewSource(seed)
	return rand.New(src)
}

//still need to add sector constraints, tag stocks by sector first, sum weights by sector <= sector max

// generateWeightConstraints returns weights that satisfy minWeight <= w[i] <= maxWeight and sum(w)=1.
// It is robust: starts with mins, splits remaining using random exponential draws, then caps+redistributes surplus.
func generateWeightConstraints(tickers []string, minWeight, maxWeight float64, randGen *rand.Rand, maxAttempts int) (map[string]float64, error) {
	n := len(tickers)
	if n == 0 {
		return nil, errors.New("no tickers provided")
	}
	if minWeight < 0 || maxWeight <= 0 || minWeight > maxWeight {
		return nil, errors.New("invalid min/max weights")
	}
	// basic feasibility
	if float64(n)*minWeight > 1.0+1e-12 {
		return nil, errors.New("sum of minimum weights exceeds 1.0; lower minWeight or reduce number of assets")
	}

	// Try attempts
	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Start with minimum weights
		weights := make(map[string]float64, n)
		sumMin := 0.0
		for _, t := range tickers {
			weights[t] = minWeight
			sumMin += minWeight
		}
		remaining := 1.0 - sumMin
		if remaining < -1e-12 {
			return nil, errors.New("minWeight infeasible")
		}
		if remaining <= 1e-12 {
			// Already satisfied (min weights sum to 1)
			return weights, nil
		}

		// Random proportions (Exp draws avoid zeros)
		props := make([]float64, n)
		sumProps := 0.0
		for i := 0; i < n; i++ {
			v := randGen.ExpFloat64()
			props[i] = v
			sumProps += v
		}
		for i, t := range tickers {
			weights[t] += remaining * (props[i] / sumProps)
		}

		// Iteratively cap at max and redistribute surplus
		for iter := 0; iter < 20; iter++ {
			surplus := 0.0
			// find any over-cap and cap it
			for _, t := range tickers {
				if weights[t] > maxWeight {
					surplus += weights[t] - maxWeight
					weights[t] = maxWeight
				}
			}
			if surplus <= 1e-12 {
				break // done
			}
			// compute total room available under maxWeight
			room := 0.0
			for _, t := range tickers {
				if weights[t] < maxWeight-1e-12 {
					room += maxWeight - weights[t]
				}
			}
			if room <= 1e-12 {
				// cannot place surplus -> break and resample
				break
			}
			// distribute surplus proportionally to available room
			for _, t := range tickers {
				if weights[t] < maxWeight-1e-12 {
					available := maxWeight - weights[t]
					weights[t] += surplus * (available / room)
				}
			}
		}

		// Final normalization to correct tiny numeric error
		sumW := 0.0
		for _, t := range tickers {
			sumW += weights[t]
		}
		if sumW == 0 {
			continue
		}
		for _, t := range tickers {
			weights[t] = weights[t] / sumW
		}

		// Validate constraints
		valid := true
		for _, t := range tickers {
			if weights[t] < minWeight-1e-9 || weights[t] > maxWeight+1e-9 {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}
		// success
		return weights, nil
	}
	return nil, errors.New("failed to generate constrained weights after maxAttempts")
}


// consider goroutine parallelism later.
func OptimizePortfolio(returns map[string][]float64, numPortfolios int, riskFreeRate float64, minWeight float64, maxWeight float64) ([]Portfolio, Portfolio) {
    var bestPortfolio Portfolio
    bestPortfolio.Sharpe = math.Inf(-1)

    // Slice to store all generated portfolios
    portfolios := make([]Portfolio, 0, numPortfolios)
	randGen := newRand()

    // Precompute expected returns and covariance matrix once
    expectedReturns := ExpectedReturn(returns)
    covMatrix := CovarianceMatrixSample(returns)

	
    tickers := make([]string, 0, len(returns))
    for t := range returns {
        tickers = append(tickers, t)
    }
	n := len(tickers)
	if n == 0 {
		return portfolios, bestPortfolio
	}

	// If basket is small, adjust maxWeight to 1/n if that is lower than the supplied maxWeight.
	if n < 10 {
		oneOverN := 1.0 / float64(n)
		if maxWeight > oneOverN {
			maxWeight = oneOverN
		}
	}

	// Safety: if minWeight * n > 1, make minWeight smaller
	if float64(n)*minWeight > 1.0 {
		minWeight = 0.0
	}

		for i := 0; i < numPortfolios; i++ {
		// generate constrained weights
		weights, err := generateWeightConstraints(tickers, minWeight, maxWeight, randGen, 200)
		if err != nil {
			// if generation fails, fallback to equal weights (but in practice this shouldn't happen)
			eq := 1.0 / float64(n)
			weights = make(map[string]float64, n)
			for _, t := range tickers {
				weights[t] = math.Min(math.Max(eq, minWeight), maxWeight)
			}
			// normalize
			sum := 0.0
			for _, v := range weights {
				sum += v
			}
			for k := range weights {
				weights[k] = weights[k] / sum
			}
		}

        // Calculate portfolio return and risk
        var portReturn, portVariance float64
        for a, wa := range weights {
            portReturn += wa * expectedReturns[a]
            for b, wb := range weights {
                portVariance += wa * wb * covMatrix[a][b]
            }
        }
        portRisk := math.Sqrt(portVariance)

        // Calculate Sharpe ratio
        sharpe := 0.0
        if portRisk != 0 {
            sharpe = (portReturn - riskFreeRate) / portRisk
        }

        portfolio := Portfolio{
            Weights: weights,
            Return:  portReturn,
            Risk:    portRisk,
            Sharpe:  sharpe,
        }

        portfolios = append(portfolios, portfolio)

        if sharpe > bestPortfolio.Sharpe {
            bestPortfolio = portfolio
        }
    }

    return portfolios, bestPortfolio
}
