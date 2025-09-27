// StockTest.jsx
import React, { useState } from "react";

export default function StockTest() {
  const [symbol, setSymbol] = useState("");

  const fetchStockData = async () => {
    if (!symbol) return;

    try {
      const response = await fetch("http://localhost:8080/api/stocks", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ symbol: symbol.toUpperCase() }),
      });

      if (!response.ok) {
        const text = await response.text();
        console.error("API Error:", text);
        return;
      }

      const data = await response.json();
      console.log("Stock Data:", data);
    } catch (err) {
      console.error("Fetch Error:", err);
    }
  };

  return (
    <div>
      <h1>Stock API Test</h1>
      <input
        type="text"
        value={symbol}
        onChange={(e) => setSymbol(e.target.value)}
        placeholder="Enter symbol (e.g., AAPL)"
      />
      <button onClick={fetchStockData}>Fetch Stock Data</button>
    </div>
  );
}
