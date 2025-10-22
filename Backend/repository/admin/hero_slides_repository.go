package admin

import (
	"context"
	"time"

	"github.com/AliSleiman0/Lacpa/models/admin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HeroSlideRepository struct {
	collection *mongo.Collection
}

func NewHeroSlideRepository(db *mongo.Database) *HeroSlideRepository {
	return &HeroSlideRepository{
		collection: db.Collection("hero_slides"),
	}
}

// CreateSlide creates a new hero slide
func (r *HeroSlideRepository) CreateSlide(ctx context.Context, slide *admin.HeroSlide) error {
	slide.ID = primitive.NewObjectID()
	slide.CreatedAt = time.Now()
	slide.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, slide)
	return err
}

// GetSlideByID retrieves a slide by its ID
func (r *HeroSlideRepository) GetSlideByID(ctx context.Context, id string) (*admin.HeroSlide, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var slide admin.HeroSlide
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&slide)
	if err != nil {
		return nil, err
	}

	return &slide, nil
}

// GetAllSlides retrieves all slides ordered by orderIndex
func (r *HeroSlideRepository) GetAllSlides(ctx context.Context) ([]*admin.HeroSlide, error) {
	opts := options.Find().SetSort(bson.D{{Key: "orderIndex", Value: 1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var slides []*admin.HeroSlide
	if err = cursor.All(ctx, &slides); err != nil {
		return nil, err
	}

	return slides, nil
}

// GetActiveSlides retrieves all active slides ordered by orderIndex
func (r *HeroSlideRepository) GetActiveSlides(ctx context.Context) ([]*admin.HeroSlide, error) {
	opts := options.Find().SetSort(bson.D{{Key: "orderIndex", Value: 1}})
	cursor, err := r.collection.Find(ctx, bson.M{"isActive": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var slides []*admin.HeroSlide
	if err = cursor.All(ctx, &slides); err != nil {
		return nil, err
	}

	return slides, nil
}

// UpdateSlide updates an existing slide
func (r *HeroSlideRepository) UpdateSlide(ctx context.Context, id string, slide *admin.HeroSlide) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	slide.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"title":             slide.Title,
			"description":       slide.Description,
			"imgSrc":            slide.ImgSrc,
			"buttonTitle":       slide.ButtonTitle,
			"buttonLink":        slide.ButtonLink,
			"isActive":          slide.IsActive,
			"imageActive":       slide.ImageActive,
			"buttonActive":      slide.ButtonActive,
			"titleActive":       slide.TitleActive,
			"descriptionActive": slide.DescriptionActive,
			"orderIndex":        slide.OrderIndex,
			"updatedAt":         slide.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

// DeleteSlide deletes a slide by its ID
func (r *HeroSlideRepository) DeleteSlide(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// GetSlideCount returns the total number of slides
func (r *HeroSlideRepository) GetSlideCount(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}
