# GoCloudProject
This project will integrate various cloud computing technologies, nearly entirely coded in Go, to help users predict Stock Market changes based on custom variables.

Users will connect to the application online and sign in to their accounts. They will be able to determine what variables matter to them and what weight to associate those variables with. The application will then start processing.
Based on the breadth (the number of Stocks to compare and consider) and the specific variables and their complexity, the user will pay for compute time. The application will scale using Docker images hosted on AWS Services, likely ECS or Fargate.

They will then be able to get notifications based on their specific "research" when live stocks display behaviour that, based on the custom analysis, usually indicates either it or another stock will change a certain way. 

This will give users the foresight to sell or buy stocks.
