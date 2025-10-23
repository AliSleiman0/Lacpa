package repository

import (
	"context"

	"github.com/AliSleiman0/Lacpa/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type DiscoverPageRepository interface {
	GetDiscoverPageData(ctx context.Context) (*models.DiscoverPage, error)
}

type discoverPageRepository struct {
	db *mongo.Database
}

func NewDiscoverPageRepository(db *mongo.Database) DiscoverPageRepository {
	return &discoverPageRepository{
		db: db,
	}
}

func (r *discoverPageRepository) GetDiscoverPageData(ctx context.Context) (*models.DiscoverPage, error) {
	collection := r.db.Collection("discover_page")
	var discoverPage models.DiscoverPage
	err := collection.FindOne(ctx, map[string]interface{}{}).Decode(&discoverPage)
	if err != nil {
		return nil, err
	}
	return &discoverPage, nil
}
