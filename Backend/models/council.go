package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CouncilPositionType represents fixed council positions
type CouncilPositionType string

const (
	PositionPresident      CouncilPositionType = "President"
	PositionVicePresident  CouncilPositionType = "Vice President"
	PositionBoardTreasurer CouncilPositionType = "Board Treasurer"
	PositionBoardSecretary CouncilPositionType = "Board Secretary"
	PositionBoardMember    CouncilPositionType = "Board Member"
	PositionNonCouncil     CouncilPositionType = "Non-Council Member" // Regular member
)

// Council position capacity constraints
const (
	MaxPresidents       = 1
	MaxVicePresidents   = 1
	MaxBoardTreasurers  = 1
	MaxBoardSecretaries = 1
	MaxBoardMembers     = 6
	MaxNonCouncil       = -1 // Unlimited
)

// CouncilPosition represents a member's position in the council
type CouncilPosition struct {
	ID        primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	MemberID  primitive.ObjectID  `json:"member_id" bson:"member_id"`   // Reference to IndividualMember
	CouncilID primitive.ObjectID  `json:"council_id" bson:"council_id"` // Reference to Council
	Position  CouncilPositionType `json:"position" bson:"position"`     // Position type
	StartDate time.Time           `json:"start_date" bson:"start_date"` // Term start date
	EndDate   time.Time           `json:"end_date" bson:"end_date"`     // Term end date
	IsActive  bool                `json:"is_active" bson:"is_active"`   // Currently serving
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time           `json:"updated_at" bson:"updated_at"`
}

// Council represents a council term/session
type Council struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`               // e.g., "LACPA Council 2024-2026"
	StartDate   time.Time          `json:"start_date" bson:"start_date"`   // Council term start
	EndDate     time.Time          `json:"end_date" bson:"end_date"`       // Council term end
	IsActive    bool               `json:"is_active" bson:"is_active"`     // Currently active council
	Description string             `json:"description" bson:"description"` // Council description
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// CouncilComposition represents the current council makeup
type CouncilComposition struct {
	CouncilID      primitive.ObjectID `json:"council_id"`
	President      *CouncilPosition   `json:"president"`
	VicePresident  *CouncilPosition   `json:"vice_president"`
	BoardTreasurer *CouncilPosition   `json:"board_treasurer"`
	BoardSecretary *CouncilPosition   `json:"board_secretary"`
	BoardMembers   []CouncilPosition  `json:"board_members"` // Should have exactly 6
	TotalPositions int                `json:"total_positions"`
}

// CouncilMemberWithDetails combines position with member information
type CouncilMemberWithDetails struct {
	Position CouncilPosition  `json:"position"`
	Member   IndividualMember `json:"member"`
}

// CouncilCompositionWithDetails includes full member details
type CouncilCompositionWithDetails struct {
	CouncilID      primitive.ObjectID         `json:"council_id"`
	CouncilName    string                     `json:"council_name"`
	President      *CouncilMemberWithDetails  `json:"president"`
	VicePresident  *CouncilMemberWithDetails  `json:"vice_president"`
	BoardTreasurer *CouncilMemberWithDetails  `json:"board_treasurer"`
	BoardSecretary *CouncilMemberWithDetails  `json:"board_secretary"`
	BoardMembers   []CouncilMemberWithDetails `json:"board_members"`
	TotalPositions int                        `json:"total_positions"`
}

// Helper Methods

// IsCouncilPosition checks if a position is a council position (not non-council)
func (p CouncilPositionType) IsCouncilPosition() bool {
	return p != PositionNonCouncil
}

// GetMaxCapacity returns the maximum number of members for this position
func (p CouncilPositionType) GetMaxCapacity() int {
	switch p {
	case PositionPresident:
		return MaxPresidents
	case PositionVicePresident:
		return MaxVicePresidents
	case PositionBoardTreasurer:
		return MaxBoardTreasurers
	case PositionBoardSecretary:
		return MaxBoardSecretaries
	case PositionBoardMember:
		return MaxBoardMembers
	case PositionNonCouncil:
		return MaxNonCouncil // Unlimited
	default:
		return 0
	}
}

// IsLeadershipPosition checks if position is president or vice president
func (p CouncilPositionType) IsLeadershipPosition() bool {
	return p == PositionPresident || p == PositionVicePresident
}

// GetPositionPriority returns priority for ordering (lower = higher priority)
func (p CouncilPositionType) GetPositionPriority() int {
	switch p {
	case PositionPresident:
		return 1
	case PositionVicePresident:
		return 2
	case PositionBoardTreasurer:
		return 3
	case PositionBoardSecretary:
		return 4
	case PositionBoardMember:
		return 5
	case PositionNonCouncil:
		return 6
	default:
		return 99
	}
}

// IsTermActive checks if the position term is currently active
func (cp *CouncilPosition) IsTermActive() bool {
	now := time.Now()
	return cp.IsActive &&
		now.After(cp.StartDate) &&
		(cp.EndDate.IsZero() || now.Before(cp.EndDate))
}

// GetAllCouncilPositionTypes returns all valid council position types
func GetAllCouncilPositionTypes() []CouncilPositionType {
	return []CouncilPositionType{
		PositionPresident,
		PositionVicePresident,
		PositionBoardTreasurer,
		PositionBoardSecretary,
		PositionBoardMember,
		PositionNonCouncil,
	}
}

// ValidateCouncilComposition checks if council composition is valid
func (cc *CouncilComposition) ValidateCouncilComposition() []string {
	errors := []string{}

	if cc.President == nil {
		errors = append(errors, "Council must have a President")
	}
	if cc.VicePresident == nil {
		errors = append(errors, "Council must have a Vice President")
	}
	if cc.BoardTreasurer == nil {
		errors = append(errors, "Council must have a Board Treasurer")
	}
	if cc.BoardSecretary == nil {
		errors = append(errors, "Council must have a Board Secretary")
	}
	if len(cc.BoardMembers) != MaxBoardMembers {
		errors = append(errors, "Council must have exactly 6 Board Members")
	}

	return errors
}

// GetRemainingSlots returns how many slots are available for a position type
func (cc *CouncilComposition) GetRemainingSlots(positionType CouncilPositionType) int {
	maxCapacity := positionType.GetMaxCapacity()
	if maxCapacity == -1 {
		return -1 // Unlimited
	}

	current := 0
	switch positionType {
	case PositionPresident:
		if cc.President != nil {
			current = 1
		}
	case PositionVicePresident:
		if cc.VicePresident != nil {
			current = 1
		}
	case PositionBoardTreasurer:
		if cc.BoardTreasurer != nil {
			current = 1
		}
	case PositionBoardSecretary:
		if cc.BoardSecretary != nil {
			current = 1
		}
	case PositionBoardMember:
		current = len(cc.BoardMembers)
	}

	return maxCapacity - current
}
