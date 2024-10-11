# Bitcoin Pricing WebSocket Platform  

## Description
A Go-based platform that provides real-time Bitcoin pricing information via a WebSocket channel. This platform allows clients to connect and receive live updates on Bitcoin prices.

## Features
- Real-time Bitcoin price updates
- WebSocket channel for efficient data transmission
- Easy to integrate with other applications

## Installation
1. **Clone the repository**:
    ```sh
    git clone https://github.com/andreleao1/pricing-service-go.git
    ```
2. **Navigate to the project directory**:
    ```sh
    cd princing-service-go
    ```
3. **Install dependencies**:
    ```sh
    go mod tidy
    ```
4. **Run the application**:
    ```sh
    go run main.go
    ```

## Usage
1. **Start the WebSocket server**:
    ```sh
    go run main.go
    ```
2. **Connect to the WebSocket server**:
    - Use a WebSocket client to connect to `ws://localhost:8080/ws`
    - You will start receiving real-time Bitcoin price updates

## API Endpoints
- **WebSocket Endpoint**: `ws://localhost:8080/ws`
    - Connect to this endpoint to receive real-time Bitcoin price updates.

## Contributing
1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -m 'Add some feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Open a pull request.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.