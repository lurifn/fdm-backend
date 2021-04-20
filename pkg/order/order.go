package order

import (
	"fmt"
	"github.com/LuanaFn/FDM-protocol/pkg/log"
	"github.com/lurifn/fdm-backend/pkg/email"
	"net/smtp"
)

type EmailConfig struct {
	NoReplyEmail    string
	NoReplyPassword string
	NoReplySMTP     string
	NoReplyPort     string
	BusinessEmail   string
}

// Create sends an email with the provided order according to the configurations loaded
func Create(order string, config EmailConfig) error {
	log.Info.Printf("Order: %s\n", order)
	auth := email.LoginAuth(config.NoReplyEmail, config.NoReplyPassword)
	message := fmt.Sprintf(
		"To: %s\r\nSubject: New order!\r\n\r\n%s\r\n",
		config.BusinessEmail,
		order,
	)

	return smtp.SendMail(
		config.NoReplySMTP+":"+config.NoReplyPort,
		auth,
		config.NoReplyEmail,
		[]string{config.BusinessEmail},
		[]byte(message),
	)
}
