# Contributing to Lacpa

Thank you for your interest in contributing to Lacpa! This document provides guidelines for contributing to the project.

## Development Setup

1. Fork and clone the repository
2. Follow the [QUICKSTART.md](QUICKSTART.md) guide to set up your development environment
3. Create a new branch for your feature or bugfix

```bash
git checkout -b feature/your-feature-name
```

## Code Standards

### Go Code Style

- Follow standard Go conventions and idioms
- Run `go fmt` before committing
- Run `go vet` to catch common mistakes
- Write clear, descriptive variable and function names
- Add comments for exported functions and complex logic

```bash
make fmt
make vet
```

### Project Structure

Maintain the existing project structure:

```
.
├── config/         # Configuration and setup
├── models/         # Data models
├── repository/     # Data access layer
├── handler/        # HTTP handlers
├── LACPA_Web/      # Frontend files
│   ├── js/         # JavaScript files
│   └── css/        # Stylesheets
└── main.go         # Application entry point
```

### Adding New Features

#### Backend (Go)

1. **Models**: Add new data structures in `models/`
2. **Repository**: Add data access methods in `repository/`
3. **Handler**: Add HTTP handlers in `handler/`
4. **Routes**: Register routes in `handler/handler.go` `SetupRoutes` method

Example flow for adding a new entity:

```go
// 1. models/user.go
type User struct {
    ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name string            `json:"name" bson:"name"`
}

// 2. repository/repository.go
func (r *MongoRepository) CreateUser(ctx context.Context, user *models.User) error {
    // Implementation
}

// 3. handler/handler.go
func (h *Handler) CreateUser(c *fiber.Ctx) error {
    // Implementation
}

// 4. Register route
users := router.Group("/users")
users.Post("/", h.CreateUser)
```

#### Frontend (HTMX + Tailwind)

1. Update `LACPA_Web/index.html` for UI changes
2. Add custom JavaScript in `LACPA_Web/js/app.js` if needed
3. Use Tailwind utility classes for styling
4. Use HTMX attributes for dynamic behavior

Example HTMX pattern:

```html
<button hx-post="/api/endpoint"
        hx-target="#result"
        hx-swap="innerHTML"
        class="bg-blue-600 text-white px-4 py-2 rounded">
    Submit
</button>
```

## Testing

### Writing Tests

- Place test files next to the code they test
- Use descriptive test names
- Test both success and error cases
- Keep tests focused and independent

```go
func TestCreateItem(t *testing.T) {
    t.Run("valid item creation", func(t *testing.T) {
        // Test implementation
    })
    
    t.Run("invalid item data", func(t *testing.T) {
        // Test implementation
    })
}
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./models/
```

## Git Workflow

### Commit Messages

Use clear, descriptive commit messages:

```
feat: add user authentication
fix: resolve item deletion bug
docs: update API documentation
refactor: simplify handler logic
test: add repository tests
```

### Pull Request Process

1. Ensure all tests pass
2. Update documentation if needed
3. Format your code: `make fmt`
4. Verify code quality: `make vet`
5. Push to your fork
6. Create a Pull Request with a clear description

## API Design Guidelines

### RESTful Conventions

- Use appropriate HTTP methods (GET, POST, PUT, DELETE)
- Return appropriate status codes
- Use JSON for request/response bodies
- Follow existing endpoint patterns

### Status Codes

- `200 OK` - Successful GET, PUT
- `201 Created` - Successful POST
- `204 No Content` - Successful DELETE
- `400 Bad Request` - Invalid input
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

### Response Format

Success response:
```json
{
    "id": "...",
    "name": "...",
    "created_at": "..."
}
```

Error response:
```json
{
    "error": "Error description"
}
```

## Frontend Guidelines

### HTMX Patterns

- Use appropriate HTTP methods via `hx-post`, `hx-get`, etc.
- Target specific elements with `hx-target`
- Control swap behavior with `hx-swap`
- Handle loading states and errors

### Tailwind CSS

- Use utility classes for styling
- Maintain consistent spacing and colors
- Follow responsive design patterns
- Keep custom CSS minimal

## Database Guidelines

### MongoDB Conventions

- Use meaningful collection names (plural)
- Include timestamps (created_at, updated_at)
- Use ObjectID for primary keys
- Index frequently queried fields

### Repository Pattern

- All database operations go through the repository
- Repository methods should be interface-based
- Handle errors appropriately
- Use context for cancellation and timeouts

## Documentation

### Code Documentation

- Document all exported functions and types
- Explain complex algorithms or business logic
- Include examples in comments when helpful

```go
// CreateItem inserts a new item into the database.
// It automatically sets CreatedAt and UpdatedAt timestamps.
// Returns an error if the item cannot be created.
func (r *MongoRepository) CreateItem(ctx context.Context, item *models.Item) error {
    // Implementation
}
```

### Project Documentation

- Update README.md for major features
- Update ARCHITECTURE.md for structural changes
- Update QUICKSTART.md if setup process changes
- Add examples for new features

## Common Tasks

### Adding a New API Endpoint

1. Define the model in `models/`
2. Add repository methods in `repository/`
3. Create handler in `handler/`
4. Register route in `SetupRoutes`
5. Add tests
6. Update API documentation

### Updating the Frontend

1. Modify `LACPA_Web/index.html`
2. Add JavaScript if needed in `LACPA_Web/js/app.js`
3. Test in browser
4. Ensure HTMX attributes are correct

### Adding Environment Variables

1. Add to `.env.example`
2. Use `os.Getenv()` in code with default values
3. Document in README.md

## Need Help?

- Check existing code for examples
- Read the [ARCHITECTURE.md](ARCHITECTURE.md) for system design
- Open an issue for questions or discussions
- Review closed PRs for similar changes

## License

By contributing to Lacpa, you agree that your contributions will be licensed under the MIT License.
