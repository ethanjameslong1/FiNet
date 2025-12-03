package database

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	StockDataSource                 string = "root:Front2Back!@tcp(stock_prediction_db:3306)/stock_prediction_db?parseTime=true"
	SQL_INSERT_PREDICTION           string = `INSERT INTO prediction (predictable_symbol, predictor_symbol, correlation, prediction_model, user_id) VALUES (?, ?, ?, ?, ?)`
	SQL_SELECT_ALL_PREDICTIONS      string = `SELECT id, predictable_symbol, predictor_symbol, correlation, prediction_model FROM prediction`
	SQL_SELECT_PREDICTIONS_FOR_USER string = `SELECT id, predictable_symbol, predictor_symbol, correlation, prediction_model FROM prediction WHERE user_id = ?`
	SQL_SELECT_PREDICTION_MODEL     string = `SELECT id FROM prediction_models WHERE model_name = ?`
)

type Prediction struct {
	ID                string  `json:"id"`
	PredictableSymbol string  `json:"predictable_symbol"`
	PredictorSymbol   string  `json:"predictor_symbol"`
	Correlation       float64 `json:"correlation"`
	ModelID           int     `json:"prediction_model"`
}

func (s *DBService) AddPrediction(ctx context.Context, predictableSym string, predictorSym string, cor float64, model string, userID int) error {
	modelR := s.db.QueryRowContext(ctx, SQL_SELECT_PREDICTION_MODEL, model)
	var modelID int
	err := modelR.Scan(&modelID)
	if err != nil {
		return fmt.Errorf("Error finding prediction model: %w", err)
	}

	r, err := s.db.ExecContext(ctx, SQL_INSERT_PREDICTION, predictableSym, predictorSym, cor, modelID, userID)
	if err != nil {
		return fmt.Errorf("Error inserting Prediction: %w", err)
	}
	rows, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error finding Rows Affected: %w", err)
	}
	if rows != 1 {
		return fmt.Errorf("Expected rows impacted to be 1, rows impacted %d", rows)
	}
	return nil
}

func (s *DBService) GetAllPredictions(ctx context.Context) ([]Prediction, error) {
	rows, err := s.db.QueryContext(ctx, SQL_SELECT_ALL_PREDICTIONS)
	if err != nil {
		return nil, fmt.Errorf("error querying for all predictions: %w", err)
	}
	defer rows.Close()

	var predictions []Prediction
	for rows.Next() {
		var p Prediction
		if err := rows.Scan(&p.ID, &p.PredictableSymbol, &p.PredictorSymbol, &p.Correlation, &p.ModelID); err != nil {
			return nil, fmt.Errorf("error scanning prediction row: %w", err)
		}
		predictions = append(predictions, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error after iterating prediction rows: %w", err)
	}

	return predictions, nil
}

func (s *DBService) RemoveAllPredictionsForUser(ctx context.Context, userID int) ([]Prediction, error) {
	rows, err := s.db.QueryContext(ctx, SQL_SELECT_PREDICTIONS_FOR_USER, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying for all predictions: %w", err)
	}
	defer rows.Close()

	var predictions []Prediction
	for rows.Next() {
		var p Prediction
		if err := rows.Scan(&p.ID, &p.PredictableSymbol, &p.PredictorSymbol, &p.Correlation, &p.ModelID); err != nil {
			return nil, fmt.Errorf("error scanning prediction row: %w", err)
		}
		predictions = append(predictions, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error after iterating prediction rows: %w", err)
	}

	return predictions, nil
}
