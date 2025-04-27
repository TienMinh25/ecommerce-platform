package adaptor

type IGmailSmtpAdapter interface {
	SendMail(data SendMailRequest) error
}
