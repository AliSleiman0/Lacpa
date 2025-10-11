# Quick Start Guide

This guide will help you get Lacpa up and running in minutes.

## Prerequisites

- Go 1.21 or higher installed
- Docker and Docker Compose (optional, for MongoDB)
- Git

## Step 1: Clone the Repository

```bash
git clone https://github.com/AliSleiman0/Lacpa.git
cd Lacpa
```

## Step 2: Set Up Environment

Copy the example environment file:

```bash
cp .env.example .env
```

The default configuration should work out of the box.

## Step 3: Start MongoDB

### Option A: Using Docker (Recommended)

```bash
make docker-up
# or
docker-compose up -d
```

### Option B: Local MongoDB

Make sure MongoDB is running on `localhost:27017`

## Step 4: Install Dependencies

```bash
make deps
# or
go mod download
```

## Step 5: Run the Application

```bash
make run
# or
go run main.go
```

You should see:
```
Server starting on port 3000
```

## Step 6: Access the Application

Open your browser and navigate to:

```
http://localhost:3000
```

## Using the Application

### Creating an Item

1. Fill in the "Name" field
2. Fill in the "Description" field
3. Click "Create Item"
4. The item will appear in the list below

### Viewing Items

Items are automatically loaded when you open the page. Click "Refresh" to reload them.

### Deleting an Item

Click the "Delete" button on any item. You'll be asked to confirm before deletion.

## Development Mode

For hot reload during development:

```bash
# Install Air
make install-tools

# Run with hot reload
make dev
```

## Testing the API

### Health Check

```bash
curl http://localhost:3000/api/health
```

### Create Item

```bash
curl -X POST http://localhost:3000/api/items \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Item","description":"This is a test"}'
```

### List Items

```bash
curl http://localhost:3000/api/items
```

### Get Single Item

```bash
curl http://localhost:3000/api/items/{item-id}
```

### Update Item

```bash
curl -X PUT http://localhost:3000/api/items/{item-id} \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated Item","description":"Updated description"}'
```

### Delete Item

```bash
curl -X DELETE http://localhost:3000/api/items/{item-id}
```

## Stopping the Application

1. Press `Ctrl+C` in the terminal where the application is running
2. Stop MongoDB (if using Docker):

```bash
make docker-down
# or
docker-compose down
```

## Common Issues

### Port Already in Use

If port 3000 is already in use, change it in your `.env` file:

```env
PORT=8080
```

### Cannot Connect to MongoDB

Make sure MongoDB is running:

```bash
# Check Docker containers
docker ps

# Or check if MongoDB is running locally
sudo systemctl status mongod  # Linux
brew services list | grep mongodb  # macOS
```

### Build Errors

Clean and rebuild:

```bash
make clean
make deps
make build
```

## Next Steps

- Read the [ARCHITECTURE.md](ARCHITECTURE.md) to understand the system design
- Check the [README.md](README.md) for detailed documentation
- Explore the codebase starting with `main.go`
- Add your own models and endpoints

## Available Make Commands

```bash
make run          # Run the application
make build        # Build the binary
make test         # Run tests
make clean        # Clean build artifacts
make dev          # Run with hot reload
make docker-up    # Start MongoDB
make docker-down  # Stop MongoDB
make fmt          # Format code
make vet          # Vet code
make deps         # Download dependencies
make tidy         # Tidy dependencies
make setup        # Full setup
```

## Support

For issues or questions, please open an issue on GitHub.
