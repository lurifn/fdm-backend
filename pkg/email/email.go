package email

import (
	"errors"
	"fmt"
	"net/smtp"
)

var errEmailLogin = errors.New("error trying to login to email server")

type loginAuth struct {
	username, password string
}

// LoginAuth returns a smtp.Auth implementation of type LOGIN to be used in emails.
func LoginAuth(username, password string) smtp.Auth {
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
