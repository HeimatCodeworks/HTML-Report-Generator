package email_sender

import (
	"fmt"
	"net/smtp"
	"regexp"
	"strings"
)

const (
	gmailSMTPHost = "smtp.gmail.com"
	gmailSMTPPort = "587"
)

func minifyHTML(html string) string {

	re := regexp.MustCompile(``)
	minified := re.ReplaceAllString(html, "")

	re = regexp.MustCompile(`\n|\t`)
	minified = re.ReplaceAllString(minified, "")

	re = regexp.MustCompile(`>\s+<`)
	minified = re.ReplaceAllString(minified, "><")

	return strings.TrimSpace(minified)
}

func SendEmail(from, password string, to, cc []string, subject, body string) error {
	auth := smtp.PlainAuth("", from, password, gmailSMTPHost)

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = strings.Join(to, ", ")
	if len(cc) > 0 {
		headers["Cc"] = strings.Join(cc, ", ")
	}
	headers["Subject"] = subject
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")

	minifiedBody := minifyHTML(body)
	msg.WriteString(minifiedBody)

	allRecipients := append(to, cc...)

	err := smtp.SendMail(gmailSMTPHost+":"+gmailSMTPPort, auth, from, allRecipients, []byte(msg.String()))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
