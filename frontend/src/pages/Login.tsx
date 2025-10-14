import { Link } from "react-router-dom";
import { Button } from "../components/Button";

const LoginPage = () => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 font-sans">
      <div className="bg-white dark:bg-gray-800 p-8 rounded-2xl shadow-lg w-full max-w-md">
        {/* Header */}
        <h3 className="text-2xl font-bold text-center mb-6 text-gray-900 dark:text-gray-100">
          Sign In
        </h3>
        <p className="text-center text-sm text-gray-500 dark:text-gray-400 mb-8">
          Welcome back! Please enter your credentials.
        </p>

        {/* Login Form */}
        <form method="POST" action="/login" className="space-y-5">
          {/* Username Input */}
          <div>
            <label
              htmlFor="username"
              className="block text-sm font-medium text-gray-700 dark:text-gray-300"
            >
              Username
            </label>
            <input
              type="text"
              id="username"
              name="username"
              required
              className="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 
                         rounded-lg shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100
                         focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              placeholder="Enter your username"
            />
          </div>

          {/* Password Input */}
          <div>
            <label
              htmlFor="password"
              className="block text-sm font-medium text-gray-700 dark:text-gray-300"
            >
              Password
            </label>
            <input
              type="password"
              id="password"
              name="password"
              required
              className="mt-1 block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 
                         rounded-lg shadow-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100
                         focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              placeholder="Enter your password"
            />
          </div>

          {/* Submit Button */}
          <Button type="submit">Sign In</Button>
        </form>

        {/* Register Link */}
        <p className="mt-6 text-center text-sm text-gray-600 dark:text-gray-400">
          Don&apos;t have an account?{" "}
          <Link
            to="/register"
            className="text-blue-600 hover:underline dark:text-blue-400"
          >
            Register
          </Link>
        </p>
      </div>
    </div>
  );
};

export default LoginPage;
