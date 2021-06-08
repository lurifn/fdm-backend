package product

import (
	"encoding/json"
	"fmt"
	"github.com/LuanaFn/FDM-protocol/pkg/log"
	"github.com/lurifn/fdm-backend/pkg/repository"
	"io"
	"net/http"
)

func create(w http.ResponseWriter, r *http.Request, repo repository.Repository) {
	// Message
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error.Println(err)
		handleError("Error processing order", w, http.StatusBadRequest)

		return
	}
	err = r.Body.Close()
	if err != nil {
		log.Warning.Println("Error closing request body:", err)
	}

	subject := Product{}
	err = json.Unmarshal(body, &subject)
	if err != nil {
		log.Error.Println(err)
		handleError("Error parsing message, please check if body is valid", w, http.StatusBadRequest)

		return
	}

	err = subject.Save(repo)
	if err != nil {
		log.Error.Println(err)
		handleError("Error sending order", w, http.StatusInternalServerError)

		return
	}

	log.Info.Print("Product saved Successfully!")
	subjectBytes, err := json.Marshal(subject)
	if err != nil {
		log.Error.Println(err)
		handleError("Error writing response", w, http.StatusCreated)

		return
	}

	_, err = w.Write(subjectBytes)
}

func handleError(msg string, w http.ResponseWriter, serverStatus int) {
	errorResponse := fmt.Sprintf("{\"error\":\"%s\"}", msg)

	w.WriteHeader(serverStatus)

	_, wErr := w.Write([]byte(errorResponse))
	if wErr != nil {
		log.Error.Panicf("error communicating error to consumer: \nError:%+v\nOriginal error:%s", wErr, msg)
	}
}

// HandleHTTPRequests register the handlers for the APIs in this package
// To expose the APIs you must run http.ListenAndServe after calling this.
func HandleHTTPRequests(repo repository.Repository) {
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			create(w, r, repo)
		default:
			log.Error.Print("error: invalid request ", r.Method)
			handleError("invalid request", w, http.StatusBadRequest)
		}
	})
}
