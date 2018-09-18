package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"github.com/gobuffalo/packr"
)

// SMTPConfig holds the host information and credentials for sending e-mails using SMTP.
type SMTPConfig struct {
	Host     string `toml:"host"`
	Port     uint16 `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

// SendEmail sends the e-mail body to the given receiver over SMTP.
func (cfg *SMTPConfig) SendEmail(subject, body string, to ...string) error {

	for i := range to {
		to[i] = strings.TrimSpace(to[i])
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	fromString := fmt.Sprintf("Goirate <%s>", cfg.Username)
	commaSeparatedTo := strings.Join(to, ",")

	msg := fmt.Sprintf("From: %v\nTo: %s\nSubject: %s\n%s\n\n%s\n", fromString, commaSeparatedTo, subject, mime, body)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host),
		cfg.Username,
		to,
		[]byte(msg),
	)

	return err
}

// LoadTorrentTemplate generates the torrent notification e-mail by loading the template
// and populating it with the given data.
func LoadTorrentTemplate(data interface{}) (string, error) {

	box := packr.NewBox("../../templates")

	html, err := box.MustString("torrent.html")

	if err != nil {
		return "", err
	}

	templ, err := template.New("tml").Parse(html)

	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	if err = templ.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
