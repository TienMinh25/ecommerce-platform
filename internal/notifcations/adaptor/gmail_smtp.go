package adaptor

import (
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"net/smtp"
	"strings"
)

type gmailSmtpAdapter struct {
	tracer pkg.Tracer
	env    *env.EnvManager
}

func NewGmailSmtpAdapter(tracer pkg.Tracer, env *env.EnvManager) IGmailSmtpAdapter {
	return &gmailSmtpAdapter{
		tracer: tracer,
		env:    env,
	}
}

func (g *gmailSmtpAdapter) SendMail(data SendMailRequest) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "3000"
	senderEmail := "dominh12223@gmail.com"
	senderPassword := "uhed ysci skqp dnar"

	subject := "Xác nhận đăng ký tài khoản"
	body := fmt.Sprintf("Xin chào %s,\n\nMã xác nhận của bạn là: %s\nVui lòng nhập mã này để hoàn tất đăng ký.\n\nCảm ơn!",
		data.FullName, data.OTP)

	headers := []string{
		"From: " + senderEmail,
		"To: " + data.To,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"utf-8\"",
		"",
		body,
	}
	message := strings.Join(headers, "\r\n")

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{data.To}, []byte(message))
	if err != nil {
		return fmt.Errorf("lỗi gửi mail: %v", err)
	}

	fmt.Println("Email xác nhận đã gửi đến", data.To)
	return nil
}
