import { useState } from "react";
import Navbar from "../../components/Navbar";
import Footer from "../../components/Footer";
import PortfolioChart from "./PortfolioChart";
import CPICalculator from "../simulator/CPICalculator";

const Portfolio = () => {
  const [cpiData, setCpiData] = useState({
    initialInvestment: 0,
    monthlyContribution: 0,
    years: 0,
    expectedReturn: 0,
    variance: 0,
  });

  //API fetch for portfolio data here
  return (
    <div className="flex flex-col min-h-screen bg-gray-50 dark:bg-gray-900 font-sans">
      <Navbar />

      {/* Ticker Selector */}
      <section className="container mx-auto p-6 bg-white dark:bg-gray-800 rounded-xl shadow mb-10 mg-top-6">
        <h3 className="text-lg font-bold mb-4 text-gray-900 dark:text-gray-100">
          Build Portfolio
        </h3>
        <form className="flex gap-4 items-center">
          <input
            type="text"
            placeholder="Enter ticker (AAPL)"
            className="w-full p-2 border rounded-lg dark:bg-gray-700 dark:text-white"
          />
          <button className="bg-blue-600 text-white px-4 py-2 rounded-lg">
            Add Ticker to Portfolio
          </button>
          <button className="bg-red-600 text-white px-4 py-2 rounded-lg">
            Run Portfolio Optimizer
          </button>
        </form>
      </section>

      <CPICalculator onChange={setCpiData} defaults={cpiData} />

      <PortfolioChart cpiData={cpiData} />

      <Footer />
    </div>
  );
};

export default Portfolio;
