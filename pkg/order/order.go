package order

import (
	"fmt"
	"github.com/LuanaFn/FDM-protocol/pkg/log"
	"github.com/lurifn/fdm-backend/configs"
	"github.com/lurifn/fdm-backend/pkg/email"
	"net/smtp"
)

// Create sends an email with the provided order according to the configurations loaded
func Create(order string) error {
	config, err := configs.Load()
	if err != nil {
		log.Error.Println("error loading configs: ", err)
		return err
	}
	log.Info.Printf("Order: %s\n", order)
	auth := email.LoginAuth(config.Email.From.Email, config.Email.From.Password)
	message := fmt.Sprintf(
		"To: %s\r\nSubject: New order!\r\n\r\n%s\r\n",
		config.Email.To,
		order,
	)

	return smtp.SendMail(
		config.Email.From.SMTP+":"+config.Email.From.Port,
		auth,
		config.Email.From.Email,
		config.Email.To,
		[]byte(message),
	)
}
