package mailer

import (
	"bytes"
	"fmt"
	"github.com/ainsleyclark/go-mail/drivers"
	"github.com/ainsleyclark/go-mail/mail"
	"github.com/vanng822/go-premailer/premailer"
	defaultMail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

type Mailer struct {
	Domain      string
	Templates   string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
	Jobs        chan Message
	Results     chan Result
	API         string
	ApiKey      string
	ApiUrl      string
}
type Message struct {
	From           string
	FromName       string
	To             string
	BCC            []string
	CC             []string
	Subject        string
	Template       string
	TemplateFormat TemplateFormat
	Attachments    []string
	Data           interface{}
}

type Result struct {
	Success bool
	Error   error
}
type TemplateFormat string

const (
	HTMLTemplateFormat      TemplateFormat = "html"
	PlainTextTemplateFormat TemplateFormat = "plain/text"
)

func (m *Mailer) ListenForMessage() {
	for {
		msg := <-m.Jobs
		err := m.Send(msg)
		if err != nil {
			m.Results <- Result{Success: false, Error: err}
		} else {
			m.Results <- Result{Success: true, Error: nil}
		}
	}
}

func (m *Mailer) Send(msg Message) error {
	if len(m.API) > 0 && len(m.ApiKey) > 0 && len(m.ApiUrl) > 0 && m.API != "smtp" {
		return m.SelectAPI(msg)
	}
	return m.SentSMTPMessage(msg)
}

func (m *Mailer) SentSMTPMessage(msg Message) error {

	formattedMessage, err := m.createHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := m.createPlanMessage(msg)
	if err != nil {
		return err
	}

	svr := defaultMail.NewSMTPClient()
	svr.Host = m.Host
	svr.Port = m.Port
	svr.Username = m.Username
	svr.Password = m.Password
	svr.KeepAlive = false
	svr.ConnectTimeout = 15 * time.Second
	svr.SendTimeout = 15 * time.Second
	svr.Encryption = m.getEncryption(m.Encryption)
	smtpClient, err := svr.Connect()
	if err != nil {
		return err
	}

	email := defaultMail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)
	if msg.TemplateFormat == HTMLTemplateFormat {
		email.SetBody(defaultMail.TextHTML, formattedMessage)
	}
	if msg.TemplateFormat == PlainTextTemplateFormat {
		email.AddAlternative(defaultMail.TextPlain, plainMessage)
	}
	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}
	if len(msg.CC) > 0 {
		for _, x := range msg.Attachments {
			email.AddCc(x)
		}
	}
	if len(msg.BCC) > 0 {
		for _, x := range msg.Attachments {
			email.AddBcc(x)
		}
	}
	err = email.Send(smtpClient)
	if err != nil {
		return err
	}
	return nil
}

// createHTMLMessage Build HTML Message using html template.
func (m *Mailer) createHTMLMessage(msg Message) (string, error) {
	renderTemplate := fmt.Sprintf("%s/%s.html.tmpl", m.Templates, msg.Template)
	t, err := template.New("email-html").ParseFiles(renderTemplate)
	if err != nil {
		return "", err
	}
	var tml bytes.Buffer
	if err = t.ExecuteTemplate(&tml, "body", msg.Data); err != nil {
		return "", err
	}
	fmtMSG := tml.String()
	fmtMSG, err = m.inlineCSS(fmtMSG)
	return fmtMSG, nil
}

// createPlanMessage : Build HTML Message using plan template.
func (m *Mailer) createPlanMessage(msg Message) (string, error) {
	renderTemplate := fmt.Sprintf("%s/%s.plain.tmpl", m.Templates, msg.Template)
	t, err := template.New("email-html").ParseFiles(renderTemplate)
	if err != nil {
		return "", err
	}
	var tml bytes.Buffer
	if err = t.ExecuteTemplate(&tml, "body", msg.Data); err != nil {
		return "", err
	}
	plainMSG := tml.String()
	return plainMSG, nil
}

func (m *Mailer) getEncryption(str string) defaultMail.Encryption {
	switch str {
	case "tls":
		return defaultMail.EncryptionSTARTTLS
	case "ssl":
		return defaultMail.EncryptionSSL
	case "none":
		return defaultMail.EncryptionNone
	default:
		return defaultMail.EncryptionSTARTTLS
	}
}

func (m *Mailer) inlineCSS(str string) (string, error) {
	opt := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}
	prem, err := premailer.NewPremailerFromString(str, &opt)
	if err != nil {
		return "", err
	}
	newHtml, err := prem.Transform()
	if err != nil {
		return "", err
	}
	return newHtml, nil
}

func (m *Mailer) SelectAPI(msg Message) error {
	switch m.API {
	case "mailgun", "sparkpost", "sendgrid", "postal", "postmark":
		return m.SendViaAPI(msg, m.API)
	default:
		return fmt.Errorf("Not supported api: %s ", m.API)
	}

}

func (m *Mailer) SendViaAPI(msg Message, transport string) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	cfg := mail.Config{
		URL:         m.ApiUrl,
		APIKey:      m.ApiKey,
		Domain:      m.Domain,
		FromAddress: msg.From,
		FromName:    msg.FromName,
	}

	formattedMessage, err := m.createHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := m.createPlanMessage(msg)

	if err != nil {
		return err
	}

	tx := &mail.Transmission{
		Recipients: []string{msg.To},
		Subject:    msg.Subject,
	}
	if msg.TemplateFormat == PlainTextTemplateFormat {
		tx.PlainText = plainMessage
	}
	if msg.TemplateFormat == HTMLTemplateFormat {
		tx.HTML = formattedMessage
	}
	_mailer, err := m.SelectAPIDriver(transport, cfg)
	if err != nil {
		return err
	}
	// add attachments
	err = m.addAPIAttachments(msg, tx)
	if err != nil {
		return err
	}

	_, err = _mailer.Send(tx)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mailer) SelectAPIDriver(transport string, config mail.Config) (mail.Mailer, error) {
	switch transport {
	case "sparkpost":
		return drivers.NewSparkPost(config)
	case "mailgun":
		return drivers.NewMailgun(config)
	case "postal":
		return drivers.NewPostal(config)
	case "postmark":
		return drivers.NewPostmark(config)
	case "sendgrid":
		return drivers.NewSendGrid(config)
	default:
		return nil, fmt.Errorf("No valid transport specified : %s ", transport)
	}
}

func (m *Mailer) addAPIAttachments(msg Message, tx *mail.Transmission) error {
	if len(msg.Attachments) > 0 {
		var attachments []mail.Attachment

		for _, x := range msg.Attachments {
			var attach mail.Attachment
			content, err := os.ReadFile(x)
			if err != nil {
				return err
			}

			fileName := filepath.Base(x)
			attach.Bytes = content
			attach.Filename = fileName
			attachments = append(attachments, attach)
		}

		tx.Attachments = attachments
	}

	return nil
}
