package repository

import (
	"context"
	"errors"
	"time"

	"github.com/AliSleiman0/Lacpa/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CouncilRepository defines operations for council management
type CouncilRepository interface {
	// Council operations
	GetActiveCouncil(ctx context.Context) (*models.Council, error)
	GetCouncilByID(ctx context.Context, id primitive.ObjectID) (*models.Council, error)
	GetAllCouncils(ctx context.Context) ([]models.Council, error)
	CreateCouncil(ctx context.Context, council *models.Council) error
	UpdateCouncil(ctx context.Context, id primitive.ObjectID, council *models.Council) error
	DeactivateCouncil(ctx context.Context, id primitive.ObjectID) error

	// Council position operations
	GetCouncilComposition(ctx context.Context, councilID primitive.ObjectID) (*models.CouncilComposition, error)
	GetCouncilCompositionWithDetails(ctx context.Context, councilID primitive.ObjectID) (*models.CouncilCompositionWithDetails, error)
	AssignCouncilPosition(ctx context.Context, position *models.CouncilPosition) error
	RemoveCouncilPosition(ctx context.Context, positionID primitive.ObjectID) error
	UpdateCouncilPosition(ctx context.Context, positionID primitive.ObjectID, position *models.CouncilPosition) error
	GetMemberCouncilHistory(ctx context.Context, memberID primitive.ObjectID) ([]models.CouncilPosition, error)
	GetPositionByID(ctx context.Context, positionID primitive.ObjectID) (*models.CouncilPosition, error)

	// Validation
	ValidatePositionAvailability(ctx context.Context, councilID primitive.ObjectID, positionType models.CouncilPositionType) (bool, error)
	GetAvailablePositions(ctx context.Context, councilID primitive.ObjectID) (map[models.CouncilPositionType]int, error)
}

type councilRepository struct {
	db                 *mongo.Database
	councilCollection  *mongo.Collection
	positionCollection *mongo.Collection
	memberCollection   *mongo.Collection
}

// NewCouncilRepository creates a new council repository instance
func NewCouncilRepository(db *mongo.Database) CouncilRepository {
	return &councilRepository{
		db:                 db,
		councilCollection:  db.Collection("councils"),
		positionCollection: db.Collection("council_positions"),
		memberCollection:   db.Collection("individual_members"),
	}
}

// GetActiveCouncil retrieves the currently active council
func (r *councilRepository) GetActiveCouncil(ctx context.Context) (*models.Council, error) {
	var council models.Council
	err := r.councilCollection.FindOne(ctx, bson.M{"is_active": true}).Decode(&council)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no active council found")
		}
		return nil, err
	}
	return &council, nil
}

// GetCouncilByID retrieves a council by its ID
func (r *councilRepository) GetCouncilByID(ctx context.Context, id primitive.ObjectID) (*models.Council, error) {
	var council models.Council
	err := r.councilCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&council)
	if err != nil {
		return nil, err
	}
	return &council, nil
}

// GetAllCouncils retrieves all councils
func (r *councilRepository) GetAllCouncils(ctx context.Context) ([]models.Council, error) {
	cursor, err := r.councilCollection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var councils []models.Council
	if err := cursor.All(ctx, &councils); err != nil {
		return nil, err
	}
	return councils, nil
}

// CreateCouncil creates a new council
func (r *councilRepository) CreateCouncil(ctx context.Context, council *models.Council) error {
	council.CreatedAt = time.Now()
	council.UpdatedAt = time.Now()

	// If this council is active, deactivate all other councils
	if council.IsActive {
		_, err := r.councilCollection.UpdateMany(ctx,
			bson.M{"is_active": true},
			bson.M{"$set": bson.M{"is_active": false, "updated_at": time.Now()}},
		)
		if err != nil {
			return err
		}
	}

	_, err := r.councilCollection.InsertOne(ctx, council)
	return err
}

// UpdateCouncil updates an existing council
func (r *councilRepository) UpdateCouncil(ctx context.Context, id primitive.ObjectID, council *models.Council) error {
	council.UpdatedAt = time.Now()

	// If activating this council, deactivate others
	if council.IsActive {
		_, err := r.councilCollection.UpdateMany(ctx,
			bson.M{"_id": bson.M{"$ne": id}, "is_active": true},
			bson.M{"$set": bson.M{"is_active": false, "updated_at": time.Now()}},
		)
		if err != nil {
			return err
		}
	}

	_, err := r.councilCollection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": council},
	)
	return err
}

// DeactivateCouncil deactivates a council
func (r *councilRepository) DeactivateCouncil(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.councilCollection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"is_active": false, "updated_at": time.Now()}},
	)
	return err
}

// GetCouncilComposition retrieves complete council composition
func (r *councilRepository) GetCouncilComposition(ctx context.Context, councilID primitive.ObjectID) (*models.CouncilComposition, error) {
	composition := &models.CouncilComposition{
		CouncilID: councilID,
	}

	// Get all active positions for this council
	cursor, err := r.positionCollection.Find(ctx, bson.M{
		"council_id": councilID,
		"is_active":  true,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var positions []models.CouncilPosition
	if err := cursor.All(ctx, &positions); err != nil {
		return nil, err
	}

	// Organize positions by type
	for _, pos := range positions {
		switch pos.Position {
		case models.PositionPresident:
			composition.President = &pos
		case models.PositionVicePresident:
			composition.VicePresident = &pos
		case models.PositionBoardTreasurer:
			composition.BoardTreasurer = &pos
		case models.PositionBoardSecretary:
			composition.BoardSecretary = &pos
		case models.PositionBoardMember:
			composition.BoardMembers = append(composition.BoardMembers, pos)
		}
	}

	composition.TotalPositions = len(positions)
	return composition, nil
}

// GetCouncilCompositionWithDetails retrieves council composition with full member details
func (r *councilRepository) GetCouncilCompositionWithDetails(ctx context.Context, councilID primitive.ObjectID) (*models.CouncilCompositionWithDetails, error) {
	// Get council info
	council, err := r.GetCouncilByID(ctx, councilID)
	if err != nil {
		return nil, err
	}

	composition := &models.CouncilCompositionWithDetails{
		CouncilID:   councilID,
		CouncilName: council.Name,
	}

	// Get all positions for this council (active or inactive)
	cursor, err := r.positionCollection.Find(ctx, bson.M{
		"council_id": councilID,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var positions []models.CouncilPosition
	if err := cursor.All(ctx, &positions); err != nil {
		return nil, err
	}

	// Fetch member details for each position
	for _, pos := range positions {
		var member models.IndividualMember
		err := r.memberCollection.FindOne(ctx, bson.M{"_id": pos.MemberID}).Decode(&member)
		if err != nil {
			continue // Skip if member not found
		}

		memberWithDetails := models.CouncilMemberWithDetails{
			Position: pos,
			Member:   member,
		}

		switch pos.Position {
		case models.PositionPresident:
			composition.President = &memberWithDetails
		case models.PositionVicePresident:
			composition.VicePresident = &memberWithDetails
		case models.PositionBoardTreasurer:
			composition.BoardTreasurer = &memberWithDetails
		case models.PositionBoardSecretary:
			composition.BoardSecretary = &memberWithDetails
		case models.PositionBoardMember:
			composition.BoardMembers = append(composition.BoardMembers, memberWithDetails)
		}
	}

	composition.TotalPositions = len(positions)
	return composition, nil
}

// AssignCouncilPosition assigns a member to a council position
func (r *councilRepository) AssignCouncilPosition(ctx context.Context, position *models.CouncilPosition) error {
	// Check capacity
	available, err := r.ValidatePositionAvailability(ctx, position.CouncilID, position.Position)
	if err != nil {
		return err
	}
	if !available {
		return errors.New("no available slots for this position type")
	}

	// Check if member already has an active position in this council
	count, err := r.positionCollection.CountDocuments(ctx, bson.M{
		"member_id":  position.MemberID,
		"council_id": position.CouncilID,
		"is_active":  true,
	})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("member already has an active position in this council")
	}

	// Insert position
	position.CreatedAt = time.Now()
	position.UpdatedAt = time.Now()
	position.IsActive = true

	result, err := r.positionCollection.InsertOne(ctx, position)
	if err != nil {
		return err
	}

	// Update member document with council position reference
	positionID := result.InsertedID.(primitive.ObjectID)
	_, err = r.memberCollection.UpdateOne(ctx,
		bson.M{"_id": position.MemberID},
		bson.M{
			"$set": bson.M{
				"current_council_position_id": positionID,
				"council_position":            position.Position,
				"is_council_member":           true,
				"updated_at":                  time.Now(),
			},
		},
	)

	return err
}

// RemoveCouncilPosition removes a member from a council position
func (r *councilRepository) RemoveCouncilPosition(ctx context.Context, positionID primitive.ObjectID) error {
	// Get position to find member ID
	var position models.CouncilPosition
	err := r.positionCollection.FindOne(ctx, bson.M{"_id": positionID}).Decode(&position)
	if err != nil {
		return err
	}

	// Deactivate position
	_, err = r.positionCollection.UpdateOne(ctx,
		bson.M{"_id": positionID},
		bson.M{
			"$set": bson.M{
				"is_active":  false,
				"end_date":   time.Now(),
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return err
	}

	// Update member document
	_, err = r.memberCollection.UpdateOne(ctx,
		bson.M{"_id": position.MemberID},
		bson.M{
			"$set": bson.M{
				"current_council_position_id": nil,
				"council_position":            models.PositionNonCouncil,
				"is_council_member":           false,
				"updated_at":                  time.Now(),
			},
		},
	)

	return err
}

// UpdateCouncilPosition updates a council position
func (r *councilRepository) UpdateCouncilPosition(ctx context.Context, positionID primitive.ObjectID, position *models.CouncilPosition) error {
	position.UpdatedAt = time.Now()
	_, err := r.positionCollection.UpdateOne(ctx,
		bson.M{"_id": positionID},
		bson.M{"$set": position},
	)
	return err
}

// GetMemberCouncilHistory retrieves all council positions for a member
func (r *councilRepository) GetMemberCouncilHistory(ctx context.Context, memberID primitive.ObjectID) ([]models.CouncilPosition, error) {
	cursor, err := r.positionCollection.Find(ctx,
		bson.M{"member_id": memberID},
		options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var positions []models.CouncilPosition
	if err := cursor.All(ctx, &positions); err != nil {
		return nil, err
	}
	return positions, nil
}

// GetPositionByID retrieves a council position by ID
func (r *councilRepository) GetPositionByID(ctx context.Context, positionID primitive.ObjectID) (*models.CouncilPosition, error) {
	var position models.CouncilPosition
	err := r.positionCollection.FindOne(ctx, bson.M{"_id": positionID}).Decode(&position)
	if err != nil {
		return nil, err
	}
	return &position, nil
}

// ValidatePositionAvailability checks if a position slot is available
func (r *councilRepository) ValidatePositionAvailability(ctx context.Context, councilID primitive.ObjectID, positionType models.CouncilPositionType) (bool, error) {
	maxCapacity := positionType.GetMaxCapacity()
	if maxCapacity == -1 {
		return true, nil // Unlimited capacity
	}

	// Count current active positions of this type
	count, err := r.positionCollection.CountDocuments(ctx, bson.M{
		"council_id": councilID,
		"position":   positionType,
		"is_active":  true,
	})
	if err != nil {
		return false, err
	}

	return count < int64(maxCapacity), nil
}

// GetAvailablePositions returns available slots for each position type
func (r *councilRepository) GetAvailablePositions(ctx context.Context, councilID primitive.ObjectID) (map[models.CouncilPositionType]int, error) {
	available := make(map[models.CouncilPositionType]int)

	for _, posType := range models.GetAllCouncilPositionTypes() {
		if posType == models.PositionNonCouncil {
			continue
		}

		maxCapacity := posType.GetMaxCapacity()
		count, err := r.positionCollection.CountDocuments(ctx, bson.M{
			"council_id": councilID,
			"position":   posType,
			"is_active":  true,
		})
		if err != nil {
			return nil, err
		}

		available[posType] = maxCapacity - int(count)
	}

	return available, nil
}
