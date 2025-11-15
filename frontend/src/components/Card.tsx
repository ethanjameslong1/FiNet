import React from "react";
import { Link } from "react-router-dom";

interface CardProps {
  title: string;
  description: string;
  icon?: React.ReactNode;
  header?: React.ReactNode;
  buttonText?: string;
  buttonTo?: string;
  onClick?: () => void;
  className?: string;
  hoverScale?: boolean;
}

const Card = ({
  title,
  description,
  icon,
  header,
  buttonText,
  buttonTo,
  onClick,
  className = "",
  hoverScale = false,
}: CardProps) => {
  return (
    <div
      className={`bg-white dark:bg-gray-800 p-8 rounded-2xl shadow-lg
                 hover:shadow-xl transition-shadow duration-300
                 ${hoverScale ? "hover:scale-105 transition-transform" : ""}
                 ${className}`}
    >
      {header ? (
        <div className="mb-6">{header}</div>
      ) : (
        <>
          {icon && (
            <div className="flex items-center justify-center h-16 w-16 bg-blue-100 rounded-full mb-6">
              {icon}
            </div>
          )}
          <h3 className="text-2xl font-semibold text-gray-900 dark:text-gray-100 mb-3">
            {title}
          </h3>
        </>
      )}

      <p className="text-gray-600 dark:text-gray-300 mb-6">{description}</p>

      {buttonText &&
        (buttonTo ? (
          <Link
            to={buttonTo}
            className="inline-block px-6 py-3 text-base font-medium text-white bg-blue-600 rounded-lg shadow-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all"
          >
            {buttonText}
          </Link>
        ) : onClick ? (
          <button
            onClick={onClick}
            className="inline-block px-6 py-3 text-base font-medium text-white bg-blue-600 rounded-lg shadow-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all"
          >
            {buttonText}
          </button>
        ) : null)}
    </div>
  );
};

export default Card;
