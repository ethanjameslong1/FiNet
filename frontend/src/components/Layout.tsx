import React from "react";
import Navbar from "./Navbar";
import Footer from "./Footer";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout = ({ children }: LayoutProps) => {
  return (
    <div className="min-h-screen flex flex-col font-sans bg-gray-50 dark:bg-gray-900">
      <Navbar />

      <main className={`flex-grow container mx-auto px-6 py-8`}>
        {children}
      </main>
      <Footer />
    </div>
  );
};

export default Layout;
