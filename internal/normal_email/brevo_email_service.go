package normal_email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type EmailService interface {
	SendNotificationEmail(toEmail, toName, subject, htmlContent string) error
}

type brevoService struct {
	apiKey    string
	fromEmail string
	fromName  string
	client    *http.Client
}

func NewBrevoService() EmailService {
	apiKey := os.Getenv("BREVO_API_KEY")
	fromEmail := os.Getenv("EMAIL_FROM")
	fromName := os.Getenv("EMAIL_FROM_NAME")

	if apiKey == "" {
		log.Fatal("BREVO_API_KEY environment variable is required")
	}

	if fromEmail == "" {
		fromEmail = "notifications@findix.railway.app"
	}

	if fromName == "" {
		fromName = "Findix FMIIS System"
	}

	log.Printf("📧 Brevo service initialized with From: %s <%s>", fromName, fromEmail)

	return &brevoService{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		fromName:  fromName,
		client:    &http.Client{Timeout: 30 * time.Second},
	}
}

// Brevo API request structures
type brevoEmailRequest struct {
	Sender      brevoSender `json:"sender"`
	To          []brevoTo   `json:"to"`
	Subject     string      `json:"subject"`
	HTMLContent string      `json:"htmlContent"`
}

type brevoSender struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type brevoTo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Brevo API response structure
type brevoResponse struct {
	MessageID string `json:"messageId"`
}

func (s *brevoService) SendNotificationEmail(toEmail, toName, subject, htmlContent string) error {
	// Validate inputs
	if toEmail == "" {
		return fmt.Errorf("recipient email cannot be empty")
	}

	log.Printf("📤 Attempting to send email to: %s (%s)", toEmail, toName)

	// Prepare email request
	emailRequest := brevoEmailRequest{
		Sender: brevoSender{
			Email: s.fromEmail,
			Name:  s.fromName,
		},
		To: []brevoTo{
			{
				Email: toEmail,
				Name:  toName,
			},
		},
		Subject:     subject,
		HTMLContent: htmlContent,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(emailRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal email request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", s.apiKey)
	req.Header.Set("Accept", "application/json")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email request: %v", err)
	}
	defer resp.Body.Close()

	// Parse response
	var response brevoResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("Warning: Could not decode Brevo response: %v", err)
	}

	// Check status code
	if resp.StatusCode >= 400 {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)

		errorMsg := "unknown error"
		if message, ok := errorResponse["message"].(string); ok {
			errorMsg = message
		}

		return fmt.Errorf("brevo API error: %s (status: %d)", errorMsg, resp.StatusCode)
	}

	log.Printf("✅ Email sent successfully to %s. Message ID: %s", toEmail, response.MessageID)
	return nil
}
