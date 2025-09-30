import { Link } from "react-router-dom";

const RegisterPage = () => {
  return (
    <div className="bg-gray-100 flex items-center justify-center min-h-screen font-sans">
      <div className="bg-white p-8 rounded-lg shadow-md w-full max-w-sm">
        {/* Registration Form */}
        <form method="POST" action="/register" className="space-y-4">
          <h3 className="text-2xl font-bold text-center mb-6 text-purple-700">
            Register for File Base
          </h3>

          {/* Username Input */}
          <div>
            <label
              htmlFor="username"
              className="block text-sm font-medium text-green-600"
            >
              Username
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
            Register
          </button>
        </form>

        {/* Link to Login Page */}
        <p className="mt-6 text-center text-sm text-gray-600">
          Already have an account?{" "}
          <Link to="/" className="text-purple-600 hover:underline">
            Sign In
          </Link>
        </p>
      </div>
    </div>
  );
};

export default RegisterPage;
