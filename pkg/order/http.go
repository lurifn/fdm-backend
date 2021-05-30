package order

import (
	"fmt"
	"github.com/LuanaFn/FDM-protocol/pkg/log"
	"github.com/lurifn/fdm-backend/pkg/repository"
	"io"
	"net/http"
)

func create(w http.ResponseWriter, r *http.Request, repo repository.Repository) {
	// Message
	message, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error.Println(err)
		handleError("Error processing order", w)

		return
	}
	err = r.Body.Close()
	if err != nil {
		log.Warning.Println("Error closing request body: ", err)
	}

	err = Create(string(message), repo)
	if err != nil {
		log.Error.Println(err)
		handleError("Error sending order", w)

		return
	}

	log.Info.Print("Order Sent Successfully!")
}

func handleError(msg string, w http.ResponseWriter) {
	errorResponse := fmt.Sprintf("{\"error\":\"%s\"}", msg)

	w.WriteHeader(http.StatusInternalServerError)

	_, wErr := w.Write([]byte(errorResponse))
	if wErr != nil {
		log.Error.Panicf("error communicating error to consumer: \nError:%+v\nOriginal error:%s", wErr, msg)
	}
}

// HandleHTTPRequests register the handlers for the APIs in this package
// To expose the APIs you must run http.ListenAndServe after calling this.
func HandleHTTPRequests(repo repository.Repository) {
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			create(w, r, repo)
		default:
			log.Error.Print("error: invalid request ", r.Method)
			handleError("invalid request", w)
			w.WriteHeader(http.StatusBadRequest)
		}
	})
}
