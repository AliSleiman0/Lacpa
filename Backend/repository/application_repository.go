package repository

import (
	"context"
	"time"

	"github.com/AliSleiman0/Lacpa/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ApplicationRepository interface {
	// Application Requirements
	GetAllRequirements(ctx context.Context) ([]models.ApplicationRequirement, error)
	GetRequirementsByType(ctx context.Context, appType models.ApplicationType) ([]models.ApplicationRequirement, error)
	CreateRequirement(ctx context.Context, requirement *models.ApplicationRequirement) error
	UpdateRequirement(ctx context.Context, id primitive.ObjectID, requirement *models.ApplicationRequirement) error
	DeleteRequirement(ctx context.Context, id primitive.ObjectID) error

	// Individual Applications
	CreateIndividualApplication(ctx context.Context, application *models.IndividualApplication) error
	GetIndividualApplicationByID(ctx context.Context, id primitive.ObjectID) (*models.IndividualApplication, error)
	GetAllIndividualApplications(ctx context.Context) ([]models.IndividualApplication, error)
	GetIndividualApplicationsByStatus(ctx context.Context, status models.ApplicationStatus) ([]models.IndividualApplication, error)
	UpdateIndividualApplicationStatus(ctx context.Context, id primitive.ObjectID, status models.ApplicationStatus, notes string, reviewedBy primitive.ObjectID) error

	// Firm Applications
	CreateFirmApplication(ctx context.Context, application *models.FirmApplication) error
	GetFirmApplicationByID(ctx context.Context, id primitive.ObjectID) (*models.FirmApplication, error)
	GetAllFirmApplications(ctx context.Context) ([]models.FirmApplication, error)
	GetFirmApplicationsByStatus(ctx context.Context, status models.ApplicationStatus) ([]models.FirmApplication, error)
	UpdateFirmApplicationStatus(ctx context.Context, id primitive.ObjectID, status models.ApplicationStatus, notes string, reviewedBy primitive.ObjectID) error
}

type applicationRepository struct {
	requirementCollection           *mongo.Collection
	individualApplicationCollection *mongo.Collection
	firmApplicationCollection       *mongo.Collection
}

func NewApplicationRepository(db *mongo.Database) ApplicationRepository {
	return &applicationRepository{
		requirementCollection:           db.Collection("application_requirements"),
		individualApplicationCollection: db.Collection("individual_applications"),
		firmApplicationCollection:       db.Collection("firm_applications"),
	}
}

// ============= Application Requirements =============

func (r *applicationRepository) GetAllRequirements(ctx context.Context) ([]models.ApplicationRequirement, error) {
	cursor, err := r.requirementCollection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "application_type", Value: 1}, {Key: "order_index", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var requirements []models.ApplicationRequirement
	if err := cursor.All(ctx, &requirements); err != nil {
		return nil, err
	}

	return requirements, nil
}

func (r *applicationRepository) GetRequirementsByType(ctx context.Context, appType models.ApplicationType) ([]models.ApplicationRequirement, error) {
	cursor, err := r.requirementCollection.Find(ctx, bson.M{"application_type": appType}, options.Find().SetSort(bson.D{{Key: "order_index", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var requirements []models.ApplicationRequirement
	if err := cursor.All(ctx, &requirements); err != nil {
		return nil, err
	}

	return requirements, nil
}

func (r *applicationRepository) CreateRequirement(ctx context.Context, requirement *models.ApplicationRequirement) error {
	requirement.CreatedAt = time.Now()
	requirement.UpdatedAt = time.Now()

	result, err := r.requirementCollection.InsertOne(ctx, requirement)
	if err != nil {
		return err
	}

	requirement.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *applicationRepository) UpdateRequirement(ctx context.Context, id primitive.ObjectID, requirement *models.ApplicationRequirement) error {
	requirement.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"title":            requirement.Title,
			"description":      requirement.Description,
			"icon":             requirement.Icon,
			"application_type": requirement.ApplicationType,
			"is_required":      requirement.IsRequired,
			"order_index":      requirement.OrderIndex,
			"updated_at":       requirement.UpdatedAt,
		},
	}

	_, err := r.requirementCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *applicationRepository) DeleteRequirement(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.requirementCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// ============= Individual Applications =============

func (r *applicationRepository) CreateIndividualApplication(ctx context.Context, application *models.IndividualApplication) error {
	application.Status = models.ApplicationStatusPending
	application.SubmittedAt = time.Now()
	application.CreatedAt = time.Now()
	application.UpdatedAt = time.Now()

	result, err := r.individualApplicationCollection.InsertOne(ctx, application)
	if err != nil {
		return err
	}

	application.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *applicationRepository) GetIndividualApplicationByID(ctx context.Context, id primitive.ObjectID) (*models.IndividualApplication, error) {
	var application models.IndividualApplication
	err := r.individualApplicationCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&application)
	if err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepository) GetAllIndividualApplications(ctx context.Context) ([]models.IndividualApplication, error) {
	cursor, err := r.individualApplicationCollection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "submitted_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var applications []models.IndividualApplication
	if err := cursor.All(ctx, &applications); err != nil {
		return nil, err
	}

	return applications, nil
}

func (r *applicationRepository) GetIndividualApplicationsByStatus(ctx context.Context, status models.ApplicationStatus) ([]models.IndividualApplication, error) {
	cursor, err := r.individualApplicationCollection.Find(ctx, bson.M{"status": status}, options.Find().SetSort(bson.D{{Key: "submitted_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var applications []models.IndividualApplication
	if err := cursor.All(ctx, &applications); err != nil {
		return nil, err
	}

	return applications, nil
}

func (r *applicationRepository) UpdateIndividualApplicationStatus(ctx context.Context, id primitive.ObjectID, status models.ApplicationStatus, notes string, reviewedBy primitive.ObjectID) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":       status,
			"review_notes": notes,
			"reviewed_by":  reviewedBy,
			"reviewed_at":  now,
			"updated_at":   now,
		},
	}

	_, err := r.individualApplicationCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// ============= Firm Applications =============

func (r *applicationRepository) CreateFirmApplication(ctx context.Context, application *models.FirmApplication) error {
	application.Status = models.ApplicationStatusPending
	application.SubmittedAt = time.Now()
	application.CreatedAt = time.Now()
	application.UpdatedAt = time.Now()

	result, err := r.firmApplicationCollection.InsertOne(ctx, application)
	if err != nil {
		return err
	}

	application.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *applicationRepository) GetFirmApplicationByID(ctx context.Context, id primitive.ObjectID) (*models.FirmApplication, error) {
	var application models.FirmApplication
	err := r.firmApplicationCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&application)
	if err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepository) GetAllFirmApplications(ctx context.Context) ([]models.FirmApplication, error) {
	cursor, err := r.firmApplicationCollection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "submitted_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var applications []models.FirmApplication
	if err := cursor.All(ctx, &applications); err != nil {
		return nil, err
	}

	return applications, nil
}

func (r *applicationRepository) GetFirmApplicationsByStatus(ctx context.Context, status models.ApplicationStatus) ([]models.FirmApplication, error) {
	cursor, err := r.firmApplicationCollection.Find(ctx, bson.M{"status": status}, options.Find().SetSort(bson.D{{Key: "submitted_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var applications []models.FirmApplication
	if err := cursor.All(ctx, &applications); err != nil {
		return nil, err
	}

	return applications, nil
}

func (r *applicationRepository) UpdateFirmApplicationStatus(ctx context.Context, id primitive.ObjectID, status models.ApplicationStatus, notes string, reviewedBy primitive.ObjectID) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":       status,
			"review_notes": notes,
			"reviewed_by":  reviewedBy,
			"reviewed_at":  now,
			"updated_at":   now,
		},
	}

	_, err := r.firmApplicationCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
