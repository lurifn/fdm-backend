package email

import (
	"errors"
	"fmt"
	"github.com/LuanaFn/FDM-protocol/pkg/log"
	"net/smtp"
)

type Email struct {
	NoReplyEmail    string
	NoReplyPassword string
	NoReplySMTP     string
	NoReplyPort     string
	BusinessEmail   string
}

var errEmailLogin = errors.New("error trying to login to email server")

type loginAuth struct {
	username, password string
}

// LoginAuth returns a smtp.Auth implementation of type LOGIN to be used in emails.
func login(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("%w: unknown fromServer '%x'", errEmailLogin, fromServer)
		}
	}

	return nil, nil
}

func (e Email) Save(message string) error {
	auth := login(e.NoReplyEmail, e.NoReplyPassword)
	msg := fmt.Sprintf(
		"To: %s\r\nSubject: New order!\r\n\r\n%s\r\n",
		e.BusinessEmail,
		message,
	)

	log.Debug.Printf(`Sending message to "%s": "%s"`, e.BusinessEmail, message)

	return smtp.SendMail(
		e.NoReplySMTP+":"+e.NoReplyPort,
		auth,
		e.NoReplyEmail,
		[]string{e.BusinessEmail},
		[]byte(msg),
	)
}
