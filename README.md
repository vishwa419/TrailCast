# Hiking Weather Forecast

This application provides weather data for popular hiking parks. It allows users to view weather forecasts for different parks and trails, as well as weather data based on custom coordinates. The system fetches weather data from an API and presents it in a user-friendly interface.

## Features
- Displays a list of popular hiking parks with their weather forecasts.
- Shows detailed weather data for each park, including temperature and precipitation for the next few days.
- Provides recommended hiking days based on the weather.
- Allows users to get weather data based on custom coordinates.
- Displays hiking trails data for the parks.
- Provides API endpoints for weather and trail data.

## Technologies Used
- Go (Golang)
- PostgreSQL for database management
- HTML templates for rendering dynamic pages
- External weather API (for weather data)
- External trails API (for trails data)

## Installation

To get started with the application, follow these steps:

### 1. Install PostgreSQL

Follow the instructions below based on your operating system to install and configure PostgreSQL:

#### On Linux:
1. Install PostgreSQL via your package manager:

    ```bash
    sudo apt update
    sudo apt install postgresql postgresql-contrib
    ```

2. Start the PostgreSQL service:

    ```bash
    sudo systemctl start postgresql
    ```

3. Log in to PostgreSQL as the `postgres` user:

    ```bash
    sudo -u postgres psql
    ```

4. Create the `parksdb` database and set the password for the `postgres` user:

    ```sql
    CREATE DATABASE parksdb;
    ALTER USER postgres WITH PASSWORD 'yourpassword';
    \q
    ```

5. Ensure PostgreSQL is running and accepting connections:

    ```bash
    sudo systemctl enable postgresql
    sudo systemctl start postgresql
    ```

### 2. Clone the Repository

Clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/hike_weather.git
cd hike_weather
```

### 3. Set Up Go Environment

Ensure that you have Go installed on your system. If not, you can install Go from [here](https://golang.org/dl/).

- Set the `GOPATH` and `GOROOT` (if needed).
- Install dependencies by running:

  ```bash
  go mod tidy
  ```

### 4. Update Database Connection String

In the Go source code, update the database connection string in the `main.go` (or wherever it is located) to include your PostgreSQL credentials. For example:

```go
connStr := "user=postgres password=yourpassword dbname=parksdb sslmode=disable"
```

Replace `yourpassword` with the password you set for the `postgres` user.

### 5. Run the Application

Start the Go application:

```bash
go run main.go
```

By default, the server will run on `http://localhost:8080`.

## API Endpoints

- **GET `/weather?park=<park_name>`**: Retrieves the weather forecast for the specified park.
- **GET `/weather?lat=<latitude>&lng=<longitude>`**: Retrieves weather data based on latitude and longitude.
- **GET `/trails?park=<park_name>`**: Retrieves trail data for a specific park.
- **GET `/trails?lat=<latitude>&lng=<longitude>`**: Retrieves trail data based on latitude and longitude.

## Project Structure

- `main.go`: The entry point for the application.
- `controllers.go`: Contains the controllers for handling routes and database interactions.
- `models.go`: Contains the data models used in the application.
- `templates/`: Directory for HTML templates.
- `static/`: Directory for static assets like CSS and JavaScript files.
- `go.mod`: Go modules file for managing dependencies.

## Troubleshooting

1. **Error: Database connection failed**:
   - Ensure that PostgreSQL is running (`sudo systemctl start postgresql`).
   - Verify the database name, user, and password are correct in your connection string.

2. **Error: Missing external API**:
   - If the weather and trails API are required, ensure that you have API keys set up and the corresponding code properly configured.

3. **Error: Port already in use**:
   - If port `8080` is in use, you can change the port by modifying the Go application’s `ListenAndServe` function:

     ```go
     http.ListenAndServe(":8081", nil)  // Change 8081 to another available port
     ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
