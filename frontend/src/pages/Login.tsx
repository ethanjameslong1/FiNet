import { Link } from "react-router-dom";
import { useState, useEffect } from "react";

export function TestApi() {
  const [msg, setMsg] = useState("");

  useEffect(() => {
    fetch(`${import.meta.env.VITE_API_URL}/test`)
      .then((res) => res.json())
      .then((data) => setMsg(data.message))
      .catch((err) => console.error(err));
  }, []);

  return <div>Backend Says: {msg} </div>;
}

const LoginPage = () => {
  return (
    <div className="bg-gray-100 flex items-center justify-center min-h-screen font-sans">
      <div className="bg-white p-8 rounded-lg shadow-md w-full max-w-sm">
        <div>
          <TestApi />
        </div>
        <form method="POST" action="/login" className="space-y-4">
          <h3 className="text-2xl font-bold text-center mb-6 text-purple-700">
            Sign In to Search File Base
          </h3>

          {/* Username Input */}
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
              name="username"
              required
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm 
                         focus:outline-none focus:ring-purple-500 focus:border-purple-500 sm:text-sm"
              placeholder="Enter your username"
            />
          </div>

          {/* Password Input */}
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
              name="password"
              required
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm 
                         focus:outline-none focus:ring-purple-500 focus:border-purple-500 sm:text-sm"
              placeholder="Enter your password"
            />
          </div>

          {/* Submit Button */}
          <button
            type="submit"
            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm 
                       text-sm font-medium text-white bg-purple-600 hover:bg-purple-700 
                       focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 
                       transition duration-150 ease-in-out"
          >
            Sign In
          </button>
        </form>

        {/* Register User Link */}
        <p className="mt-6 text-center text-sm text-gray-600">
          Don&apos;t have an account?{" "}
          <Link to="/register" className="text-purple-600 hover:underline">
            Register User
          </Link>
        </p>
      </div>
    </div>
  );
};

export default LoginPage;
