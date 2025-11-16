GoFiNet Stock Analysis Web App

GoFiNet is a web application built with Go that allows users to analyze financial stocks. Users can register and log in to a secure session to submit analysis requests based on specific stocks and the weight of various stock metrics. The application then processes this information to provide financial insights and predictions. The backend is designed for concurrency and scalability, featuring robust session management and a containerized architecture for easy deployment.
Key Features 🔑

    User Authentication: Secure user registration and login functionality to protect user data and analysis history.

    Session Management: Persistent user sessions are managed with a 24-hour lifetime. An automated cleanup service runs hourly in the background to remove expired sessions, ensuring database efficiency.

    Custom Stock Analysis: Users can request stock analysis by specifying stock tickers and custom weights for different financial variables, tailoring the analysis to their strategy.

    View Predictions: Authenticated users can view the results and historical predictions from their previous stock analysis requests.

    Containerized Deployment: Utilizes Docker and Docker Compose for a streamlined setup, consistent development/production environments, and scalability. The application and its databases (user sessions and stock data) are isolated in their own containers.

Tech Stack 🛠️

This project is built with a focus on performance, concurrency, and modern deployment practices.

    Go (Golang): The core backend language, chosen for its high performance, simplicity, and excellent support for concurrency.

    net/http: Go's native library is used to build the web server, define API endpoints (/login, /stock, /predictions, etc.), and handle all HTTP routing and requests.

    Go Routines (Concurrency): A concurrent goroutine is used to run a non-blocking background task that periodically cleans up expired user sessions from the database without interfering with the main application flow.

    Docker & Docker Compose: The application and its dependencies are fully containerized using Docker. Docker Compose orchestrates the multi-container setup (including the Go application, user database, and stock database), simplifying the entire development and deployment process.

    Alpha Vantage API: This external API is used to fetch the real-time and historical stock data required for financial analysis and predictions.

    SQL Database: A relational database is used to store user credentials, session data, and stock prediction information, managed via the custom database package.
