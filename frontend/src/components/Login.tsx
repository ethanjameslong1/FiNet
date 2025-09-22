import { useState } from "react";
import axios from "axios";

export default function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault(); // prevent page reload

    try {
      const response = await axios.post("/api/login", { username, password });
      console.log("Login success:", response.data);
      // TODO: redirect user or store token/session
    } catch (err) {
      console.error(err);
      setError("Invalid username or password");
    }
  };

  // async function fetchStockData(symbols: string[]) {
  //   try {
  //     const res = await fetch("http://localhost:8080/api/stock", {
  //       method: "POST",
  //       headers: { "Content-Type": "application/json" },
  //       body: JSON.stringify({ symbols, interval: "weekly" }),
  //     });

  //     if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);
  //     const data = await res.json();
  //     console.log("Stock data:", data);
  //     return data;
  //   } catch (err) {
  //     console.error("Fetch error:", err);
  //   }
  // }

  // // Example usage
  // fetchStockData(["AAPL", "GOOG", "MSFT"]);

  return (
    <div className="bg-gray-100 flex items-center justify-center min-h-screen font-sans">
      <div className="bg-white p-8 rounded-lg shadow-md w-full max-w-sm">
        <form onSubmit={handleSubmit} className="space-y-4">
          <h3 className="text-2xl font-bold text-center mb-6 text-purple-600">
            Sign In to Search File Base
          </h3>

          {error && <p className="text-red-500 text-sm text-center">{error}</p>}

          <div>
            <label
              htmlFor="username"
              className="block text-sm font-medium text-green-600"
            >
              User
            </label>
            <input
              type="text"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
              placeholder="Enter your username"
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-purple-500 focus:border-purple-500 sm:text-sm"
            />
          </div>

          <div>
            <label
              htmlFor="password"
              className="block text-sm font-medium text-green-600"
            >
              Password
            </label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              placeholder="Enter your password"
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-purple-500 focus:border-purple-500 sm:text-sm"
            />
          </div>

          <button
            type="submit"
            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-purple-600 hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 transition duration-150 ease-in-out"
          >
            Sign In
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-gray-600">
          Don't have an account?{" "}
          <a
            href="/register"
            className="font-medium text-purple-600 hover:text-purple-500"
          >
            Register User
          </a>
        </p>
      </div>
    </div>
  );
}
