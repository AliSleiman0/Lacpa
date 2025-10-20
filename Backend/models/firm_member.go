package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FirmMember represents a firm/company that is a member of LACPA
type FirmMember struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`

	// Basic Information
	LacpaID  string `json:"lacpa_id" bson:"lacpa_id"`   // Unique LACPA Firm ID: "F-1234"
	FirmName string `json:"firm_name" bson:"firm_name"` // "Deloitte Lebanon"
	LogoURL  string `json:"logo_url" bson:"logo_url"`   // Path to firm logo image

	// Classification
	FirmType   string `json:"firm_type" bson:"firm_type"`     // "Audit Firm", "Accounting Firm", "Consultancy"
	FirmSize   string `json:"firm_size" bson:"firm_size"`     // "Big 4", "Large", "Medium", "Small"
	BadgeEmoji string `json:"badge_emoji" bson:"badge_emoji"` // Visual badge: "🏢"
	BadgeColor string `json:"badge_color" bson:"badge_color"` // Tailwind color: "blue-500"

	// Contact Information
	PrimaryPhone   string `json:"primary_phone" bson:"primary_phone"`     // "+961 01 123 456"
	SecondaryPhone string `json:"secondary_phone" bson:"secondary_phone"` // Optional secondary phone
	PrimaryEmail   string `json:"primary_email" bson:"primary_email"`     // "info@firm.com"
	Website        string `json:"website" bson:"website"`                 // "www.firm.com"
	LinkedInURL    string `json:"linkedin_url" bson:"linkedin_url"`       // Company LinkedIn page
	FacebookURL    string `json:"facebook_url" bson:"facebook_url"`       // Company Facebook page
	TwitterURL     string `json:"twitter_url" bson:"twitter_url"`         // Company Twitter
	InstagramURL   string `json:"instagram_url" bson:"instagram_url"`     // Company Instagram

	// Primary Contact Person
	ContactPersonName     string `json:"contact_person_name" bson:"contact_person_name"`         // "John Doe"
	ContactPersonTitle    string `json:"contact_person_title" bson:"contact_person_title"`       // "Managing Partner"
	ContactPersonPhone    string `json:"contact_person_phone" bson:"contact_person_phone"`       // Direct line
	ContactPersonEmail    string `json:"contact_person_email" bson:"contact_person_email"`       // Direct email
	ContactPersonLinkedIn string `json:"contact_person_linkedin" bson:"contact_person_linkedin"` // Personal LinkedIn

	// Address Components
	FullAddress  string `json:"full_address" bson:"full_address"`   // Full formatted address
	BuildingName string `json:"building_name" bson:"building_name"` // "ABC Tower"
	Floor        string `json:"floor" bson:"floor"`                 // "12th Floor"
	Street       string `json:"street" bson:"street"`               // "Main Street"
	Area         string `json:"area" bson:"area"`                   // "Downtown"
	City         string `json:"city" bson:"city"`                   // "Beirut"
	District     string `json:"district" bson:"district"`           // "Beirut District"
	Governorate  string `json:"governorate" bson:"governorate"`     // "Beirut"
	PostalCode   string `json:"postal_code" bson:"postal_code"`     // "1107 2080"
	Country      string `json:"country" bson:"country"`             // "Lebanon"

	// Business Information
	YearEstablished   int    `json:"year_established" bson:"year_established"`       // 1995
	NumberOfPartners  int    `json:"number_of_partners" bson:"number_of_partners"`   // Number of partners
	NumberOfEmployees int    `json:"number_of_employees" bson:"number_of_employees"` // Total employees
	NumberOfCPAs      int    `json:"number_of_cpas" bson:"number_of_cpas"`           // Certified accountants
	AnnualRevenue     string `json:"annual_revenue" bson:"annual_revenue"`           // "$1M - $5M" (range)

	// Services & Specializations
	ServicesOffered []string `json:"services_offered" bson:"services_offered"` // ["Audit", "Tax", "Advisory"]
	Industries      []string `json:"industries" bson:"industries"`             // ["Banking", "Real Estate", "Technology"]
	Specializations []string `json:"specializations" bson:"specializations"`   // ["IFRS", "US GAAP", "Tax Planning"]
	Certifications  []string `json:"certifications" bson:"certifications"`     // ["ISO 9001", "SOC 2"]
	Accreditations  []string `json:"accreditations" bson:"accreditations"`     // Professional accreditations

	// Company Description
	ShortDescription string   `json:"short_description" bson:"short_description"` // Brief tagline
	FullDescription  string   `json:"full_description" bson:"full_description"`   // Detailed company profile
	MissionStatement string   `json:"mission_statement" bson:"mission_statement"` // Company mission
	VisionStatement  string   `json:"vision_statement" bson:"vision_statement"`   // Company vision
	CoreValues       []string `json:"core_values" bson:"core_values"`             // ["Integrity", "Excellence"]

	// License/Registration
	CommercialLicense  string    `json:"commercial_license" bson:"commercial_license"`   // Business license #
	TaxIDNumber        string    `json:"tax_id_number" bson:"tax_id_number"`             // Tax identification
	RegistrationNumber string    `json:"registration_number" bson:"registration_number"` // Company registration
	LicenseIssueDate   time.Time `json:"license_issue_date" bson:"license_issue_date"`   // License issued
	LicenseExpiryDate  time.Time `json:"license_expiry_date" bson:"license_expiry_date"` // License expiry

	// Membership Information
	MembershipStartDate time.Time `json:"membership_start_date" bson:"membership_start_date"` // Joined LACPA
	MembershipStatus    string    `json:"membership_status" bson:"membership_status"`         // "Active", "Suspended", "Expired"
	MembershipTier      string    `json:"membership_tier" bson:"membership_tier"`             // "Gold", "Silver", "Bronze"
	RenewalDate         time.Time `json:"renewal_date" bson:"renewal_date"`                   // Next renewal due
	IsActive            bool      `json:"is_active" bson:"is_active"`                         // Active status
	DuesStatus          string    `json:"dues_status" bson:"dues_status"`                     // "Paid", "Pending", "Overdue"

	// Individual Members Associated with Firm
	AssociatedMemberIDs []primitive.ObjectID `json:"associated_member_ids" bson:"associated_member_ids"`               // References to IndividualMember
	PrimaryPartnerID    *primitive.ObjectID  `json:"primary_partner_id,omitempty" bson:"primary_partner_id,omitempty"` // Main partner reference

	// Sharing & Profile
	ProfileURL string `json:"profile_url" bson:"profile_url"` // Public firm profile URL
	QRCodeURL  string `json:"qr_code_url" bson:"qr_code_url"` // QR code for sharing

	// Privacy Settings
	ShowPhone         bool `json:"show_phone" bson:"show_phone"`                   // Show phone publicly
	ShowEmail         bool `json:"show_email" bson:"show_email"`                   // Show email publicly
	ShowWebsite       bool `json:"show_website" bson:"show_website"`               // Show website publicly
	ShowAddress       bool `json:"show_address" bson:"show_address"`               // Show address publicly
	ShowRevenue       bool `json:"show_revenue" bson:"show_revenue"`               // Show revenue publicly
	ShowEmployeeCount bool `json:"show_employee_count" bson:"show_employee_count"` // Show employee count

	// Search Optimization
	SearchTags []string `json:"search_tags" bson:"search_tags"` // ["deloitte", "audit", "big4"]

	// Additional Metadata
	ProfileViews      int       `json:"profile_views" bson:"profile_views"`           // Number of profile views
	EventsSponsored   int       `json:"events_sponsored" bson:"events_sponsored"`     // Count of events sponsored
	TrainingsProvided int       `json:"trainings_provided" bson:"trainings_provided"` // Training sessions provided
	LastUpdatedAt     time.Time `json:"last_updated_at" bson:"last_updated_at"`       // Last profile update

	// Financial Contribution
	AnnualContribution float64 `json:"annual_contribution" bson:"annual_contribution"` // Annual LACPA contribution
	TotalContributions float64 `json:"total_contributions" bson:"total_contributions"` // Lifetime contributions
	SponsorshipLevel   string  `json:"sponsorship_level" bson:"sponsorship_level"`     // "Platinum", "Gold", "Silver"
}

// FirmMetrics represents aggregate statistics for firm filtering
type FirmMetrics struct {
	TotalFirms       int `json:"total_firms"`        // Total count
	AuditFirmsCount  int `json:"audit_firms_count"`  // Count by type
	Big4Count        int `json:"big4_count"`         // Big 4 firms
	LargeFirmsCount  int `json:"large_firms_count"`  // Large firms
	MediumFirmsCount int `json:"medium_firms_count"` // Medium firms
	SmallFirmsCount  int `json:"small_firms_count"`  // Small firms
}

// Helper Methods

// IsExpiringSoon checks if license expires within 30 days
func (f *FirmMember) IsExpiringSoon() bool {
	if f.LicenseExpiryDate.IsZero() {
		return false
	}
	daysUntilExpiry := time.Until(f.LicenseExpiryDate).Hours() / 24
	return daysUntilExpiry > 0 && daysUntilExpiry <= 30
}

// IsDuesOverdue checks if membership dues are overdue
func (f *FirmMember) IsDuesOverdue() bool {
	return f.DuesStatus == "Overdue"
}

// IsRenewalDue checks if renewal is due within 30 days
func (f *FirmMember) IsRenewalDue() bool {
	if f.RenewalDate.IsZero() {
		return false
	}
	daysUntilRenewal := time.Until(f.RenewalDate).Hours() / 24
	return daysUntilRenewal > 0 && daysUntilRenewal <= 30
}

// GetYearsInBusiness calculates years since establishment
func (f *FirmMember) GetYearsInBusiness() int {
	if f.YearEstablished == 0 {
		return 0
	}
	return time.Now().Year() - f.YearEstablished
}

// IsBig4 checks if firm is one of the Big 4
func (f *FirmMember) IsBig4() bool {
	return f.FirmSize == "Big 4"
}

// HasWebsite checks if firm has a website
func (f *FirmMember) HasWebsite() bool {
	return f.Website != "" && f.ShowWebsite
}
