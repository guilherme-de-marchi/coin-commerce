# Coin Commerce - Golang Study Project

Welcome to Coin Commerce, a fictional personal project developed in Golang for study purposes. This application simulates a system for buying and selling fictional cryptoassets, featuring a distributed architecture with microservices, a custom load balancer, and communication through RabbitMQ using gRPC-encoded messages.

## Project Structure

The project is organized as follows:

- **RESTful API**: Responsible for receiving and routing requests to the appropriate microservices.
- **Users Microservice**: Handles user-related operations, including registration, authentication, and profile information.
- **Orders Microservice**: Manages operations related to buying and selling cryptoassets.
- **Load Balancer**: Distributes request loads among microservices for efficient performance.
- **RabbitMQ**: Facilitates communication between components through gRPC-encoded messages defined by .proto files.
- **PostgreSQL Database**: Stores essential data for the system's operation.

## Prerequisites

The only requirement for local execution is to have Golang installed in your environment.

## Configuration and Docker Execution

1. Clone the repository to your local environment.
   ```bash
   git clone https://github.com/guilherme-de-marchi/coin-commerce.git
   cd coin-commerce
   ```

2. Compile the program.
   ```bash
   go build -o coin-commerce
   ```

3. Run the Docker containers.
   ```bash
   docker compose up
   ```

## Usage

The REST API provides endpoints for user and order operations. Refer to the API documentation for details on available endpoints and how to use them.

## Contribution

If you are interested in contributing to the development of this study project, feel free to open issues, submit pull requests, or contact the development team.

## License

This project is licensed under the MIT License.
