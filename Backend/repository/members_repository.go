package repository

import (
	"context"

	"github.com/AliSleiman0/Lacpa/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MembersRepository defines the interface for member operations
type MembersRepository interface {
	// Individual Members
	GetIndividualMemberByID(ctx context.Context, id primitive.ObjectID) (*models.IndividualMember, error)
	GetAllIndividualMembers(ctx context.Context, page, pageSize int) ([]*models.IndividualMember, int64, error)
	GetIndividualMembersByType(ctx context.Context, memberType string, page, pageSize int) ([]*models.IndividualMember, int64, error)
	CreateIndividualMember(ctx context.Context, member *models.IndividualMember) error
	UpdateIndividualMember(ctx context.Context, member *models.IndividualMember) error
	DeleteIndividualMember(ctx context.Context, id primitive.ObjectID) error
	CountIndividualMembers(ctx context.Context) (int64, error)
	GetIndividualMemberMetrics(ctx context.Context) (*models.MemberMetrics, error)

	// Firm Members
	GetFirmMemberByID(ctx context.Context, id primitive.ObjectID) (*models.FirmMember, error)
	GetAllFirmMembers(ctx context.Context, page, pageSize int) ([]*models.FirmMember, int64, error)
	GetFirmMembersByType(ctx context.Context, firmType string, page, pageSize int) ([]*models.FirmMember, int64, error)
	GetFirmMembersBySize(ctx context.Context, firmSize string, page, pageSize int) ([]*models.FirmMember, int64, error)
	CreateFirmMember(ctx context.Context, firm *models.FirmMember) error
	UpdateFirmMember(ctx context.Context, firm *models.FirmMember) error
	DeleteFirmMember(ctx context.Context, id primitive.ObjectID) error
	CountFirmMembers(ctx context.Context) (int64, error)
	GetFirmMemberMetrics(ctx context.Context) (*models.FirmMetrics, error)
}

// membersRepository implements MembersRepository interface
type membersRepository struct {
	db                   *mongo.Database
	individualMembersCol *mongo.Collection
	firmMembersCol       *mongo.Collection
}

// NewMembersRepository creates a new members repository instance
func NewMembersRepository(db *mongo.Database) MembersRepository {
	return &membersRepository{
		db:                   db,
		individualMembersCol: db.Collection("individual_members"),
		firmMembersCol:       db.Collection("firm_members"),
	}
}

// ========================================
// INDIVIDUAL MEMBERS METHODS
// ========================================

// GetIndividualMemberByID retrieves a single individual member by ID
func (r *membersRepository) GetIndividualMemberByID(ctx context.Context, id primitive.ObjectID) (*models.IndividualMember, error) {
	var member models.IndividualMember
	err := r.individualMembersCol.FindOne(ctx, bson.M{"_id": id}).Decode(&member)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// GetAllIndividualMembers retrieves all individual members with pagination
func (r *membersRepository) GetAllIndividualMembers(ctx context.Context, page, pageSize int) ([]*models.IndividualMember, int64, error) {
	// Calculate skip value
	skip := int64((page - 1) * pageSize)

	// Set find options with pagination
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "last_name", Value: 1}, {Key: "first_name", Value: 1}})

	// Execute query
	cursor, err := r.individualMembersCol.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var members []*models.IndividualMember
	if err = cursor.All(ctx, &members); err != nil {
		return nil, 0, err
	}

	// Get total count
	total, err := r.individualMembersCol.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

// GetIndividualMembersByType retrieves individual members filtered by type with pagination
func (r *membersRepository) GetIndividualMembersByType(ctx context.Context, memberType string, page, pageSize int) ([]*models.IndividualMember, int64, error) {
	// Calculate skip value
	skip := int64((page - 1) * pageSize)

	// Build filter
	filter := bson.M{}
	if memberType != "" && memberType != "all" {
		filter["member_type"] = memberType
	}

	// Set find options with pagination
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "last_name", Value: 1}, {Key: "first_name", Value: 1}})

	// Execute query
	cursor, err := r.individualMembersCol.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var members []*models.IndividualMember
	if err = cursor.All(ctx, &members); err != nil {
		return nil, 0, err
	}

	// Get total count for this filter
	total, err := r.individualMembersCol.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

// CreateIndividualMember creates a new individual member
func (r *membersRepository) CreateIndividualMember(ctx context.Context, member *models.IndividualMember) error {
	result, err := r.individualMembersCol.InsertOne(ctx, member)
	if err != nil {
		return err
	}
	member.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// UpdateIndividualMember updates an existing individual member
func (r *membersRepository) UpdateIndividualMember(ctx context.Context, member *models.IndividualMember) error {
	filter := bson.M{"_id": member.ID}
	update := bson.M{"$set": member}
	_, err := r.individualMembersCol.UpdateOne(ctx, filter, update)
	return err
}

// DeleteIndividualMember deletes an individual member by ID
func (r *membersRepository) DeleteIndividualMember(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.individualMembersCol.DeleteOne(ctx, filter)
	return err
}

// CountIndividualMembers counts total individual members
func (r *membersRepository) CountIndividualMembers(ctx context.Context) (int64, error) {
	return r.individualMembersCol.CountDocuments(ctx, bson.M{})
}

// GetIndividualMemberMetrics retrieves statistics about individual members
func (r *membersRepository) GetIndividualMemberMetrics(ctx context.Context) (*models.MemberMetrics, error) {
	metrics := &models.MemberMetrics{}

	// Get total count
	total, err := r.individualMembersCol.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	metrics.TotalMembers = int(total)

	// Count by type
	apprentices, _ := r.individualMembersCol.CountDocuments(ctx, bson.M{"member_type": "Apprentices"})
	metrics.ApprenticesCount = int(apprentices)

	practicing, _ := r.individualMembersCol.CountDocuments(ctx, bson.M{"member_type": "Practicing"})
	metrics.PracticingCount = int(practicing)

	nonPracticing, _ := r.individualMembersCol.CountDocuments(ctx, bson.M{"member_type": "Non-Practicing"})
	metrics.NonPracticingCount = int(nonPracticing)

	retired, _ := r.individualMembersCol.CountDocuments(ctx, bson.M{"member_type": "Retired"})
	metrics.RetiredCount = int(retired)

	return metrics, nil
}

// ========================================
// FIRM MEMBERS METHODS
// ========================================

// GetFirmMemberByID retrieves a single firm member by ID
func (r *membersRepository) GetFirmMemberByID(ctx context.Context, id primitive.ObjectID) (*models.FirmMember, error) {
	var firm models.FirmMember
	err := r.firmMembersCol.FindOne(ctx, bson.M{"_id": id}).Decode(&firm)
	if err != nil {
		return nil, err
	}
	return &firm, nil
}

// GetAllFirmMembers retrieves all firm members with pagination
func (r *membersRepository) GetAllFirmMembers(ctx context.Context, page, pageSize int) ([]*models.FirmMember, int64, error) {
	// Calculate skip value
	skip := int64((page - 1) * pageSize)

	// Set find options with pagination
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "firm_name", Value: 1}})

	// Execute query
	cursor, err := r.firmMembersCol.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var firms []*models.FirmMember
	if err = cursor.All(ctx, &firms); err != nil {
		return nil, 0, err
	}

	// Get total count
	total, err := r.firmMembersCol.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return firms, total, nil
}

// GetFirmMembersByType retrieves firm members filtered by type with pagination
func (r *membersRepository) GetFirmMembersByType(ctx context.Context, firmType string, page, pageSize int) ([]*models.FirmMember, int64, error) {
	// Calculate skip value
	skip := int64((page - 1) * pageSize)

	// Build filter
	filter := bson.M{}
	if firmType != "" && firmType != "all" {
		filter["firm_type"] = firmType
	}

	// Set find options with pagination
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "firm_name", Value: 1}})

	// Execute query
	cursor, err := r.firmMembersCol.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var firms []*models.FirmMember
	if err = cursor.All(ctx, &firms); err != nil {
		return nil, 0, err
	}

	// Get total count for this filter
	total, err := r.firmMembersCol.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return firms, total, nil
}

// GetFirmMembersBySize retrieves firm members filtered by size with pagination
func (r *membersRepository) GetFirmMembersBySize(ctx context.Context, firmSize string, page, pageSize int) ([]*models.FirmMember, int64, error) {
	// Calculate skip value
	skip := int64((page - 1) * pageSize)

	// Build filter
	filter := bson.M{}
	if firmSize != "" && firmSize != "all" {
		filter["firm_size"] = firmSize
	}

	// Set find options with pagination
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "firm_name", Value: 1}})

	// Execute query
	cursor, err := r.firmMembersCol.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results
	var firms []*models.FirmMember
	if err = cursor.All(ctx, &firms); err != nil {
		return nil, 0, err
	}

	// Get total count for this filter
	total, err := r.firmMembersCol.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return firms, total, nil
}

// CreateFirmMember creates a new firm member
func (r *membersRepository) CreateFirmMember(ctx context.Context, firm *models.FirmMember) error {
	result, err := r.firmMembersCol.InsertOne(ctx, firm)
	if err != nil {
		return err
	}
	firm.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// UpdateFirmMember updates an existing firm member
func (r *membersRepository) UpdateFirmMember(ctx context.Context, firm *models.FirmMember) error {
	filter := bson.M{"_id": firm.ID}
	update := bson.M{"$set": firm}
	_, err := r.firmMembersCol.UpdateOne(ctx, filter, update)
	return err
}

// DeleteFirmMember deletes a firm member by ID
func (r *membersRepository) DeleteFirmMember(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.firmMembersCol.DeleteOne(ctx, filter)
	return err
}

// CountFirmMembers counts total firm members
func (r *membersRepository) CountFirmMembers(ctx context.Context) (int64, error) {
	return r.firmMembersCol.CountDocuments(ctx, bson.M{})
}

// GetFirmMemberMetrics retrieves statistics about firm members
func (r *membersRepository) GetFirmMemberMetrics(ctx context.Context) (*models.FirmMetrics, error) {
	metrics := &models.FirmMetrics{}

	// Get total count
	total, err := r.firmMembersCol.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	metrics.TotalFirms = int(total)

	// Count by type
	auditFirms, _ := r.firmMembersCol.CountDocuments(ctx, bson.M{"firm_type": "Audit Firm"})
	metrics.AuditFirmsCount = int(auditFirms)

	// Count by size
	big4, _ := r.firmMembersCol.CountDocuments(ctx, bson.M{"firm_size": "Big 4"})
	metrics.Big4Count = int(big4)

	large, _ := r.firmMembersCol.CountDocuments(ctx, bson.M{"firm_size": "Large"})
	metrics.LargeFirmsCount = int(large)

	medium, _ := r.firmMembersCol.CountDocuments(ctx, bson.M{"firm_size": "Medium"})
	metrics.MediumFirmsCount = int(medium)

	small, _ := r.firmMembersCol.CountDocuments(ctx, bson.M{"firm_size": "Small"})
	metrics.SmallFirmsCount = int(small)

	return metrics, nil
}
