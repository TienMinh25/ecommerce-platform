package adaptor

import (
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/pkg"
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
	return nil
}
