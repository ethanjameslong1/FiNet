import { Link } from "react-router-dom";

interface UserHomepageProps {
  name: string; // replace {{.UserData.Name}} with a prop
}

const UserHomepage = ({ name }: UserHomepageProps) => {
  return (
    <div className="bg-gray-100 min-h-screen font-sans">
      {/* Header Navigation */}
      <header className="bg-white shadow-md">
        <nav className="container mx-auto px-6 py-4 flex justify-between items-center">
          <Link to="/home" className="text-2xl font-bold text-purple-600">
            Finet
          </Link>
          <Link
            to="/logout"
            className="bg-purple-600 text-white py-2 px-4 rounded-md hover:bg-purple-700 transition duration-300"
          >
            Logout
          </Link>
        </nav>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-6 py-12">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-800">Welcome, {name}!</h1>
          <p className="text-gray-600 mt-2">What would you like to do today?</p>
        </div>

        {/* Feature Cards */}
        <div className="grid md:grid-cols-2 gap-8 max-w-4xl mx-auto">
          {/* Raw Stock Data Card */}
          <div className="bg-white p-8 rounded-lg shadow-lg hover:shadow-xl transition-shadow duration-300">
            <div className="flex items-center justify-center h-16 w-16 bg-purple-100 rounded-full mb-6">
              {/* SVG Icon for Data */}
              <svg
                className="h-8 w-8 text-purple-600"
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
            </div>
            <h3 className="text-2xl font-semibold text-gray-800 mb-3">
              Raw Stock Data
            </h3>
            <p className="text-gray-600 mb-6">
              Access and review unprocessed stock market data directly from the
              source.
            </p>
            <a
              href="/rawdata"
              className="inline-block w-full text-center bg-gray-200 text-gray-800 font-medium py-3 px-6 rounded-md hover:bg-gray-300 transition duration-300"
            >
              Go to Raw Data
            </a>
          </div>

          {/* Stock Analysis Card */}
          <div className="bg-white p-8 rounded-lg shadow-lg hover:shadow-xl transition-shadow duration-300">
            <div className="flex items-center justify-center h-16 w-16 bg-purple-100 rounded-full mb-6">
              {/* SVG Icon for Analysis */}
              <svg
                className="h-8 w-8 text-purple-600"
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
            </div>
            <h3 className="text-2xl font-semibold text-gray-800 mb-3">
              Stock Analysis
            </h3>
            <p className="text-gray-600 mb-6">
              Utilize our tools to perform in-depth analysis and generate
              predictive insights.
            </p>
            <a
              href="/stock"
              className="inline-block w-full text-center bg-purple-600 text-white font-medium py-3 px-6 rounded-md hover:bg-purple-700 transition duration-300"
            >
              Go to Analysis
            </a>
          </div>
        </div>
      </main>
    </div>
  );
};

export default UserHomepage;
