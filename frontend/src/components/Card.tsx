import React from "react";
import { Button } from "./Button";

interface CardProps {
  title: string;
  description: string;
  icon?: React.ReactNode; // optional icon
  header?: React.ReactNode; // optional custom header (overrides title + icon)
  buttonText?: string; // optional button text
  buttonTo?: string; // optional link
  className?: string; // extra wrapper classes
  hoverScale?: boolean; // optionally enable hover scaling
}

const Card = ({
  title,
  description,
  icon,
  header,
  buttonText,
  buttonTo,
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
      {/* Optional custom header */}
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

      {buttonText && buttonTo && <Button to={buttonTo}>{buttonText}</Button>}
    </div>
  );
};

export default Card;
