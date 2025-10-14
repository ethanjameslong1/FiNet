import React from "react";
import { Link } from "react-router-dom";

interface ButtonProps {
  children: React.ReactNode;
  to?: string; // if provided, renders a Link
  onClick?: () => void; // optional click handler
  type?: "button" | "submit"; // standard button type
  fullWidth?: boolean; // optional full width
}

export const Button = ({
  children,
  to,
  onClick,
  type = "button",
  fullWidth = true,
}: ButtonProps) => {
  const baseClasses = `
      ${fullWidth ? "w-full" : "inline-flex"} 
      justify-center items-center 
      px-4 py-2 
      rounded-lg shadow-sm 
      text-white bg-blue-600 hover:bg-blue-700 
    `;

  if (to) {
    return (
      <Link to={to} className={baseClasses}>
        {children}
      </Link>
    );
  }

  return (
    <button type={type} onClick={onClick} className={baseClasses}>
      {children}
    </button>
  );
};
