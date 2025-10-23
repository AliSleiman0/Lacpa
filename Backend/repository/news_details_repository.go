package repository

import (
	"context"

	"github.com/AliSleiman0/Lacpa/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type NewsDetailsPageRepository interface {
	GetNewsDetailsPageData(ctx context.Context) (*models.News, error)
}

type newsPageRepository struct {
	db *mongo.Database
}

func NewNewsDetailsPageRepository(db *mongo.Database) NewsDetailsPageRepository {
	return &newsPageRepository{
		db: db,
	}
}

func (r *newsPageRepository) GetNewsDetailsPageData(ctx context.Context) (*models.News, error) {
	collection := r.db.Collection("news_page")
	var newsPage models.News
	err := collection.FindOne(ctx, map[string]interface{}{}).Decode(&newsPage)
	if err != nil {
		return nil, err
	}
	return &newsPage, nil
}
