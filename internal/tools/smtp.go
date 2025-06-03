package tools

import (
	"github.com/audricimanuel/errorutils"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/config"
	"net/smtp"
	"strings"
)

type (
	EmailType string

	SMTPClient interface {
		SendEmail(mimeType string, emailType EmailType, subject, message string, receivers []string, cc []string) error
	}

	SMTPClientImpl struct {
		Host           string
		Port           string
		User           string
		Password       string
		UseTLS         bool
		CSEmailAddress string
	}
)

const (
	EMAIL_TYPE_OTP EmailType = "otp"
)

func NewSMTPClient(cfg config.Config) SMTPClient {
	smtpConfig := cfg.SMTPConfig
	return &SMTPClientImpl{
		Host:           smtpConfig.Host,
		Port:           smtpConfig.Port,
		User:           smtpConfig.User,
		Password:       smtpConfig.Password,
		CSEmailAddress: smtpConfig.CSEmailAddress,
	}
}

func (s *SMTPClientImpl) SendEmail(mimeType string, emailType EmailType, subject, message string, receivers []string, cc []string) error {
	typeSender := map[EmailType]string{
		EMAIL_TYPE_OTP: s.CSEmailAddress,
		"default":      s.CSEmailAddress,
	}

	sender := typeSender["default"]
	if email, ok := typeSender[emailType]; ok {
		sender = email
	}

	mime := ""
	if mimeType == "html" {
		mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	}

	body := "From: " + sender + "\n" +
		"To: " + strings.Join(receivers, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n" +
		mime + "\n" +
		message

	smtpAuth := smtp.PlainAuth("", s.User, s.Password, s.Host)

	if err := smtp.SendMail(s.Host+":"+s.Port, smtpAuth, sender, append(receivers, cc...), []byte(body)); err != nil {
		return errorutils.ErrorInternalServer.CustomMessage(err.Error())
	}

	return nil
}
