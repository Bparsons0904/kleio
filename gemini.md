
# Gemini Code Understanding

This document provides an overview of the Kleio project, a self-hosted vinyl record collection manager.

## About the Project

Kleio is a web application that allows users to manage their vinyl record collection. It integrates with the Discogs API to fetch collection data and provides features to track record plays, cleaning history, and stylus wear. The application is composed of a Go backend and a SolidJS frontend.

## Tech Stack

-   **Backend**: Go
-   **Frontend**: SolidJS, TypeScript, SCSS
-   **Database**: SQLite
-   **Containerization**: Docker

## Project Structure

The project is organized into the following main directories:

-   `clio/`: Contains the SolidJS frontend application.
    -   `src/`: The source code for the frontend.
        -   `components/`: Reusable UI components.
        -   `pages/`: The main pages of the application.
        -   `provider/`: Application-wide state management.
        -   `utils/`: Utility functions, including API communication.
-   `cmd/`: Contains the main entry point for the backend application.
    -   `api/`: The main application package.
-   `internal/`: Contains the core logic of the backend application.
    -   `controller/`: Handles the business logic of the application.
    -   `database/`: Manages database interactions.
    -   `server/`: Defines the HTTP server and API routes.
-   `assets/`: Contains static assets like images.
-   `kleio_data/`: Default directory for the SQLite database.

## Getting Started

To run the project, you can use Docker (recommended) or build from source.

### Docker

1.  Make sure you have Docker installed.
2.  Run the following command to start the application:
    ```bash
    docker-compose up -d
    ```
3.  The application will be available at `http://localhost:38080`.

### Manual Build

1.  **Backend**:
    ```bash
    go mod download
    go build -o kleio cmd/api/main.go
    ./kleio
    ```
2.  **Frontend**:
    ```bash
    cd clio
    npm install
    npm run build
    ```

## Key Files

### Backend

-   `cmd/api/main.go`: The main entry point for the backend application.
-   `internal/server/server.go`: Initializes the server and database.
-   `internal/server/routes.go`: Defines all the API routes.
-   `internal/database/database.go`: Contains the database connection logic.

### Frontend

-   `clio/src/App.tsx`: The main application component.
-   `clio/src/index.tsx`: The entry point for the frontend application.
-   `clio/vite.config.ts`: The build configuration for the frontend.
-   `clio/package.json`: Defines the frontend dependencies and scripts.
