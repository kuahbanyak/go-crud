package email
import (
	"fmt"
	"net/smtp"
)
type SMTPClient struct {
	host     string
	port     string
	username string
	password string
}
type EmailService interface {
	SendEmail(to, subject, body string) error
}
func NewSMTPClient(host, port, username, password string) *SMTPClient {
	return &SMTPClient{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}
func (s *SMTPClient) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, body))
	return smtp.SendMail(s.host+":"+s.port, auth, s.username, []string{to}, msg)
}

