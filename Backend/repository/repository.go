package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository serves as the central orchestrator for all repository operations.
//
// ROLE: Aggregate Interface Pattern
//   - Combines all individual repository interfaces (ItemRepository, UserRepository, etc.)
//     into one unified interface
//   - Provides a single dependency point for handlers/services instead of injecting
//     multiple separate repositories
//   - Simplifies dependency injection: inject one Repository instead of N repositories
//   - Makes the codebase more maintainable as new repositories are added
//
// USAGE:
//
//	repo := repository.NewMongoRepository(db)
//	handler := handler.NewHandler(repo)  // Single dependency
//
//	// Access any repository method:
//	items, err := repo.GetAllItems(ctx)
//	users, err := repo.GetAllUsers(ctx)
type Repository interface {
	ItemRepository
	UserRepository
	// Future repositories will be added here:
	// OrderRepository
}

// MongoRepositoryManager implements the Repository interface and acts as a composition root, to be registered in main.go.
//
// ROLE: Composition Manager
// - Implements the unified Repository interface by embedding all sub-repository interfaces
// - Acts as a container that holds concrete implementations of all repositories
// - Uses Go's interface embedding to automatically expose all methods from embedded repositories
// - Provides a single struct that satisfies the aggregate Repository interface
//
// PATTERN: This follows the Composition Root pattern where all dependencies are
// composed in one place, making dependency management centralized and clean
type MongoRepositoryManager struct {
	ItemRepository
	UserRepository
	// Future repositories will be added here as embedded interfaces
}

// NewMongoRepository is the factory function that creates and initializes all repositories.
//
// ROLE: Factory Function / Dependency Injection Root
// - Single entry point to create the entire repository system
// - Handles initialization of all sub-repositories with the shared database connection
// - Ensures all repositories are properly configured with the same MongoDB database
// - Returns the unified Repository interface, hiding implementation details
// - Centralizes repository creation logic, making it easy to modify or extend
//
// PARAMETERS:
//
//	db: Shared MongoDB database connection that will be passed to all sub-repositories
//
// RETURNS:
//
//	Repository: Unified interface providing access to all repository operations
//
// EXTENSIBILITY:
//
//	To add a new repository (e.g., OrderRepository):
//	1. Add OrderRepository to the Repository interface above
//	2. Add OrderRepository to the MongoRepositoryManager struct above
//	3. Add OrderRepository: NewOrderRepository(db) to the return statement below
func NewMongoRepository(db *mongo.Database) Repository {
	return &MongoRepositoryManager{
		ItemRepository: NewItemRepository(db),
		UserRepository: NewUserRepository(db),
		// Future repositories will be initialized here:
		// OrderRepository: NewOrderRepository(db),
	}
}
