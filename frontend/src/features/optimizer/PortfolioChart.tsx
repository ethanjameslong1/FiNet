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
import { Line } from "react-chartjs-2";

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

interface PortfolioChartProps {
  cpiData: {
    initialInvestment: number;
    monthlyContribution: number;
    years: number;
    expectedReturn: number;
    variance: number;
  };
}

const PortfolioChart = ({ cpiData }: PortfolioChartProps) => {
  const labels = Array.from(
    { length: cpiData.years },
    (_, i) => `${new Date().getFullYear() + i}`
  );

  const dataPoints = labels.map((_, i) => {
    const t = i + 1;

    return (
      cpiData.initialInvestment *
        Math.pow(1 + cpiData.expectedReturn / 100, t) +
      cpiData.monthlyContribution *
        ((Math.pow(1 + cpiData.expectedReturn / 100, t) - 1) /
          (cpiData.expectedReturn / 100)) *
        (1 + cpiData.expectedReturn / 100)
    );
  });

  const chartData = {
    labels,
    datasets: [
      {
        label: "CPI",
        data: dataPoints,
        borderColor: "rgba(75, 192, 192, 1)",
        backgroundColor: "rgba(75, 192, 192, 0.2)",
      },
    ],
  };

  return (
    <div>
      <div className="bg-white dark:bg-gray-800 p-6 rounded-xl shadow">
        <h3 className="text-lg font-bold mb-4 text-gray-900 dark:text-gray-100">
          Portfolio Growth Over Time
        </h3>
        <Line data={chartData} />
      </div>
    </div>
  );
};

export default PortfolioChart;
