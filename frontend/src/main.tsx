import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";

import Login from "./pages/Login.tsx";
import UserHomepage from "./pages/Home.tsx";
import Register from "./pages/Register.tsx";
import Portfolio from "./features/optimizer/Portfolio.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/home" element={<UserHomepage name="User" />} />
        <Route path="/portfolio" element={<Portfolio />} />
      </Routes>
    </Router>
  </StrictMode>
);
