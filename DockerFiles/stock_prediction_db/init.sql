
CREATE TABLE IF NOT EXISTS predictor_stocks (
    id UUID PRIMARY KEY DEFAULT (UUID()),
    stock_symbol VARCHAR(10) NOT NULL,
    parameter_id INT NOT NULL,
    parameter_delta DECIMAL(8,6) NOT NULL,
    FOREIGN KEY (parameter_id) REFERENCES parameters(id)
);

CREATE TABLE IF NOT EXISTS predictable_stocks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    stock_symbol VARCHAR(10) NOT NULL,
    predicted_slope DECIMAL(8, 6) NOT NULL,
    prediction_model_id INT, --if this could just somehow hold 1-3 id values that would solve my problem.
    FOREIGN KEY (prediction_model_id) REFERENCES prediction_models(id)
);

CREATE TABLE IF NOT EXISTS prediction_trigger_links (
    PRIMARY KEY (predictable_stock_id, predictor_stock_event_id),       -- Composite PK ensures uniqueness
    predictable_stock_id INT NOT NULL,                                  -- Foreign Key to `predictable_stocks.id`
    predictor_stock_event_id UUID NOT NULL,                             -- Foreign Key to `predictor_stocks.id`
    FOREIGN KEY (predictable_stock_id) REFERENCES predictable_stocks(id) ON DELETE CASCADE,
    FOREIGN KEY (predictor_stock_event_id) REFERENCES predictor_stocks(id) ON DELETE CASCADE
);

--HELPER TABLES, NOT RELAVENT TO ACTUAL ANALYSIS
CREATE TABLE IF NOT EXISTS parameters (
    id INT AUTO_INCREMENT PRIMARY KEY,
    parameter_name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT
);


INSERT INTO parameters (parameter_name, description) VALUES
('Open', 'The price at which a stock starts trading at the beginning of a trading day.'),
('High', 'The highest price at which a stock traded during a trading day.'),
('Low', 'The lowest price at which a stock traded during a trading day.'),
('Close', 'The final price at which a stock trades at the end of a trading day.'),
('AdjClose', 'The closing price of the stock after accounting for corporate actions like dividends or splits. It is often used for historical analysis as it reflects the true value of a stock.'),
('Volume', 'The total number of shares of a stock that were traded during a trading day.'),
('DivAmount', 'The amount of dividend paid per share of a stock.');

CREATE TABLE IF NOT EXISTS prediction_models (
    id INT AUTO_INCREMENT PRIMARY KEY,
    model_name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO prediction_models (model_name, description) VALUES ("First Draft", "compares only 4 stocks, checks the predicted movements first and checks backwards to detect causes, likely bad")

