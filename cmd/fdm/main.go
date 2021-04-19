package main

import (
	"github.com/LuanaFn/FDM-protocol/pkg/log"
	"github.com/lurifn/fdm-backend/pkg/order"
	"net/http"
)

func main() {
	log.Info.Println("Initializing app...")
	order.HandleHTTPRequests()

	c := make(chan int)
	go func() {
		err := http.ListenAndServe(":8080", nil)
		log.Error.Fatal(err)
		c <- 1
	}()

	log.Debug.Println("Listening on port 8080")
	log.Debug.Print(<- c)
}
