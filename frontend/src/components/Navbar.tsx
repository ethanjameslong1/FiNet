import { Link } from "react-router-dom";
import { Button } from "./Button";
import { SPACING } from "../constants/global";

interface NavbarProps {
  username?: string; // optional display name
}

const Navbar = ({ username }: NavbarProps) => {
  return (
    <header className="shadow-md bg-white dark:bg-gray-800">
      <nav
        className={`container mx-auto flex justify-between items-center px-${SPACING.large} py-${SPACING.medium}`}
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
          <Button to="/logout" fullWidth={false}>
            Logout
          </Button>
        </div>
      </nav>
    </header>
  );
};

export default Navbar;
