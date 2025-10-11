# Lacpa

A full-stack web application built with Go Fiber backend and HTMX/Tailwind CSS frontend, using MongoDB for data persistence.

## Tech Stack

### Backend
- **Go** - Programming language
- **Fiber** - Fast HTTP web framework
- **MongoDB** - NoSQL database
- **MongoDB Go Driver** - Official MongoDB driver for Go

### Frontend
- **HTMX** - Modern interactions without JavaScript frameworks
- **Tailwind CSS** - Utility-first CSS framework
- **Vanilla JavaScript** - For custom interactions

## Project Structure

```
.
├── main.go                 # Application entry point
├── config/                 # Configuration files
│   └── database.go        # MongoDB connection setup
├── models/                # Data models
│   └── item.go           # Item model
├── repository/            # Data access layer
│   └── repository.go     # MongoDB repository implementation
├── handler/               # HTTP handlers
│   └── handler.go        # API endpoint handlers
└── LACPA_Web/            # Frontend files
    ├── index.html        # Main HTML page
    ├── js/
    │   └── app.js       # Custom JavaScript
    └── css/             # Custom styles (if needed)
```

## Prerequisites

- Go 1.21 or higher
- MongoDB 4.4 or higher

## Installation

1. Clone the repository:
```bash
git clone https://github.com/AliSleiman0/Lacpa.git
cd Lacpa
```

2. Install Go dependencies:
```bash
go mod download
```

3. Create a `.env` file from the example:
```bash
cp .env.example .env
```

4. Update the `.env` file with your MongoDB connection string if needed:
```env
PORT=3000
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=lacpa
```

## Running the Application

### Option 1: Using Docker Compose (Recommended)

1. Start MongoDB using Docker Compose:
```bash
docker-compose up -d
```

2. Start the server:
```bash
go run main.go
```

3. Open your browser and navigate to:
```
http://localhost:3000
```

### Option 2: Using Local MongoDB

1. Make sure MongoDB is running on your system

2. Start the server:
```bash
go run main.go
```

3. Open your browser and navigate to:
```
http://localhost:3000
```

## API Endpoints

### Health Check
- `GET /api/health` - Check if the server is running

### Items
- `POST /api/items` - Create a new item
- `GET /api/items` - Get all items
- `GET /api/items/:id` - Get a specific item
- `PUT /api/items/:id` - Update an item
- `DELETE /api/items/:id` - Delete an item

## Development

### Building for Production

```bash
go build -o lacpa
./lacpa
```

### Running Tests

```bash
go test ./...
```

## Features

- ✅ RESTful API with Go Fiber
- ✅ MongoDB integration with repository pattern
- ✅ Clean architecture with separated layers
- ✅ HTMX for dynamic frontend interactions
- ✅ Tailwind CSS for modern, responsive UI
- ✅ CRUD operations for items
- ✅ Environment-based configuration

## License

This project is licensed under the MIT License.