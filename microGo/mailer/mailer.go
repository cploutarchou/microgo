package mailer

import (
	"bytes"
	"fmt"
	"github.com/vanng822/go-premailer/premailer"
	"github.com/xhit/go-simple-mail/v2"
	"html/template"
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
	ApiUr       string
}
type Message struct {
	From        string
	FromName    string
	To          string
	BCC         []string
	CC          []string
	Subject     string
	Template    string
	Attachments []string
	Data        interface{}
}

type Result struct {
	Success bool
	Error   error
}

func (m *Mailer) ListenForMessage() {
	for {
		msg := <-m.Jobs
		err := m.Send(msg)
		if err != nil {
			m.Results <- Result{
				Success: false,
				Error:   err,
			}
		} else {
			m.Results <- Result{
				Success: true,
				Error:   nil,
			}
		}
	}
}

func (m *Mailer) Send(msg Message) error {
	// TODO: SMTP OR API Server support.
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

	svr := mail.NewSMTPClient()
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

	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)
	email.SetBody(mail.TextHTML, formattedMessage)
	email.AddAlternative(mail.TextPlain, plainMessage)
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

//createHTMLMessage Build HTML Message using html template.
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

//createPlanMessage : Build HTML Message using plan template.
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

func (m *Mailer) getEncryption(str string) mail.Encryption {
	switch str {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSL
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
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
