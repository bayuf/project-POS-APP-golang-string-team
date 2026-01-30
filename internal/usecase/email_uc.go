package usecase

import (
	"strings"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"go.uber.org/zap"
)

type EmailService struct {
	logger *zap.Logger
	config *utils.Configuration
}

func NewEmailService(logger *zap.Logger, config *utils.Configuration) *EmailService {
	return &EmailService{
		logger: logger,
		config: config,
	}
}

func (uc *EmailService) SendEmail(payload dto.Email) error {
	var (
		body    string
		subject string
	)
	if payload.Type == "otp" {
		subject = "POS Golang: Kode Verifikasi Keamanan"
		htmlBody := `<div style="font-family: Arial, sans-serif; max-width: 600px; margin: auto; border: 1px solid #e0e0e0; padding: 20px;">
    <h2 style="color: #333; text-align: center;">Kode Verifikasi</h2>
    <p>Halo, {{username}}</p>
    <p>Kami menerima permintaan untuk akses ke akun Anda. Gunakan kode OTP di bawah ini untuk melanjutkan:</p>
    <div style="background-color: #f4f4f4; padding: 20px; text-align: center; font-size: 32px; font-weight: bold; letter-spacing: 5px; color: #2c3e50; margin: 20px 0;">
        {{code}}
    </div>
    <p style="color: #555; font-size: 14px;">Kode ini hanya berlaku selama <b>5 menit</b>. Jangan memberikan kode ini kepada siapapun.</p>
    <hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
    <p style="font-size: 12px; color: #888;">Jika Anda tidak merasa melakukan permintaan ini, abaikan email ini atau hubungi tim dukungan kami.</p>
    </div>`

		// replacer
		body = strings.NewReplacer("{{code}}", payload.Code, "{{username}}", payload.Username).Replace(htmlBody)

	}

	if payload.Type == "password" {
		subject = "POS Golang: Password Sementara"
		htmlBody := `
		<div style="font-family: Arial, sans-serif; max-width: 600px; margin: auto; border: 1px solid #e0e0e0; padding: 20px;">
    <h2 style="color: #333;">Halo, {{username}}</h2>
    <p>Akun Anda telah berhasil dibuat. Berikut adalah kredensial login sementara Anda:</p>
    <table style="width: 100%; background-color: #f9f9f9; padding: 15px; border-radius: 8px;">
        <tr>
            <td style="color: #888;">Email:</td>
            <td style="font-weight: bold;"> {{email}} </td>
        </tr>
        <tr>
            <td style="color: #888;">Password Sementara:</td>
            <td style="font-weight: bold; color: #e74c3c;">{{password}}</td>
        </tr>
    </table>
    <p style="margin-top: 20px;">Demi keamanan, kami sangat menyarankan Anda untuk segera mengganti password ini setelah login pertama kali.</p>

    <hr style="border: 0; border-top: 1px solid #eee;">
    <p style="font-size: 12px; color: #888;">Terima kasih telah bergabung bersama kami.</p>
    </div>`

		body = strings.NewReplacer(
			"{{username}}", payload.Username,
			"{{email}}", payload.Email,
			"{{password}}", payload.Password,
		).Replace(htmlBody)
	}

	// mapping
	finalpayload := dto.Email{
		Username: payload.Username,
		Email:    payload.Email,
		Subject:  subject,
		Body:     body,
		Code:     payload.Code,
		Password: payload.Password,
	}
	err := utils.SendEmail(finalpayload, uc.config)
	if err != nil {
		uc.logger.Error("Failed to send email", zap.Error(err))
		return err
	}

	return nil
}
