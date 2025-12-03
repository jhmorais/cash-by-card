package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

type SMTPMailer struct {
	Host     string // e.g. "smtp.gmail.com"
	Port     int    // e.g. 587
	Username string
	Password string
	From     string
}

func NewSMTPMailer(host string, port int, username, password, from string) *SMTPMailer {
	return &SMTPMailer{Host: host, Port: port, Username: username, Password: password, From: from}
}

func (m *SMTPMailer) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)

	headers := map[string]string{
		"From":         m.From,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/plain; charset=\"utf-8\"",
	}
	var msg strings.Builder
	for k, v := range headers {
		fmt.Fprintf(&msg, "%s: %s\r\n", k, v)
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)
	return smtp.SendMail(addr, auth, m.From, []string{to}, []byte(msg.String()))
}
