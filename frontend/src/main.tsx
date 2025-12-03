import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";

import Login from "./pages/Login.tsx";
import Logout from "./pages/Logout.tsx";
import UserHomepage from "./pages/Home.tsx";
import Register from "./pages/Register.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <Router>
      <Routes>
        <Route path="/" element={<UserHomepage />} />
        <Route path="/register" element={<Register />} />
        <Route path="/home" element={<UserHomepage />} />
        {/* <Route path="/portfolio" element={<Portfolio />} /> */}
        <Route path="/login" element={<Login />} />
        <Route path="/logout" element={<Logout />} />
      </Routes>
    </Router>
  </StrictMode>
);
