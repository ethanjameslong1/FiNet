import { Link } from "react-router-dom";
import Logout from "../pages/LogoutButton";

interface NavbarProps {
  username?: string; // optional display name
}

const Navbar = ({ username }: NavbarProps) => {
  return (
    <header className="shadow-md bg-white dark:bg-gray-800">
      <nav
        className={`container mx-auto flex justify-between items-center px-4 py-4`}
      >
        <Link
          to="/home"
          className="text-2xl font-bold text-gray-900 dark:text-gray-100"
        >
          Finet
        </Link>

        <div className="flex items-center gap-4">
          {username && (
            <span className="text-gray-700 dark:text-gray-300">
              Hello, {username}
            </span>
          )}
          <Logout />
        </div>
      </nav>
    </header>
  );
};

export default Navbar;
