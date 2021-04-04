package main

import (
	"github.com/LuanaFn/FDM-protocol/pkg/log"
	"github.com/lurifn/fdm-backend/configs"
	"github.com/lurifn/fdm-backend/pkg/order"
	"net/http"
)

func main() {
	err := configs.Config.Load()
	if err != nil {
		log.Error.Panic("error loading configs: ", err)
	}

	order.HandleRequests()

	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Error.Panic("error opening order service: ", err)
	}
}
