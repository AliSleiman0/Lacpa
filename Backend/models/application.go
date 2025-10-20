package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ApplicationType represents the type of membership application
type ApplicationType string

const (
	ApplicationTypeIndividual ApplicationType = "Individual"
	ApplicationTypeFirm       ApplicationType = "Firm"
)

// ApplicationStatus represents the status of an application
type ApplicationStatus string

const (
	ApplicationStatusPending     ApplicationStatus = "Pending"
	ApplicationStatusUnderReview ApplicationStatus = "Under Review"
	ApplicationStatusApproved    ApplicationStatus = "Approved"
	ApplicationStatusRejected    ApplicationStatus = "Rejected"
)

// ApplicationRequirement represents a requirement for membership application
type ApplicationRequirement struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	Icon            string             `bson:"icon" json:"icon"`                         // Font Awesome icon class
	ApplicationType ApplicationType    `bson:"application_type" json:"application_type"` // Individual or Firm
	IsRequired      bool               `bson:"is_required" json:"is_required"`
	OrderIndex      int                `bson:"order_index" json:"order_index"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// IndividualApplication represents an individual membership application
type IndividualApplication struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	// Personal Information
	FirstName   string    `bson:"first_name" json:"first_name"`
	MiddleName  string    `bson:"middle_name,omitempty" json:"middle_name,omitempty"`
	LastName    string    `bson:"last_name" json:"last_name"`
	DateOfBirth time.Time `bson:"date_of_birth" json:"date_of_birth"`
	Nationality string    `bson:"nationality" json:"nationality"`

	// Contact Information
	Email       string `bson:"email" json:"email"`
	Phone       string `bson:"phone" json:"phone"`
	MobilePhone string `bson:"mobile_phone,omitempty" json:"mobile_phone,omitempty"`

	// Address
	Street     string `bson:"street" json:"street"`
	City       string `bson:"city" json:"city"`
	District   string `bson:"district" json:"district"`
	Country    string `bson:"country" json:"country"`
	PostalCode string `bson:"postal_code,omitempty" json:"postal_code,omitempty"`

	// Professional Information
	ProfessionalTitle string   `bson:"professional_title" json:"professional_title"`
	Qualifications    []string `bson:"qualifications" json:"qualifications"`
	YearsOfExperience int      `bson:"years_of_experience" json:"years_of_experience"`
	CurrentEmployer   string   `bson:"current_employer,omitempty" json:"current_employer,omitempty"`

	// Documents (file paths or URLs)
	CVDocument            string   `bson:"cv_document,omitempty" json:"cv_document,omitempty"`
	CertificatesDocuments []string `bson:"certificates_documents,omitempty" json:"certificates_documents,omitempty"`
	IDDocument            string   `bson:"id_document,omitempty" json:"id_document,omitempty"`

	// Application Details
	Status      ApplicationStatus   `bson:"status" json:"status"`
	SubmittedAt time.Time           `bson:"submitted_at" json:"submitted_at"`
	ReviewedAt  *time.Time          `bson:"reviewed_at,omitempty" json:"reviewed_at,omitempty"`
	ReviewedBy  *primitive.ObjectID `bson:"reviewed_by,omitempty" json:"reviewed_by,omitempty"`
	ReviewNotes string              `bson:"review_notes,omitempty" json:"review_notes,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// FirmApplication represents a firm membership application
type FirmApplication struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Firm Information
	FirmName           string `bson:"firm_name" json:"firm_name"`
	TradeName          string `bson:"trade_name,omitempty" json:"trade_name,omitempty"`
	RegistrationNumber string `bson:"registration_number" json:"registration_number"`
	YearEstablished    int    `bson:"year_established" json:"year_established"`

	// Contact Information
	Email   string `bson:"email" json:"email"`
	Phone   string `bson:"phone" json:"phone"`
	Website string `bson:"website,omitempty" json:"website,omitempty"`

	// Address
	Street     string `bson:"street" json:"street"`
	City       string `bson:"city" json:"city"`
	District   string `bson:"district" json:"district"`
	Country    string `bson:"country" json:"country"`
	PostalCode string `bson:"postal_code,omitempty" json:"postal_code,omitempty"`

	// Firm Details
	NumberOfPartners  int      `bson:"number_of_partners" json:"number_of_partners"`
	NumberOfEmployees int      `bson:"number_of_employees" json:"number_of_employees"`
	ServicesOffered   []string `bson:"services_offered" json:"services_offered"`

	// Representative Information
	RepresentativeName  string `bson:"representative_name" json:"representative_name"`
	RepresentativeTitle string `bson:"representative_title" json:"representative_title"`
	RepresentativeEmail string `bson:"representative_email" json:"representative_email"`
	RepresentativePhone string `bson:"representative_phone" json:"representative_phone"`

	// Documents (file paths or URLs)
	RegistrationDocument string   `bson:"registration_document,omitempty" json:"registration_document,omitempty"`
	LicenseDocuments     []string `bson:"license_documents,omitempty" json:"license_documents,omitempty"`
	TaxCertificate       string   `bson:"tax_certificate,omitempty" json:"tax_certificate,omitempty"`

	// Application Details
	Status      ApplicationStatus   `bson:"status" json:"status"`
	SubmittedAt time.Time           `bson:"submitted_at" json:"submitted_at"`
	ReviewedAt  *time.Time          `bson:"reviewed_at,omitempty" json:"reviewed_at,omitempty"`
	ReviewedBy  *primitive.ObjectID `bson:"reviewed_by,omitempty" json:"reviewed_by,omitempty"`
	ReviewNotes string              `bson:"review_notes,omitempty" json:"review_notes,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
