package models

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IndividualMember represents a single LACPA member in the system
type IndividualMember struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`

	// Basic Information (State 0 - Card Front)
	LacpaID    string `json:"lacpa_id" bson:"lacpa_id"`       // Unique LACPA ID: "3666"
	FirstName  string `json:"first_name" bson:"first_name"`   // "Boushra"
	MiddleName string `json:"middle_name" bson:"middle_name"` // "El"
	LastName   string `json:"last_name" bson:"last_name"`     // "Obeid"
	FullName   string `json:"full_name" bson:"full_name"`     // Computed or stored: "Boushra El Obeid"

	AvatarURL  string `json:"avatar_url" bson:"avatar_url"`   // Path to profile image
	MemberType string `json:"member_type" bson:"member_type"` // "Apprentices", "Practicing", "Non-Practicing", "Retired"
	BadgeEmoji string `json:"badge_emoji" bson:"badge_emoji"` // Visual badge: "🎓"
	BadgeColor string `json:"badge_color" bson:"badge_color"` // Tailwind color: "emerald-500"

	// Contact Information (State 1 - Contact Details)
	Phone       string `json:"phone" bson:"phone"`               // "+961 01 123 456"
	Email       string `json:"email" bson:"email"`               // "boushra@gmail.com"
	LinkedInURL string `json:"linkedin_url" bson:"linkedin_url"` // "linkedin.com/in/boushra..."
	Firm        string `json:"firm" bson:"firm"`                 // "Nabil El Haj" - Associated firm/company

	// Address Components
	FullAddress string `json:"full_address" bson:"full_address"` // "Mount Lebanon - Metn - Zalqa - Al-Samrani"
	Governorate string `json:"governorate" bson:"governorate"`   // "Mount Lebanon"
	District    string `json:"district" bson:"district"`         // "Metn"
	City        string `json:"city" bson:"city"`                 // "Zalqa"
	Area        string `json:"area" bson:"area"`                 // "Al-Samrani"
	Country     string `json:"country" bson:"country"`           // "Lebanon"

	// Biography/Professional Summary (State 2 - Bio)
	Biography           string `json:"biography" bson:"biography"`                       // Full professional bio
	ProfessionalSummary string `json:"professional_summary" bson:"professional_summary"` // Short summary

	// Professional Information
	Title             string   `json:"title" bson:"title"`                             // "CPA", "Senior Auditor"
	Position          string   `json:"position" bson:"position"`                       // Current job position
	Specializations   []string `json:"specializations" bson:"specializations"`         // ["Auditing", "Tax"]
	Services          []string `json:"services" bson:"services"`                       // Services offered
	YearsOfExperience int      `json:"years_of_experience" bson:"years_of_experience"` // Years in profession

	// License/Certification
	LicenseNumber     string    `json:"license_number" bson:"license_number"`           // Professional license #
	LicenseIssueDate  time.Time `json:"license_issue_date" bson:"license_issue_date"`   // License issued date
	LicenseExpiryDate time.Time `json:"license_expiry_date" bson:"license_expiry_date"` // License expiry date

	// Membership Information
	MembershipStartDate time.Time `json:"membership_start_date" bson:"membership_start_date"` // Joined LACPA
	MembershipStatus    string    `json:"membership_status" bson:"membership_status"`         // "Active", "Suspended", "Expired"
	RenewalDate         time.Time `json:"renewal_date" bson:"renewal_date"`                   // Next renewal due
	IsActive            bool      `json:"is_active" bson:"is_active"`                         // Active status
	DuesStatus          string    `json:"dues_status" bson:"dues_status"`                     // "Paid", "Pending", "Overdue"

	// Council Information (Reference-based, not embedded)
	CurrentCouncilPositionID *primitive.ObjectID `json:"current_council_position_id,omitempty" bson:"current_council_position_id,omitempty"` // Reference to active CouncilPosition
	CouncilPosition          string              `json:"council_position" bson:"council_position"`                                           // Current position (cached for quick access)
	IsCouncilMember          bool                `json:"is_council_member" bson:"is_council_member"`                                         // Quick flag for filtering

	// Sharing & Profile
	ProfileURL string `json:"profile_url" bson:"profile_url"` // Public profile URL
	QRCodeURL  string `json:"qr_code_url" bson:"qr_code_url"` // QR code for sharing

	// Privacy Settings
	ShowPhone    bool `json:"show_phone" bson:"show_phone"`       // Show phone publicly
	ShowEmail    bool `json:"show_email" bson:"show_email"`       // Show email publicly
	ShowLinkedIn bool `json:"show_linkedin" bson:"show_linkedin"` // Show LinkedIn publicly
	ShowAddress  bool `json:"show_address" bson:"show_address"`   // Show address publicly

	// Search Optimization
	SearchTags []string `json:"search_tags" bson:"search_tags"` // ["boushra", "obeid", "apprentice", "auditor"]

	// Additional Metadata
	LastLoginAt      time.Time `json:"last_login_at" bson:"last_login_at"`         // Last login timestamp
	ProfileViews     int       `json:"profile_views" bson:"profile_views"`         // Number of profile views
	EventsAttended   int       `json:"events_attended" bson:"events_attended"`     // Count of events attended
	CPECredits       int       `json:"cpe_credits" bson:"cpe_credits"`             // Continuing education credits
	CommitteesServed []string  `json:"committees_served" bson:"committees_served"` // ["Audit", "Ethics"]
}

// MemberMetrics represents aggregate statistics for member filtering
type MemberMetrics struct {
	TotalMembers       int `json:"total_members"`        // Total count
	ApprenticesCount   int `json:"apprentices_count"`    // Count by type
	PracticingCount    int `json:"practicing_count"`     // Count by type
	NonPracticingCount int `json:"non_practicing_count"` // Count by type
	RetiredCount       int `json:"retired_count"`        // Count by type
}

// Helper Methods

// GetFullName constructs full name from components
func (m *IndividualMember) GetFullName() string {
	if m.FullName != "" {
		return m.FullName
	}
	name := m.FirstName
	if m.MiddleName != "" {
		name += " " + m.MiddleName
	}
	if m.LastName != "" {
		name += " " + m.LastName
	}
	return name
}

// GetDisplayAddress formats address for display
func (m *IndividualMember) GetDisplayAddress() string {
	if m.FullAddress != "" {
		return m.FullAddress
	}
	parts := []string{}
	if m.Governorate != "" {
		parts = append(parts, m.Governorate)
	}
	if m.District != "" {
		parts = append(parts, m.District)
	}
	if m.City != "" {
		parts = append(parts, m.City)
	}
	if m.Area != "" {
		parts = append(parts, m.Area)
	}
	return strings.Join(parts, " - ")
}

// IsExpiringSoon checks if license expires within 30 days
func (m *IndividualMember) IsExpiringSoon() bool {
	if m.LicenseExpiryDate.IsZero() {
		return false
	}
	daysUntilExpiry := time.Until(m.LicenseExpiryDate).Hours() / 24
	return daysUntilExpiry > 0 && daysUntilExpiry <= 30
}

// IsDuesOverdue checks if membership dues are overdue
func (m *IndividualMember) IsDuesOverdue() bool {
	return m.DuesStatus == "Overdue"
}

// IsRenewalDue checks if renewal is due within 30 days
func (m *IndividualMember) IsRenewalDue() bool {
	if m.RenewalDate.IsZero() {
		return false
	}
	daysUntilRenewal := time.Until(m.RenewalDate).Hours() / 24
	return daysUntilRenewal > 0 && daysUntilRenewal <= 30
}

// HasCouncilPosition checks if member has any council position
func (m *IndividualMember) HasCouncilPosition() bool {
	return m.IsCouncilMember && m.CouncilPosition != "" && m.CouncilPosition != "Non-Council Member"
}

// IsLeader checks if member is president or vice president
func (m *IndividualMember) IsLeader() bool {
	return m.CouncilPosition == "President" || m.CouncilPosition == "Vice President"
}
