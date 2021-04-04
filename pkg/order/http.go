package order

import (
	"fmt"
	"github.com/LuanaFn/FDM-protocol/pkg/log"
	"github.com/lurifn/fdm-backend/configs"
	"github.com/lurifn/fdm-backend/pkg/email"
	"net/http"
	"net/smtp"
)

func create(w http.ResponseWriter, r *http.Request) {
	// Message
	var message []byte
	_, err := r.Body.Read(message)
	if err != nil {
		log.Error.Println(err)
		handleError("Error building email message", w)

		return
	}

	auth := email.LoginAuth(configs.Config.Email.From.Email, configs.Config.Email.From.Password)
	err = smtp.SendMail(
		configs.Config.Email.From.SMTP+":"+configs.Config.Email.From.Port,
		auth,
		configs.Config.Email.From.Email,
		configs.Config.Email.To,
		message,
	)

	if err != nil {
		log.Error.Println(err)
		handleError("Error sending email", w)

		return
	}

	log.Info.Print("Email Sent Successfully!")
}

func handleError(msg string, w http.ResponseWriter) {
	errorResponse := fmt.Sprintf("{\"error\":\"%s\"}", msg)

	w.WriteHeader(http.StatusInternalServerError)

	_, wErr := w.Write([]byte(errorResponse))
	if wErr != nil {
		log.Error.Panicf("error communicating error to consumer: \nError:%+v\nOriginal error:%s", wErr, msg)
	}
}

/**
HandleRequests register the handlers for the APIs in this package
To expose the APIs you must run http.ListenAndServe after calling this.
*/
func HandleRequests() {
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			create(w, r)
		default:
			log.Error.Print("error: invalid request ", r.Method)
			handleError("invalid request", w)
			w.WriteHeader(http.StatusBadRequest)
		}
	})
}
