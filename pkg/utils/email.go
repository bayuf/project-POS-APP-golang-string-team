package utils

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
)

func SendEmail(payload dto.Email, config *Configuration) error {
	log.Println(payload)
	from := config.Email.From
	password := config.Email.Password // 16 digit App Password
	to := []string{payload.Email}

	// SMTP Gmail Server Config
	smtpHost := config.Email.SMTPHost
	smtpPort := config.Email.SMTPPort

	subject := fmt.Sprintf("Subject: %s\r\n", payload.Subject)
	contentType := "MIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n"

	// combine all the part string
	message := []byte(subject + contentType + payload.Body)

	// Autentikasi
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending Process
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}
	return nil
}
