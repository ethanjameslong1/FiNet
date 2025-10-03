import Navbar from "../../components/Navbar";
import Footer from "../../components/Footer";
import { Line } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  ArcElement,
  Title,
  Tooltip,
  Legend
);

const Portfolio = () => {
  const cpiData = {
    labels: ["2025", "2026", "2027", "2028", "2029"],
    datasets: [
      {
        label: "CPI",
        data: [100, 103, 106, 110, 114],
        borderColor: "rgb(37, 99, 235)", // blue-600
        backgroundColor: "rgba(37, 99, 235, 0.3)",
        tension: 0.2,
      },
    ],
  };

  //   const weightsData = {
  //     labels: ["AAPL", "MSFT", "GOOGL", "TSLA"],
  //     datasets: [
  //       {
  //         data: [40, 30, 20, 10],
  //         backgroundColor: [
  //           "rgba(37, 99, 235, 0.7)",
  //           "rgba(16, 185, 129, 0.7)",
  //           "rgba(251, 191, 36, 0.7)",
  //           "rgba(239, 68, 68, 0.7)",
  //         ],
  //       },
  //     ],
  //   };

  return (
    <div className="flex flex-col min-h-screen bg-gray-50 dark:bg-gray-900 font-sans">
      <Navbar />

      {/* Sticky KPIs */}
      <section className="sticky top-0 z-10 bg-gray-50 dark:bg-gray-700 shadow-md p-4 flex gap-6 justify-center">
        <div className="text-center">
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Portfolio Value
          </p>
          <p className="text-xl font-bold text-gray-900 dark:text-gray-100">
            $250,000
          </p>
        </div>
        <div className="text-center">
          <p className="text-sm text-gray-500 dark:text-gray-400">CAGR</p>
          <p className="text-xl font-bold text-gray-900 dark:text-gray-100">
            8.4%
          </p>
        </div>
      </section>

      {/* Charts */}
      <section className="container mx-auto p-6 space-y-10">
        <div className="bg-white dark:bg-gray-800 p-6 rounded-xl shadow">
          <h3 className="text-lg font-bold mb-4 text-gray-900 dark:text-gray-100">
            CPI Over Time
          </h3>
          <Line data={cpiData} />
        </div>

        <div className="bg-white dark:bg-gray-800 p-6 rounded-xl shadow">
          <h3 className="text-lg font-bold mb-4 text-gray-900 dark:text-gray-100">
            Portfolio Weights
          </h3>
        </div>
      </section>

      {/* Ticker Selector */}
      <section className="container mx-auto p-6 bg-white dark:bg-gray-800 rounded-xl shadow mb-10">
        <h3 className="text-lg font-bold mb-4 text-gray-900 dark:text-gray-100">
          Build Portfolio
        </h3>
        <form className="space-y-4">
          <input
            type="text"
            placeholder="Enter ticker (AAPL)"
            className="w-full p-2 border rounded-lg dark:bg-gray-700 dark:text-white"
          />
          <input
            type="number"
            placeholder="Enter weight (%)"
            className="w-full p-2 border rounded-lg dark:bg-gray-700 dark:text-white"
          />
          <button className="bg-blue-600 text-white px-4 py-2 rounded-lg">
            Add Ticker
          </button>
        </form>
      </section>

      {/* CPI Assumptions */}
      <section className="container mx-auto p-6 bg-white dark:bg-gray-800 rounded-xl shadow mb-10">
        <h3 className="text-lg font-bold mb-4 text-gray-900 dark:text-gray-100">
          CPI Assumptions
        </h3>
        <form className="space-y-4">
          <input
            type="number"
            placeholder="Expected Return (%)"
            className="w-full p-2 border rounded-lg dark:bg-gray-700 dark:text-white"
          />
          <input
            type="number"
            placeholder="Variance (%)"
            className="w-full p-2 border rounded-lg dark:bg-gray-700 dark:text-white"
          />
          <input
            type="number"
            placeholder="Time Horizon (years)"
            className="w-full p-2 border rounded-lg dark:bg-gray-700 dark:text-white"
          />
        </form>
      </section>

      <Footer />
    </div>
  );
};

export default Portfolio;
