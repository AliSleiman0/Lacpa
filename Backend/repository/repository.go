package repository

import (
	
	"go.mongodb.org/mongo-driver/mongo"
)

// MainRepository defines the interface for main page operations
type MainRepository interface {
	// Add methods as needed for main page functionality
	// GetLandingPageData(ctx context.Context) (*models.LandingPage, error)
}

// mainRepository implements MainRepository interface
type mainRepository struct {
	db *mongo.Database
}

// NewMainRepository creates a new main repository instance
func NewMainRepository(db *mongo.Database) MainRepository {
	return &mainRepository{
		db: db,
	}
}

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
	MainRepository
}
type MongoRepositoryManager struct {
	MainRepository
	// Future repositories will be added here as embedded interfaces
	// ItemRepository
	// UserRepository
}

func NewMongoRepository(db *mongo.Database) Repository {
	return &MongoRepositoryManager{
		MainRepository: NewMainRepository(db),
		// Future repositories will be initialized here:
		// OrderRepository: NewOrderRepository(db),
	}
}
