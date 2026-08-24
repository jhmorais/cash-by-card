package email

import (
	"context"
	"strings"
	"testing"
)

func TestSmtpSenderRejectsCRLFInRecipient(t *testing.T) {
	s := &smtpSender{host: "smtp.example.com", port: "587", user: "u", password: "p", from: "u@example.com"}
	err := s.SendPasswordResetEmail(context.Background(), "a@x.com\r\nBcc: victim@evil.com", "http://link")
	if err == nil || !strings.Contains(err.Error(), "invalid recipient") {
		t.Fatalf("esperado erro de recipient invalido, got '%v'", err)
	}
}

func TestLogSenderRejectsCRLFInRecipient(t *testing.T) {
	s := &logSender{}
	err := s.SendPasswordResetEmail(context.Background(), "a@x.com\r\nBcc: x@evil.com", "http://link")
	if err == nil || !strings.Contains(err.Error(), "invalid recipient") {
		t.Fatalf("esperado erro de recipient invalido, got '%v'", err)
	}
}

func TestLogSenderLogsAndSucceeds(t *testing.T) {
	s := &logSender{}
	if err := s.SendPasswordResetEmail(context.Background(), "ok@x.com", "http://link"); err != nil {
		t.Fatalf("esperado nil, got '%v'", err)
	}
}
