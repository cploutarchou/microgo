package mailer

import (
	"errors"
	"testing"
)

func TestMail_SendSMTPMessage(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
		To:          "you@there.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	err := mailer.SentSMTPMessage(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestMail_SendUsingChan(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
		To:          "you@there.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	mailer.Jobs <- msg
	res := <-mailer.Results
	if res.Error != nil {
		t.Error(errors.New("failed to send over channel"))
	}

	msg.To = "not_an_email_address"
	mailer.Jobs <- msg
	res = <-mailer.Results
	if res.Error == nil {
		t.Error(errors.New("no error received with invalid to address"))
	}
}

func TestMail_SendUsingAPI(t *testing.T) {
	msg := Message{
		To:          "you@there.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	mailer.API = "unknown"
	mailer.ApiKey = "abc123"
	mailer.ApiUrl = "https://www.fake.com"

	err := mailer.SendViaAPI(msg, "unknown")
	if err == nil {
		t.Error(err)
	}
	mailer.API = ""
	mailer.ApiKey = ""
	mailer.ApiUrl = ""
}

func TestMail_buildHTMLMessage(t *testing.T) {
	msg := Message{
		From:        "test@mymail.com",
		FromName:    "Christos Ploutarchou",
		To:          "you@mymail.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	_, err := mailer.createHTMLMessage(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestMail_buildPlainMessage(t *testing.T) {
	msg := Message{
		From:        "test@mymail.com",
		FromName:    "Christos Ploutarchou",
		To:          "you@mymail.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	_, err := mailer.createPlanMessage(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestMail_send(t *testing.T) {
	msg := Message{
		From:        "test@mymail.com",
		FromName:    "Christos Ploutarchou",
		To:          "you@mymail.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	err := mailer.Send(msg)
	if err != nil {
		t.Error(err)
	}

	mailer.API = "unknown"
	mailer.ApiKey = "abc123"
	mailer.ApiUrl = "https://www.fake.com"

	err = mailer.Send(msg)
	if err == nil {
		t.Error("did not not get an error when we should have")
	}

	mailer.API = ""
	mailer.ApiKey = ""
	mailer.ApiKey = ""
}

func TestMail_ChooseAPI(t *testing.T) {
	msg := Message{
		From:        "test@mymail.com",
		FromName:    "Christos Ploutarchou",
		To:          "you@mymail.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}
	mailer.API = "unknown"
	err := mailer.SelectAPI(msg)
	if err == nil {
		t.Error(err)
	}
}
