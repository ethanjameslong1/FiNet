CREATE TABLE IF NOT EXISTS prediction_models (
    id INT AUTO_INCREMENT PRIMARY KEY,
    model_name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- CREATE TABLE IF NOT EXISTS parameters (
--     id INT AUTO_INCREMENT PRIMARY KEY,
--     parameter_name VARCHAR(100) UNIQUE NOT NULL,
--     description TEXT
-- );

CREATE TABLE IF NOT EXISTS prediction (
    id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
    predictable_symbol VARCHAR(10) NOT NULL,
    predictor_symbol VARCHAR(10) NOT NULL,
    correlation DECIMAL(10,7) NOT NULL,
    prediction_model INT NOT NULL,
    FOREIGN KEY (prediction_model) REFERENCES prediction_models(id) ON DELETE CASCADE
);

-- CREATE TABLE IF NOT EXISTS predictable_stocks (
--     id VARCHAR(36) PRIMARY KEY DEFAULT (UUID()),
--     stock_symbol VARCHAR(10) NOT NULL,
--     predicted_slope DECIMAL(8, 6) NOT NULL,
--     prediction_model_id INT,
--     FOREIGN KEY (prediction_model_id) REFERENCES prediction_models(id)
-- );

-- CREATE TABLE IF NOT EXISTS prediction_trigger_links (
--     PRIMARY KEY (predictable_stock_id, predictor_stock_event_id),
--     predictable_stock_id VARCHAR(36) NOT NULL,
--     predictor_stock_event_id VARCHAR(36) NOT NULL,
--     FOREIGN KEY (predictable_stock_id) REFERENCES predictable_stocks(id) ON DELETE CASCADE,
--     FOREIGN KEY (predictor_stock_event_id) REFERENCES predictor_stocks(id) ON DELETE CASCADE
-- );



-- INSERT INTO parameters (parameter_name, description) VALUES
-- ('Open', 'The price at which a stock starts trading at the beginning of a trading day.'), -- 1
-- ('High', 'The highest price at which a stock traded during a trading day.'), -- 2
-- ('Low', 'The lowest price at which a stock traded during a trading day.'),-- 3
-- ('Close', 'The final price at which a stock trades at the end of a trading day.'),-- 4
-- ('AdjClose', 'The closing price of the stock after accounting for corporate actions like dividends or splits. It is often used for historical analysis as it reflects the true value of a stock.'),-- 5
-- ('Volume', 'The total number of shares of a stock that were traded during a trading day.'),-- 6
-- ('DivAmount', 'The amount of dividend paid per share of a stock.');-- 7

INSERT INTO prediction_models (model_name, description) VALUES ("First Draft", "Will only use stock close data, comparing stocks 1 to 1 for simplicity. Used for establishing infrastructure for the rest of the app")
