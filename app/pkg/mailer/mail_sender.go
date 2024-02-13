package mailer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/Budgetin-Project/user-service/app/domain/model"
	"gopkg.in/mail.v2"
)

const EmailVerificationTemplatePath = "./app/pkg/mailer/email_verification_template.html"

// EmailVerificationData holds data for email verification template in 'email_verification_template.html'
type EmailVerificationData struct {
	User             string
	VerificationLink string
	SupportEmail     string
	CompanyName      string
	Expiration       int
}

// RenderEmailVerificationTemplate renders the email verification template
func RenderEmailVerificationTemplate(data *EmailVerificationData) (string, error) {
	templateFile, err := os.ReadFile(EmailVerificationTemplatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	temp, err := template.New("email_verification").Parse(string(templateFile))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := temp.Execute(&body, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return body.String(), nil
}

// ComposeEmailMessage composes the email message
func ComposeEmailMessage(emailTo string, body string) *mail.Message {
	m := mail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":    {m.FormatAddress(os.Getenv("SMTP_SENDER_EMAIL"), os.Getenv("SMTP_SENDER_ALIAS"))},
		"To":      {emailTo},
		"Subject": {"Email Verification"},
	})
	m.SetBody("text/html", body)
	return m
}

// SendEmail sends the email
func SendEmail(emailTo string, body string) error {
	// Settings for SMTP server
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %w", err)
	}
	dial := mail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_SENDER_EMAIL"), os.Getenv("SMTP_SENDER_PASS"))

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production, this should be set to false.
	dial.TLSConfig = &tls.Config{InsecureSkipVerify: os.Getenv("APP_ENV") == "local"}

	// Send the email
	if err := dial.DialAndSend(ComposeEmailMessage(emailTo, body)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendEmailVerification sends an email with the specified mail data
func SendEmailVerification(emailTo string, userName string, verificationToken string, expiredAt time.Time) error {
	// TODO: Later 'VerificationLink', 'SupportEmail', 'CompanyName', and 'Expiration' will be retrieved from the configuration (database)
	data := EmailVerificationData{
		User:             userName,
		VerificationLink: fmt.Sprintf("http://localhost:8080/email-verification/%s", verificationToken),
		SupportEmail:     "Andresuryana17@gmail.com",
		CompanyName:      "Budgetin",
		Expiration:       model.TokenExpDuration,
	}

	body, err := RenderEmailVerificationTemplate(&data)
	if err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	if err := SendEmail(emailTo, body); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}
