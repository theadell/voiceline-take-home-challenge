# voiceline-take-home-challenge

## Overview

This repository contains my solution to the take-home coding challenge. I created a simple Go API for user authentication, allowing users to log in and register either via OAuth2 providers or with email and password.

#### Stack used
- **net/http** as an http framework
- **x/oauth2** for OAuth2 client functionality
- **SQLite** as a lightweight database
- **sqlc** for compiling SQL to type-safe go code
- **golang-migrate/migrate** for database migrations
- **alexedwards/scs** for session management
- **theadell/authress** for ID token validation

#### Structure
- **cmd/api** contiains the server (api) entry point 
- **cmd/client** contiains the client entry Point 
- **migrations** contains the database schema 

#### Demo 
You can access a live demo of the api here:

- **Live API:** [https://voiceline.adelhub.com/](https://voiceline.adelhub.com/)

View the API documentation:

- **API Documentation:** [https://voiceline.adelhub.com/docs](https://voiceline.adelhub.com/docs)

## API and Client Setup Guide
### Prerequisites
- **Go** version 1.23 or higher is required.
- **C Compiler:** A `C` compiler is required to build the project. (the sqlite3 driver relies on `cgo`).

### Steps to Run the API Locally
1. Clone the repository 
    ```sh
    git clone https://github.com/theadell/voiceline-take-home-challenge.git
    cd /voiceline-take-home-challenge
    ```
2. Download required dependencies 
   ```sh 
   go mod tidy
   ```
3. Build the API:
   ```sh 
   CGO_ENABLED=1 go build -o bin/api ./cmd/api/
   ```
   or
   ```sh
   make build 
   ```
4. Start the API and run the database migrations 
   ```sh
    ./bin/api -migrate-db
   ```
   or 
   ```sh
   make run
   ```
5. Access the API Docs 
   Once the API is running, you can view the Swagger documentation at http://localhost:8080/docs and test the endpoints in the browser


### Using the client 
I created a simple Go client in cmd/client to help test the API.

You can view the available commands and options by running:
```sh 
go run ./cmd/client help
```

- **Sign In with Google**:
You can log in using Google’s device flow by running the following command from the terminal 
    ```sh 
    go run ./cmd/client login --sso
    ```
- **Register** a New User with email and password 
  ```sh 
  go run ./cmd/client register --email your-email@example.com --password your-password
  ```
- **Login** with Email and password 
  ```sh
  go run ./cmd/client login --email your-email@example.com --password your-password
  ```

### Running the API with Docker

You can also run the API as a Docker container.

1. **Build the Docker image:**

   ```bash
   docker build -t api-image .
2. **Run the container:**

   ```bash
   docker run -p 8080:8080 api-image
