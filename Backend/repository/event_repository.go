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

// EventRepository defines operations for event management
type EventRepository interface {
	// Event CRUD operations
	GetEventByID(ctx context.Context, id primitive.ObjectID) (*models.Event, error)
	GetAllEvents(ctx context.Context, filter *models.EventFilter) ([]models.Event, error)
	GetEventsByCategory(ctx context.Context, category models.EventCategory) ([]models.Event, error)
	GetEventsGroupedByCategory(ctx context.Context) (*models.EventsGroupedByCategory, error)
	GetEventsGroupedByCategoryPaginated(ctx context.Context, page, pageSize int) (*models.EventsGroupedByCategory, error)
	CountEventsByCategory(ctx context.Context, category *models.EventCategory) (int64, error)
	CreateEvent(ctx context.Context, event *models.Event) error
	UpdateEvent(ctx context.Context, id primitive.ObjectID, event *models.Event) error
	DeleteEvent(ctx context.Context, id primitive.ObjectID) error

	// Special queries
	GetUpcomingEvents(ctx context.Context, limit int) ([]models.Event, error)
	GetActiveEvents(ctx context.Context) ([]models.Event, error)
	GetPastEvents(ctx context.Context, limit int) ([]models.Event, error)
}

type eventRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewEventRepository creates a new event repository instance
func NewEventRepository(db *mongo.Database) EventRepository {
	return &eventRepository{
		db:         db,
		collection: db.Collection("events"),
	}
}

// GetEventByID retrieves an event by its ID
func (r *eventRepository) GetEventByID(ctx context.Context, id primitive.ObjectID) (*models.Event, error) {
	var event models.Event
	err := r.collection.FindOne(ctx, bson.M{"_id": id, "is_published": true}).Decode(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// GetAllEvents retrieves all events with optional filtering
func (r *eventRepository) GetAllEvents(ctx context.Context, filter *models.EventFilter) ([]models.Event, error) {
	// Build query filter
	queryFilter := bson.M{"is_published": true}

	if filter != nil {
		if filter.Category != nil {
			queryFilter["category"] = *filter.Category
		}
		if filter.StartDate != nil {
			queryFilter["start_date"] = bson.M{"$gte": *filter.StartDate}
		}
		if filter.EndDate != nil {
			queryFilter["end_date"] = bson.M{"$lte": *filter.EndDate}
		}
	}

	// Build options
	opts := options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}})
	if filter != nil {
		if filter.Limit > 0 {
			opts.SetLimit(int64(filter.Limit))
		}
		if filter.Offset > 0 {
			opts.SetSkip(int64(filter.Offset))
		}
	}

	cursor, err := r.collection.Find(ctx, queryFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []models.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// GetEventsByCategory retrieves events by category
func (r *eventRepository) GetEventsByCategory(ctx context.Context, category models.EventCategory) ([]models.Event, error) {
	cursor, err := r.collection.Find(ctx,
		bson.M{"category": category, "is_published": true},
		options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []models.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// GetEventsGroupedByCategory retrieves all events grouped by their categories
func (r *eventRepository) GetEventsGroupedByCategory(ctx context.Context) (*models.EventsGroupedByCategory, error) {
	grouped := &models.EventsGroupedByCategory{}

	// Get events for each category
	congress, err := r.GetEventsByCategory(ctx, models.CategoryCongress)
	if err != nil {
		return nil, err
	}
	grouped.Congress = congress

	workshops, err := r.GetEventsByCategory(ctx, models.CategoryWorkshops)
	if err != nil {
		return nil, err
	}
	grouped.Workshops = workshops

	professional, err := r.GetEventsByCategory(ctx, models.CategoryProfessionalEvents)
	if err != nil {
		return nil, err
	}
	grouped.ProfessionalEvents = professional

	social, err := r.GetEventsByCategory(ctx, models.CategorySocialEvents)
	if err != nil {
		return nil, err
	}
	grouped.SocialEvents = social

	other, err := r.GetEventsByCategory(ctx, models.CategoryOtherAnnouncements)
	if err != nil {
		return nil, err
	}
	grouped.OtherAnnouncements = other

	// Calculate total count
	grouped.TotalCount = len(congress) + len(workshops) + len(professional) + len(social) + len(other)

	return grouped, nil
}

// CreateEvent creates a new event
func (r *eventRepository) CreateEvent(ctx context.Context, event *models.Event) error {
	event.CreatedAt = time.Now()
	event.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, event)
	return err
}

// UpdateEvent updates an existing event
func (r *eventRepository) UpdateEvent(ctx context.Context, id primitive.ObjectID, event *models.Event) error {
	event.UpdatedAt = time.Now()

	_, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": event},
	)
	return err
}

// DeleteEvent deletes an event (soft delete by unpublishing)
func (r *eventRepository) DeleteEvent(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"is_published": false, "updated_at": time.Now()}},
	)
	return err
}

// GetUpcomingEvents retrieves upcoming events
func (r *eventRepository) GetUpcomingEvents(ctx context.Context, limit int) ([]models.Event, error) {
	now := time.Now()

	opts := options.Find().
		SetSort(bson.D{{Key: "start_date", Value: 1}}) // Ascending order (soonest first)

	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := r.collection.Find(ctx,
		bson.M{
			"start_date":   bson.M{"$gt": now},
			"is_published": true,
		},
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []models.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// GetActiveEvents retrieves currently active events
func (r *eventRepository) GetActiveEvents(ctx context.Context) ([]models.Event, error) {
	now := time.Now()

	cursor, err := r.collection.Find(ctx,
		bson.M{
			"start_date":   bson.M{"$lte": now},
			"end_date":     bson.M{"$gte": now},
			"is_published": true,
		},
		options.Find().SetSort(bson.D{{Key: "start_date", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []models.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// GetPastEvents retrieves past events
func (r *eventRepository) GetPastEvents(ctx context.Context, limit int) ([]models.Event, error) {
	now := time.Now()

	opts := options.Find().
		SetSort(bson.D{{Key: "end_date", Value: -1}}) // Descending order (most recent first)

	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := r.collection.Find(ctx,
		bson.M{
			"end_date":     bson.M{"$lt": now},
			"is_published": true,
		},
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var events []models.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	return events, nil
}

// GetEventsGroupedByCategoryPaginated retrieves events grouped by category with pagination
// This returns a maximum of pageSize events total (not per category)
func (r *eventRepository) GetEventsGroupedByCategoryPaginated(ctx context.Context, page, pageSize int) (*models.EventsGroupedByCategory, error) {
	grouped := &models.EventsGroupedByCategory{}

	skip := (page - 1) * pageSize

	// Get all events with pagination (sorted by date, newest first)
	opts := options.Find().
		SetSort(bson.D{{Key: "start_date", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := r.collection.Find(ctx,
		bson.M{"is_published": true},
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var allEvents []models.Event
	if err := cursor.All(ctx, &allEvents); err != nil {
		return nil, err
	}

	// Group the paginated events by category
	for _, event := range allEvents {
		switch event.Category {
		case models.CategoryCongress:
			grouped.Congress = append(grouped.Congress, event)
		case models.CategoryWorkshops:
			grouped.Workshops = append(grouped.Workshops, event)
		case models.CategoryProfessionalEvents:
			grouped.ProfessionalEvents = append(grouped.ProfessionalEvents, event)
		case models.CategorySocialEvents:
			grouped.SocialEvents = append(grouped.SocialEvents, event)
		case models.CategoryOtherAnnouncements:
			grouped.OtherAnnouncements = append(grouped.OtherAnnouncements, event)
		}
	}

	// Total count is the number of events on this page
	grouped.TotalCount = len(allEvents)

	return grouped, nil
}

// CountEventsByCategory counts events by category (or all if category is nil)
func (r *eventRepository) CountEventsByCategory(ctx context.Context, category *models.EventCategory) (int64, error) {
	filter := bson.M{"is_published": true}

	if category != nil {
		filter["category"] = *category
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}
