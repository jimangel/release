/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mail

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

// GoogleGroup is a simple google group representation
type GoogleGroup string

const (
	KubernetesAnnounceGoogleGroup     GoogleGroup = "kubernetes-announce"
	KubernetesDevGoogleGroup          GoogleGroup = "dev"
	KubernetesAnnounceTestGoogleGroup GoogleGroup = "kubernetes-announce-test"
)

type Sender struct {
	apiKey     string
	sendClient SendClient
	apiClient  APIClient
	sender     *mail.Email
	recipients []*mail.Email
}

func NewSender(apiKey string) *Sender {
	return &Sender{
		apiKey:     apiKey,
		sendClient: sendgrid.NewSendClient(apiKey),
		apiClient:  &sendgridAPIClient{},
	}
}

// SetSendClient can be used to set the sendgrid sender client
func (s *Sender) SetSendClient(client SendClient) {
	s.sendClient = client
}

// SetSendClient can be used to set the sendgrid API client
func (s *Sender) SetAPIClient(client APIClient) {
	s.apiClient = client
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . SendClient
//go:generate /usr/bin/env bash -c "cat ../../hack/boilerplate/boilerplate.generatego.txt mailfakes/fake_send_client.go > mailfakes/_fake_send_client.go && mv mailfakes/_fake_send_client.go mailfakes/fake_send_client.go"
type SendClient interface {
	Send(*mail.SGMailV3) (*rest.Response, error)
}

//counterfeiter:generate . APIClient
//go:generate /usr/bin/env bash -c "cat ../../hack/boilerplate/boilerplate.generatego.txt mailfakes/fake_apiclient.go > mailfakes/_fake_apiclient.go && mv mailfakes/_fake_apiclient.go mailfakes/fake_apiclient.go"
type APIClient interface {
	API(rest.Request) (*rest.Response, error)
}

type sendgridAPIClient struct{}

func (s *sendgridAPIClient) API(request rest.Request) (*rest.Response, error) {
	return sendgrid.API(request)
}

func (s *Sender) Send(body, subject string) error {
	html := mail.NewContent("text/html", body)

	p := mail.NewPersonalization()
	p.AddTos(s.recipients...)

	msg := mail.NewV3Mail().
		SetFrom(s.sender).
		AddContent(html).
		AddPersonalizations(p)
	msg.Subject = subject

	logrus.WithField("message", msg).Trace("Message prepared")

	res, err := s.sendClient.Send(msg)
	if err != nil {
		return err
	}
	if res == nil {
		return &SendError{code: -1, resBody: "<empty API response>"}
	}
	if c := res.StatusCode; c < 200 || c >= 300 {
		return &SendError{code: res.StatusCode, resBody: res.Body, resHeaders: fmt.Sprintf("%#v", res.Headers)}
	}

	logrus.Debug("Mail successfully sent")
	return nil
}

type SendError struct {
	code       int
	resBody    string
	resHeaders string
}

func (s *SendError) Error() string {
	return fmt.Sprintf("got code %d while sending: Body: %q, Header: %q", s.code, s.resBody, s.resHeaders)
}

func (s *Sender) SetDefaultSender() error {
	// Retrieve the mail
	request := sendgrid.GetRequest(s.apiKey, "/v3/user/email", "")
	response, err := s.apiClient.API(request)
	if err != nil {
		return fmt.Errorf("getting user email: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to get users email: %s", response.Body)
	}
	type email struct {
		Email string `json:"email"`
	}
	emailResponse := &email{}
	if err := json.Unmarshal([]byte(response.Body), emailResponse); err != nil {
		return fmt.Errorf("decoding JSON response: %w", err)
	}
	logrus.Infof("Using sender address: %s", emailResponse.Email)

	// Retrieve first and last name
	request = sendgrid.GetRequest(s.apiKey, "/v3/user/profile", "")
	response, err = s.apiClient.API(request)
	if err != nil {
		return fmt.Errorf("getting user profile: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to get users profile: %s", response.Body)
	}
	type profile struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	pr := profile{}
	if err := json.Unmarshal([]byte(response.Body), &pr); err != nil {
		return fmt.Errorf("decoding JSON response: %w", err)
	}

	name := fmt.Sprintf("%s %s", pr.FirstName, pr.LastName)
	logrus.Infof("Using sender name: %s", name)

	s.sender = mail.NewEmail(name, emailResponse.Email)
	return nil
}

func (s *Sender) SetSender(name, email string) error {
	if email == "" {
		return errors.New("email must not be empty")
	}
	s.sender = mail.NewEmail(name, email)
	logrus.WithField("sender", s.sender).Debugf("Sender set")
	return nil
}

func (s *Sender) SetRecipients(recipientArgs ...string) error {
	l := len(recipientArgs)

	if l%2 != 0 {
		return errors.New("must be called with alternating recipient's names and email addresses")
	}

	recipients := make([]*mail.Email, l/2)

	for i := range recipients {
		name := recipientArgs[i*2]
		email := recipientArgs[i*2+1]
		if email == "" {
			return errors.New("email must not be empty")
		}
		recipients[i] = mail.NewEmail(name, email)
	}

	s.recipients = recipients
	logrus.WithField("recipients", s.sender).Debugf("Recipients set")

	return nil
}

// SetGoogleGroupRecipient can be used to set multiple Google Groups as recipient
func (s *Sender) SetGoogleGroupRecipients(groups ...GoogleGroup) error {
	args := []string{}
	for _, group := range groups {
		if group == "dev" {
			args = append(args, string(group), fmt.Sprintf("%s@kubernetes.io", group))
		} else {
			args = append(args, string(group), fmt.Sprintf("%s@googlegroups.com", group))
		}
	}
	return s.SetRecipients(args...)
}

// GetRecipients can be used to get the recipients
func (s *Sender) GetRecipients() []*mail.Email {
	return s.recipients
}
