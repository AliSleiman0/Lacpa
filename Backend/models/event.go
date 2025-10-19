package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EventCategory represents the category/type of an event
type EventCategory string

const (
	CategoryCongress           EventCategory = "congress"
	CategoryWorkshops          EventCategory = "workshops"
	CategoryProfessionalEvents EventCategory = "professional_events"
	CategorySocialEvents       EventCategory = "social_events"
	CategoryOtherAnnouncements EventCategory = "other_announcements"
)

// GetAllEventCategories returns all available event categories
func GetAllEventCategories() []EventCategory {
	return []EventCategory{
		CategoryCongress,
		CategoryWorkshops,
		CategoryProfessionalEvents,
		CategorySocialEvents,
		CategoryOtherAnnouncements,
	}
}

// String returns the string representation of the event category
func (ec EventCategory) String() string {
	return string(ec)
}

// GetDisplayName returns the user-friendly display name for the category
func (ec EventCategory) GetDisplayName() string {
	switch ec {
	case CategoryCongress:
		return "Congress"
	case CategoryWorkshops:
		return "Workshops"
	case CategoryProfessionalEvents:
		return "Professional Events"
	case CategorySocialEvents:
		return "Social Events"
	case CategoryOtherAnnouncements:
		return "Other Announcements"
	default:
		return "Unknown"
	}
}

// IsValidCategory checks if the category is valid
func (ec EventCategory) IsValid() bool {
	switch ec {
	case CategoryCongress, CategoryWorkshops, CategoryProfessionalEvents,
		CategorySocialEvents, CategoryOtherAnnouncements:
		return true
	default:
		return false
	}
}

// Event represents an event or news item
type Event struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Category    EventCategory      `json:"category" bson:"category"`
	StartDate   time.Time          `json:"start_date" bson:"start_date"`
	EndDate     time.Time          `json:"end_date" bson:"end_date"`
	CPEHours    int                `json:"cpe_hours" bson:"cpe_hours"` // Continuing Professional Education hours
	ImageURL    string             `json:"image_url" bson:"image_url,omitempty"`
	IsPublished bool               `json:"is_published" bson:"is_published"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// IsActive checks if the event is currently active/ongoing
func (e *Event) IsActive() bool {
	now := time.Now()
	return now.After(e.StartDate) && now.Before(e.EndDate)
}

// IsUpcoming checks if the event is in the future
func (e *Event) IsUpcoming() bool {
	return time.Now().Before(e.StartDate)
}

// IsPast checks if the event has ended
func (e *Event) IsPast() bool {
	return time.Now().After(e.EndDate)
}

// GetDuration returns the duration of the event in days
func (e *Event) GetDuration() int {
	duration := e.EndDate.Sub(e.StartDate)
	return int(duration.Hours() / 24)
}

// GetFormattedDateRange returns a formatted string of the date range
func (e *Event) GetFormattedDateRange() string {
	if e.StartDate.Format("2006-01-02") == e.EndDate.Format("2006-01-02") {
		// Same day event
		return e.StartDate.Format("02/01/2006")
	}
	// Multi-day event
	return e.StartDate.Format("02/01/2006") + " - " + e.EndDate.Format("02/01/2006")
}

// EventsGroupedByCategory represents events organized by their categories
type EventsGroupedByCategory struct {
	Congress           []Event `json:"congress"`
	Workshops          []Event `json:"workshops"`
	ProfessionalEvents []Event `json:"professional_events"`
	SocialEvents       []Event `json:"social_events"`
	OtherAnnouncements []Event `json:"other_announcements"`
	TotalCount         int     `json:"total_count"`
}

// EventFilter represents filtering options for events
type EventFilter struct {
	Category   *EventCategory `json:"category,omitempty"`
	StartDate  *time.Time     `json:"start_date,omitempty"`
	EndDate    *time.Time     `json:"end_date,omitempty"`
	IsActive   *bool          `json:"is_active,omitempty"`
	IsUpcoming *bool          `json:"is_upcoming,omitempty"`
	IsPast     *bool          `json:"is_past,omitempty"`
	Limit      int            `json:"limit,omitempty"`
	Offset     int            `json:"offset,omitempty"`
}
