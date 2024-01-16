# Spotify API

This project, developed by Ing Samuel Seda, is a code challenge from Brillio. It focuses on consuming the Spotify API, implementing OAuth 2.0 login, and performing searches by ISRC and artist, storing the results in a local MySQL database or interacting with the Spotify API.

## Getting Started

### Install Go Dependencies

```bash
go mod tidy
```

### Install Swagger Dependency

```bash
swag init --parseDependency true
```

### Install MySQL

Ensure you have MySQL installed. If not, you can download it from [MySQL Downloads](https://dev.mysql.com/downloads/).

### Set MySQL User and Password

In the `db/database.go` file, set your MySQL user and password in the `dsn` variable:

```go
dsn := "your_username:your_password@tcp(localhost:3306)/spotify?parseTime=true"
```

### Generate Spotify Client ID and Client Secret

Visit [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/applications) and create a new application to obtain your Spotify client ID and client secret.

### Set Spotify Client ID and Client Secret

In the `config/oauth.go` file, set your Client ID and Client Secret in the `clientID` & `clientSecret` variable:

```go
var (
	clientID     = "your_clientID"
	clientSecret = "your_clientSecret"
    ...
)
```

### Run the Application

```bash
go run .
```

## Endpoints and Server Information

The application exposes the following endpoints:

- **GET /:** Welcome page
- **GET /login:** Initiates the Spotify login flow
- **GET /callback:** Handles the Spotify callback after a successful login
- **GET /search/:isrc:** Searches for a track by ISRC code
- **GET /search/artist/:name:** Searches for tracks by artist name
- **GET /swagger/*any:** Swagger documentation

The server is listening on port 8080.

## Testing

1. Obtain a Spotify token by clicking "Login with Spotify" at [http://localhost:8080](http://localhost:8080).
![Captura de pantalla 2024-01-16 032728](https://github.com/samuelzyx/spotify-api/assets/12131059/b963cb93-ad29-45ef-b481-0c2c7c4e52f8)

2. Use the provided Thunder collection (`/thunder-client-collection`) to interact with the defined endpoints.
![Captura de pantalla 2024-01-16 032910](https://github.com/samuelzyx/spotify-api/assets/12131059/455bed11-0bd5-43ab-905e-474c4cd094b0)

3. Check the Swagger documentation at [http://localhost:8080/swagger/index.html#/](http://localhost:8080/swagger/index.html#/).
![Captura de pantalla 2024-01-16 032820](https://github.com/samuelzyx/spotify-api/assets/12131059/626efb85-5d24-4573-9e0d-1d032e8825de)

4. Install Workbench to login to MySql Database Viewer
Artist table
![Captura de pantalla 2024-01-16 032747](https://github.com/samuelzyx/spotify-api/assets/12131059/28b15de9-7db4-4cde-a07b-87a872e0f823)
Tracks table
![Captura de pantalla 2024-01-16 032757](https://github.com/samuelzyx/spotify-api/assets/12131059/476fc590-ad41-49f8-a8b9-531d25547380)


**Note:** The token is stored in memory and resets with each server restart, requiring a new login.
