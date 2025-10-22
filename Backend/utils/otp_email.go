package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"time"
)

// EmailConfig holds SMTP configuration
type EmailConfig struct {
	SMTPHost    string
	SMTPPort    string
	SenderEmail string
	SenderPass  string
}

// OTPEmail holds OTP data
type OTPEmail struct {
	Email     string
	OTP       string
	ExpiresAt time.Time
}

// GetEmailConfig returns the email configuration
func GetEmailConfig() EmailConfig {
	return EmailConfig{
		SMTPHost:    "smtp.gmail.com",
		SMTPPort:    "587",
		SenderEmail: "sleimana181@gmail.com",
		SenderPass:  "mqdy abgx fjzx bfms", // App password (spaces will be removed)
	}
}

// SendOTPEmail sends an OTP email to the specified recipient
// Generates a new OTP internally
func SendOTPEmail(recipientEmail string) (*OTPEmail, error) {
	otp, err := GenerateOTP()
	if err != nil {
		log.Printf("Error generating OTP: %v", err)
		return nil, fmt.Errorf("failed to generate OTP: %w", err)
	}
	return SendOTPEmailWithCode(recipientEmail, "", otp)
}

// SendOTPEmailWithName sends an OTP email to the specified recipient with personalized name
// Generates a new OTP internally
func SendOTPEmailWithName(recipientEmail, recipientName string) (*OTPEmail, error) {
	otp, err := GenerateOTP()
	if err != nil {
		log.Printf("Error generating OTP: %v", err)
		return nil, fmt.Errorf("failed to generate OTP: %w", err)
	}
	return SendOTPEmailWithCode(recipientEmail, recipientName, otp)
}

// SendOTPEmailWithCode sends a specific OTP code via email
// Use this when you've already generated an OTP and stored it in database
func SendOTPEmailWithCode(recipientEmail, recipientName, otp string) (*OTPEmail, error) {
	expiresAt := time.Now().Add(5 * time.Minute)

	// Log OTP to console for testing
	log.Printf("========================================")
	log.Printf("Generated OTP for %s: %s", recipientEmail, otp)
	log.Printf("Expires at: %s", expiresAt.Format("2006-01-02 15:04:05"))
	log.Printf("========================================")

	// Get email config
	config := GetEmailConfig()

	// Remove spaces from app password
	password := ""
	for _, char := range config.SenderPass {
		if char != ' ' {
			password += string(char)
		}
	}

	// Setup authentication
	auth := smtp.PlainAuth("", config.SenderEmail, password, config.SMTPHost)

	// Compose email with HTML template
	subject := "Your OTP Code - LACPA Verification"
	htmlBody := OTPEmailTemplate(otp, recipientName)

	// Email headers and body with HTML content
	message := []byte(
		"From: LACPA <" + config.SenderEmail + ">\r\n" +
			"To: " + recipientEmail + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			htmlBody + "\r\n",
	)

	// Send email
	smtpAddr := config.SMTPHost + ":" + config.SMTPPort
	err := smtp.SendMail(smtpAddr, auth, config.SenderEmail, []string{recipientEmail}, message)
	if err != nil {
		log.Printf("Failed to send OTP email to %s: %v", recipientEmail, err)
		return nil, fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("OTP email sent successfully to %s", recipientEmail)

	return &OTPEmail{
		Email:     recipientEmail,
		OTP:       otp,
		ExpiresAt: expiresAt,
	}, nil
}

// VerifyOTP verifies if the provided OTP matches and hasn't expired
func VerifyOTP(storedOTP *OTPEmail, providedOTP string) bool {
	if storedOTP == nil {
		log.Printf("No OTP found")
		return false
	}

	// Check if OTP has expired
	if time.Now().After(storedOTP.ExpiresAt) {
		log.Printf("OTP has expired")
		return false
	}

	// Check if OTP matches
	if storedOTP.OTP != providedOTP {
		log.Printf("OTP mismatch: expected %s, got %s", storedOTP.OTP, providedOTP)
		return false
	}

	log.Printf("OTP verified successfully")
	return true
}
