package main

import (
	"github.com/LuanaFn/FDM-protocol/pkg/log"
	"github.com/joho/godotenv"
	"github.com/lurifn/fdm-backend/pkg/order"
	"github.com/lurifn/fdm-backend/pkg/repository/email"
	"net/http"
	"os"
)

func main() {
	log.Info.Println("Initializing app...")

	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Warning.Println("Error trying to load environment variables from .env file:", err)
	}

	order.HandleHTTPRequests(email.Email{
		NoReplyEmail:    os.Getenv("NOREPLY_EMAIL_ADDRESS"),
		NoReplyPassword: os.Getenv("NOREPLY_EMAIL_PASSWORD"),
		NoReplySMTP:     os.Getenv("NOREPLY_EMAIL_SMTP"),
		NoReplyPort:     os.Getenv("NOREPLY_EMAIL_SMTP_PORT"),
		BusinessEmail:   os.Getenv("BUSINESS_EMAIL_ADDRESS"),
	})

	c := make(chan int)
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Error.Fatal(err)
		}
		c <- 1
	}()

	log.Debug.Println("Listening on port 8080")
	log.Debug.Print(<-c)
}
