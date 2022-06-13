package mailer

import (
	"bytes"
	"fmt"
	"html/template"
)

type Mailer struct {
	Domain      string
	Templates   string
	Host        string
	Port        string
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
	formattedTemplate := tml.String()
	return formattedTemplate, nil
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
