import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

const Logout = () => {
  const navigate = useNavigate();

  useEffect(() => {
    localStorage.removeItem("authToken");
    localStorage.removeItem("username");

    setTimeout(() => {
      navigate("/");
    }, 100);
  }, [navigate]);

  return (
    <div className="flex items-center justify-center h-screen">
      <p>Logging out safely...</p>
    </div>
  );
};

export default Logout;
