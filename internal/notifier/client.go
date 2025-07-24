package notifier

import (
	"log"

	"github.com/go-resty/resty/v2"
)

type NotifierInterface interface {
	SendWarning(employeeAbbreviation string, message string) error
}

type Notifier struct {
	client *resty.Client
	url    string
}

func NewNotifier(baseURL string) *Notifier {
	return &Notifier{
		client: resty.New(),
		url:    baseURL,
	}
}

func (n *Notifier) SendWarning(abbr string, message string) error {
	body := map[string]string{
		"level":                "warning",
		"employeeAbbreviation": abbr,
		"message":              message,
	}
	_, err := n.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(n.url + "/api/notify")

	log.Println("SentWarning for " + abbr)

	return err
}
