import Layout from "../components/Layout";
import Card from "../components/Card";
import { useState, useEffect } from 'react'; // Combined and corrected imports
import { useNavigate } from 'react-router-dom';

const UserHomepage = () => {
    const navigate = useNavigate();
    const [name, setName] = useState<string | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true); // Added for initial loading state

    // 1. Logic for validating the session
    const validateSession = async () => {
        const authToken = localStorage.getItem("authToken");

        if (!authToken) {
            navigate("/login"); // No token found, redirect immediately
            return;
        }

        try {
            const response = await fetch("finet/Middleware", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ authToken }),
            });
            
            const data = await response.json();
            
            if (!response.ok) {
                // Token is invalid/expired according to the backend
                localStorage.setItem("authToken", "");
                localStorage.setItem("username", "");
                setError(data.error || "Session validation failed.");
                navigate("/login");
            } else { 
                // SUCCESS case: Session is valid, update local storage with fresh data (if provided)
                localStorage.setItem("authToken", data.authToken);
                localStorage.setItem("username", data.username);
                
                // Set the name for display after successful validation
                setName(data.username); 
            }
        } catch (err) {
            console.error("Request failed:", err);
            setError("A network error occurred.");
            navigate("/login");
        } finally {
            setIsLoading(false); // Finished validation/loading
        }
    };
    // 4. useEffect to load username
useEffect(() => {
        validateSession();
    }, []); // Runs once on mount

    // 3. Optional: Show loading state while validation happens
    if (isLoading) {
        return <div className="min-h-screen flex items-center justify-center">Loading...</div>;
    }
    return (
        <Layout>
            <div className="text-center mb-12">
                <h1 className="text-4xl font-bold text-gray-900 dark:text-gray-100">
                    Welcome, {name || "User"}!
                </h1>
                <p className="text-gray-600 dark:text-gray-300 mt-2">
                    What would you like to do today?
                </p>
            </div>
            {/* Display error if one exists */}
            {error && (
                <div className="text-red-500 text-center mb-4">{error}</div>
            )}

            <div className="grid md:grid-cols-2 gap-8 max-w-4xl mx-auto">
                <Card
                    title="Raw Stock Data"
                    description="Access and review unprocessed stock market data directly from the source."
                    buttonText="Go to Raw Data"
                    onClick={() => {
                        window.location.href = "/finet/rawdata";
                    }}
                    hoverScale
                />

                <Card
                    title="Stock Analysis"
                    description="Utilize our tools to perform in-depth analysis and generate predictive insights."
                    buttonText="Go to Analysis"
                    onClick={() => {
                        window.location.href = "/finet/stock";
                    }}
                    hoverScale
                />

                <Card
                    title="Portfolio Management"
                    description="Manage your investment portfolio with our intuitive tools and insights."
                    buttonText="Optimize Portfolio"
                    buttonTo="/portfolio"
                    hoverScale
                />
            </div>
        </Layout>
    );
};

export default UserHomepage;
