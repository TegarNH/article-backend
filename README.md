# Article Management Project

A simple full-stack application for managing articles, built with a Go (Gin) backend and a Vanilla JavaScript frontend. The backend is designed for deployment on Render, and the frontend is optimized for Vercel.

-----

## Technologies Used

  - **Backend**: Go, Gin, `go-sql-driver/mysql`, `golang-migrate`
  - **Frontend**: HTML, CSS, Vanilla JavaScript, Bootstrap 5 (via CDN)
  - **Database**: MySQL / TiDB
  - **Deployment**: Render (Backend), Vercel (Frontend)

-----

## Prerequisites

Before you begin, ensure you have the following installed on your local machine:

  - [Go](https://golang.org/dl/) (version 1.21 or newer)
  - A local MySQL server (e.g., [XAMPP](https://www.apachefriends.org/index.html))
  - [golang-migrate/migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

-----

## Backend Setup (Go)

Follow these steps to get the backend server running locally.

### 1\. Database Setup

1.  Start your XAMPP/MySQL server.
2.  Open phpMyAdmin (usually `http://localhost/phpmyadmin`).
3.  Create a new database named `article`.

### 2\. Environment Variables

1.  Create a file named `.env` in the root directory of the project.

2.  Add the following content to the `.env` file. This configures the application to connect to your local database.

    ```env
    # .env
    DB_DSN="root:@tcp(127.0.0.1:3306)/article?parseTime=true"
    APP_PORT="8080"
    GIN_MODE="debug"
    ```

### 3\. Database Migration

Run the migration to create the `posts` table in your `article` database. Execute this command from the root project directory:

```bash
migrate -database "mysql://root:@tcp(127.0.0.1:3306)/article" -path db/migrations up
```

*(Note: Ensure the path to your migrations folder is correct.)*

### 4\. Install Dependencies

Download the Go modules required for the project.

```bash
go mod tidy
```

### 5\. Run the Server

Start the backend server. The server will run on `http://localhost:8080`.

```bash
go run main.go
```

*(Assuming your main Go file is `main.go` in the root. Adjust if necessary, e.g., `go run ./api/index.go`)*

-----

## Deployment

  - **Backend**: The Go application is configured to be deployed on **Render**.
  - **Frontend**: The static frontend is configured for deployment on **Vercel**.
