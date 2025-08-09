package notification

import (
    "net/smtp"
    "os"
    "fmt"
)

func SendEmail(to, subject, body string) error {
    host := os.Getenv("SMTP_HOST")
    port := os.Getenv("SMTP_PORT")
    user := os.Getenv("SMTP_USER")
    pass := os.Getenv("SMTP_PASS")
    if host == "" || user == "" {
        return nil // noop in dev
    }
    auth := smtp.PlainAuth("", user, pass, host)
    msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))
    return smtp.SendMail(host+":"+port, auth, user, []string{to}, msg)
}
