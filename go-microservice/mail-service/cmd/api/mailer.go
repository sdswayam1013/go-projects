package main

import (
	"bytes"
	"html/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func (m *Mail) SendSMTPMessage(msg Message) error {

	// If sender email is empty, use the default one
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	// If sender name is empty, use the default one
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	// Data that will be injected into the email template
	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	// Build the HTML version of the email
	formattedMessage, err := m.buildHTMLMessage(msg)

	if err != nil {
		return err
	}
	plainMessage, err := m.buildPlainTextMessage(msg)

	if err != nil {
		return err
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {

	// Path to the HTML template
	templateToRender := "./templates/mail.html.gohtml"

	// Parse the template file
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	// Buffer to hold the rendered template
	var tpl bytes.Buffer

	// Execute the template using msg.DataMap
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	// Convert buffer to string
	formattedMessage := tpl.String()

	// Inline CSS into the HTML
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	// Return the final HTML
	return formattedMessage, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {

	// Path to the HTML template
	templateToRender := "./templates/mail.plain.gohtml"

	// Parse the template file
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	// Buffer to hold the rendered template
	var tpl bytes.Buffer

	// Execute the template using msg.DataMap
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	// Convert buffer to string
	plainMessage := tpl.String()

	// Return the final HTML
	return plainMessage, nil
}

func (m *Mail) inlineCSS(s string) (string, error) {

	// Configure CSS inlining options
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	// Create a new premailer instance from the HTML string
	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	// Convert CSS into inline styles
	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	// Return the transformed HTML
	return html, nil
}

func (m *Mail) getEncryption(encryption string) mail.Encryption {

	switch encryption {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
