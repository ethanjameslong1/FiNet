import { useState } from "react";

interface CPICalculatorProps {
  onChange: (data: {
    initialInvestment: number;
    monthlyContribution: number;
    years: number;
    expectedReturn: number;
    variance: number;
  }) => void;
  defaults?: {
    initialInvestment: number;
    monthlyContribution: number;
    years: number;
    expectedReturn: number;
    variance: number;
  };
}

const CPICalculator = ({ onChange, defaults }: CPICalculatorProps) => {
  const [inputs, setInputs] = useState({
    initialInvestment: defaults?.initialInvestment || 0,
    monthlyContribution: defaults?.monthlyContribution || 0,
    years: defaults?.years || 0,
    expectedReturn: defaults?.expectedReturn || 0,
    variance: defaults?.variance || 0,
  });

  const handleChange = (field: string, value: number) => {
    const updated = { ...inputs, [field]: value };
    setInputs(updated);
    onChange(updated);
  };

  return (
    <div>
      {/* CPI Assumptions */}
      <section className="container mx-auto p-6 bg-white dark:bg-gray-800 rounded-xl shadow mb-10">
        <h3 className="text-lg font-bold mb-4 text-gray-900 dark:text-gray-100">
          Enter your assumptions for the CPI calculation.
        </h3>
        <form className="space-y-4">
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Initial Investment ($)
          </p>
          <input
            type="number"
            placeholder="Initial Investment ($)"
            value={inputs.initialInvestment}
            onChange={(e) => {
              handleChange(
                "initialInvestment",
                parseFloat(e.target.value) || 0
              );
            }}
            className="w-full p-2 border rounded-lg dark:bg-gray-700 dark:text-white"
          />
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Monthly Contribution ($)
          </p>
          <input
            type="number"
            placeholder="Monthly Contribution ($)"
            value={inputs.monthlyContribution}
            onChange={(e) => {
              handleChange(
                "monthlyContribution",
                parseFloat(e.target.value) || 0
              );
            }}
            className="w-full p-2 border rounded-lg dark:bg-gray-700 dark:text-white"
          />
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Time Horizon (years)
          </p>
          <input
            type="number"
            placeholder="Time Horizon (years)"
            value={inputs.years}
            onChange={(e) => {
              handleChange("years", parseFloat(e.target.value) || 0);
            }}
            className="w-full p-2 border rounded-lg dark:bg-gray-700 dark:text-white"
          />
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Estimated Annual Return (%)
          </p>
          <input
            type="number"
            placeholder="Estimated Annual Return (%)"
            value={inputs.expectedReturn}
            onChange={(e) => {
              handleChange("expectedReturn", parseFloat(e.target.value) || 0);
            }}
            className="w-full p-2 border rounded-lg dark:bg-gray-700 dark:text-white"
          />
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Estimated Rate Variance (%)
          </p>
          <input
            type="number"
            placeholder="Estimated rate variance (%)"
            value={inputs.variance}
            onChange={(e) => {
              handleChange("variance", parseFloat(e.target.value) || 0);
            }}
            className="w-full p-2 border rounded-lg dark:bg-gray-700 dark:text-white"
          />
        </form>
      </section>
    </div>
  );
};

export default CPICalculator;
