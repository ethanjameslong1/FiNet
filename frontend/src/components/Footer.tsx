const Footer = () => {
  return (
    <footer className="bg-white dark:bg-gray-800 shadow-inner mt-12">
      <div className="container mx-auto px-6 py-6 text-center text-gray-600 dark:text-gray-400 text-sm">
        © {new Date().getFullYear()} Finet. All rights reserved.
      </div>
    </footer>
  );
};

export default Footer;
