package adaptor

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/internal/notifications/transport/grpc/proto/notification_proto_gen"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"gopkg.in/gomail.v2"
	"html/template"
)

//go:embed template/*
var templateFolder embed.FS

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
	var templateFile []byte
	var subject string
	var err error

	if data.Purpose == notification_proto_gen.PurposeOTP_PASSWORD_RESET {
		templateFile, err = templateFolder.ReadFile("template/otp-forgot-password.html")
		subject = "Đặt lại mật khẩu của bạn"
	} else {
		templateFile, err = templateFolder.ReadFile("template/otp-verify-email.html")
		subject = "Xác nhận email của bạn"
	}

	if err != nil {
		return err
	}

	tmpl, err := template.New("email").Parse(string(templateFile))

	if err != nil {
		return err
	}

	dataMail := struct {
		Name       string
		OTPCode    string
		ExpireTime int
	}{
		Name:       data.FullName,
		OTPCode:    data.OTP,
		ExpireTime: g.env.OTPVerifyEmailTimeout,
	}

	var bodyMail bytes.Buffer
	if err = tmpl.Execute(&bodyMail, dataMail); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("To", data.To)
	m.SetHeader("From", g.env.Mail.MailFrom)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", bodyMail.String())

	// send mail
	d := gomail.NewDialer(g.env.Mail.MailHost, 587, g.env.Mail.MailUser, g.env.Mail.MailPassword)

	if err = d.DialAndSend(m); err != nil {
		return err
	}

	fmt.Println("Email xác nhận đã gửi đến", data.To)
	return nil
}
