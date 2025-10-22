package admin

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HeroSlide struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title             string             `json:"title" bson:"title"`
	Description       string             `json:"description" bson:"description"`
	ImgSrc            string             `json:"imgSrc" bson:"imgSrc"`
	ButtonTitle       string             `json:"buttonTitle" bson:"buttonTitle"`
	ButtonLink        string             `json:"buttonLink" bson:"buttonLink"`
	IsActive          bool               `json:"isActive" bson:"isActive"`
	ImageActive       bool               `json:"imageActive" bson:"imageActive"`
	ButtonActive      bool               `json:"buttonActive" bson:"buttonActive"`
	TitleActive       bool               `json:"titleActive" bson:"titleActive"`
	DescriptionActive bool               `json:"descriptionActive" bson:"descriptionActive"`
	OrderIndex        int                `json:"orderIndex" bson:"orderIndex"`
	CreatedAt         time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt         time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// CreateSlideRequest represents the request body for creating a slide
type CreateSlideRequest struct {
	Title             string `json:"title" form:"title"`
	Description       string `json:"description" form:"description"`
	ImgSrc            string `json:"imgSrc" form:"imgSrc"`
	ButtonTitle       string `json:"buttonTitle" form:"buttonTitle"`
	ButtonLink        string `json:"buttonLink" form:"buttonLink"`
	IsActive          bool   `json:"isActive" form:"isActive"`
	ImageActive       bool   `json:"imageActive" form:"imageActive"`
	ButtonActive      bool   `json:"buttonActive" form:"buttonActive"`
	TitleActive       bool   `json:"titleActive" form:"titleActive"`
	DescriptionActive bool   `json:"descriptionActive" form:"descriptionActive"`
	OrderIndex        int    `json:"orderIndex" form:"orderIndex"`
}

// UpdateSlideRequest represents the request body for updating a slide (all fields optional)
type UpdateSlideRequest struct {
	Title             *string `json:"title,omitempty" form:"title"`
	Description       *string `json:"description,omitempty" form:"description"`
	ImgSrc            *string `json:"imgSrc,omitempty" form:"imgSrc"`
	ButtonTitle       *string `json:"buttonTitle,omitempty" form:"buttonTitle"`
	ButtonLink        *string `json:"buttonLink,omitempty" form:"buttonLink"`
	IsActive          *bool   `json:"isActive,omitempty" form:"isActive"`
	ImageActive       *bool   `json:"imageActive,omitempty" form:"imageActive"`
	ButtonActive      *bool   `json:"buttonActive,omitempty" form:"buttonActive"`
	TitleActive       *bool   `json:"titleActive,omitempty" form:"titleActive"`
	DescriptionActive *bool   `json:"descriptionActive,omitempty" form:"descriptionActive"`
	OrderIndex        *int    `json:"orderIndex,omitempty" form:"orderIndex"`
}
