Welcome

## Prerequisites

Before getting started, ensure you have the following prerequisites installed on your machine:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Setup
0. Config
   copy the .env.example to .env and change the values

1. Clone the repository:
   ```bash
   git clone https://github.com/shandler/go-expert-observabilidade.git
   ```

2. Navigate to the project directory:
    ```bash
    cd observabilidade-open-telemetry
    ```

3. Run the application
    ```
   make build-up
    ```

4. request the endpoint POST
   ```
   curl -X POST http://localhost:8080 -d '{"zipCode": "07987110"}'
   ```

5. URLS para acessar os serviços, só clicar nos links
   - [http://localhost:16686](http://localhost:16686)
   - [http://127.0.0.1:9411/zipkin/](http://127.0.0.1:9411/zipkin/)
   - [http://localhost:9090](http://localhost:9090)
