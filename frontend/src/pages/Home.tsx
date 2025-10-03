import Layout from "../components/Layout";
import Card from "../components/Card";

const UserHomepage = ({ name }: { name: string }) => {
  return (
    <Layout>
      <div className="text-center mb-12">
        <h1 className="text-4xl font-bold text-gray-900 dark:text-gray-100">
          Welcome, {name}!
        </h1>
        <p className="text-gray-600 dark:text-gray-300 mt-2">
          What would you like to do today?
        </p>
      </div>

      <div className="grid md:grid-cols-2 gap-8 max-w-4xl mx-auto">
        <Card
          title="Raw Stock Data"
          description="Access and review unprocessed stock market data directly from the source."
          icon={
            <svg
              className="h-8 w-8 text-blue-600"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M9 17v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"
              />
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M13 17v-2a4 4 0 00-4-4h-2m8 6h2a4 4 0 004-4v-2a4 4 0 00-4-4h-2a4 4 0 00-4 4v2"
              />
            </svg>
          }
          buttonText="Go to Raw Data"
          buttonTo="/rawdata"
          hoverScale
        />
        <Card
          title="Stock Analysis"
          description="Utilize our tools to perform in-depth analysis and generate predictive insights."
          icon={
            <svg
              className="h-8 w-8 text-blue-600"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
              />
            </svg>
          }
          buttonText="Go to Analysis"
          buttonTo="/stock"
          hoverScale
        />
        <Card
          title="Portfolio Management"
          description="Manage your investment portfolio with our intuitive tools and insights."
          icon={
            <svg
              className="h-8 w-8 text-blue-600"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
              />
            </svg>
          }
          buttonText="Optimize Portfolio"
          buttonTo="/portfolio"
          hoverScale
        />
      </div>
    </Layout>
  );
};

export default UserHomepage;
