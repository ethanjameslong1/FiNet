const Logout = ({ buttonText = "Logout" }) => {
  return (
    <div>
      <button
        onClick={() => {
          window.location.href = "/finet/logout";
        }}
        className="inline-block px-6 py-3 text-base font-medium text-white bg-blue-600 rounded-lg shadow-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-all"
      >
        {buttonText}
      </button>
    </div>
  );
};

export default Logout;
