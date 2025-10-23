package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// PresidentInfo represents the information about the LACPA president
type PresidentInfo struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Title       string             `json:"title" bson:"title"`
	Image       string             `json:"image" bson:"image"`
	Description string             `json:"description" bson:"description"`
}

// BoardMember represents an individual board member
type BoardMember struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Position string             `json:"position" bson:"position"`
	Image    string             `json:"image" bson:"image"`
}

// TextContent represents different sections of text content on the discover page
type TextContent struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Section    string             `json:"section" bson:"section"`
	Title      string             `json:"title" bson:"title"`
	Content    string             `json:"content" bson:"content"`
	OrderIndex int                `json:"orderIndex" bson:"orderIndex"`
}

// DiscoverPage represents the complete structure of the discover LACPA page
type DiscoverPage struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	President    PresidentInfo      `json:"president" bson:"president"`
	BoardMembers []BoardMember      `json:"boardMembers" bson:"boardMembers"`
	TextSections []TextContent      `json:"textSections" bson:"textSections"`
	UpdatedAt    primitive.DateTime `json:"updatedAt" bson:"updatedAt"`
}
