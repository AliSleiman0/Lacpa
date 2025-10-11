# Lacpa Architecture

## System Overview

Lacpa is a full-stack web application with a clear separation of concerns, following best practices for both backend and frontend development.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         Client (Browser)                     │
│  ┌────────────────────────────────────────────────────────┐ │
│  │          LACPA_Web (Frontend)                          │ │
│  │  • index.html - Main UI                                │ │
│  │  • Tailwind CSS - Styling                              │ │
│  │  • HTMX - Dynamic interactions                         │ │
│  │  • app.js - Custom JavaScript                          │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP/REST API
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Go Fiber Backend                          │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  main.go - Application Entry Point                     │ │
│  │  • Initializes MongoDB connection                      │ │
│  │  • Sets up Fiber app with middleware                   │ │
│  │  • Configures routes                                   │ │
│  │  • Serves static files                                 │ │
│  └────────────────────────────────────────────────────────┘ │
│                              │                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  handler/ - HTTP Request Handlers                      │ │
│  │  • handler.go - API endpoint logic                     │ │
│  │  • Validates requests                                  │ │
│  │  • Calls repository methods                            │ │
│  │  • Returns HTTP responses                              │ │
│  └────────────────────────────────────────────────────────┘ │
│                              │                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  repository/ - Data Access Layer                       │ │
│  │  • repository.go - MongoDB operations                  │ │
│  │  • CRUD operations                                     │ │
│  │  • Database queries                                    │ │
│  └────────────────────────────────────────────────────────┘ │
│                              │                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  models/ - Data Models                                 │ │
│  │  • item.go - Item entity definition                    │ │
│  └────────────────────────────────────────────────────────┘ │
│                              │                               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │  config/ - Configuration                               │ │
│  │  • database.go - MongoDB connection setup              │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ MongoDB Driver
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      MongoDB Database                        │
│  • Collection: items                                         │
│  • Document structure: Item model                            │
└─────────────────────────────────────────────────────────────┘
```

## Request Flow

### Creating an Item

1. **User Action**: User fills form and clicks "Create Item"
2. **HTMX**: Intercepts form submission, sends POST to `/api/items`
3. **Handler**: `CreateItem` validates request body
4. **Repository**: `CreateItem` inserts document into MongoDB
5. **Response**: Handler returns JSON with created item
6. **HTMX**: Receives response and updates items list
7. **JavaScript**: Formats and displays the new item

### Listing Items

1. **Page Load**: HTMX triggers GET request to `/api/items`
2. **Handler**: `GetAllItems` calls repository
3. **Repository**: `GetAllItems` queries MongoDB
4. **Response**: Handler returns JSON array of items
5. **JavaScript**: Transforms JSON into HTML cards
6. **Browser**: Displays formatted items

### Deleting an Item

1. **User Action**: User clicks "Delete" button
2. **HTMX**: Confirms action, sends DELETE to `/api/items/:id`
3. **Handler**: `DeleteItem` validates ID and calls repository
4. **Repository**: `DeleteItem` removes document from MongoDB
5. **Response**: Handler returns success message
6. **HTMX**: Removes item from DOM with animation

## Technology Choices

### Backend

**Go Fiber**
- Fast and lightweight web framework
- Express-like API familiar to developers
- Built on fasthttp for excellent performance
- Great middleware ecosystem

**Repository Pattern**
- Separates data access from business logic
- Makes code more testable
- Allows easy database swapping
- Clean abstraction over MongoDB

**MongoDB**
- Flexible schema for evolving requirements
- JSON-like documents match Go structs
- Excellent Go driver support
- Scalable and performant

### Frontend

**HTMX**
- Modern interactivity without complex JavaScript frameworks
- Server-driven UI updates
- Progressive enhancement
- Reduces frontend complexity

**Tailwind CSS**
- Utility-first approach speeds development
- Consistent design system
- Small production bundle with CDN
- Responsive by default

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/api/health` | Health check |
| POST   | `/api/items` | Create item |
| GET    | `/api/items` | List all items |
| GET    | `/api/items/:id` | Get specific item |
| PUT    | `/api/items/:id` | Update item |
| DELETE | `/api/items/:id` | Delete item |

## Environment Configuration

```env
PORT=3000                           # Server port
MONGO_URI=mongodb://localhost:27017 # MongoDB connection
MONGO_DATABASE=lacpa                # Database name
```

## Folder Structure Explanation

- **config/** - Application configuration (database setup)
- **handler/** - HTTP handlers for API endpoints
- **models/** - Data structures/entities
- **repository/** - Database access layer (interface + implementation)
- **LACPA_Web/** - Frontend static files
  - **js/** - JavaScript files
  - **css/** - Custom stylesheets (currently using Tailwind CDN)
- **main.go** - Application entry point

## Design Principles

1. **Separation of Concerns**: Each layer has a specific responsibility
2. **Dependency Injection**: Handler receives repository interface
3. **Clean Architecture**: Dependencies point inward
4. **RESTful API**: Standard HTTP methods and status codes
5. **Progressive Enhancement**: Frontend works with and without JavaScript
