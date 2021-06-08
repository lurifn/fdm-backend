package main

import (
	"github.com/LuanaFn/FDM-protocol/pkg/log"
	"github.com/joho/godotenv"
	"github.com/lurifn/fdm-backend/pkg/order"
	"github.com/lurifn/fdm-backend/pkg/product"
	"github.com/lurifn/fdm-backend/pkg/repository/mongodb"
	"net/http"
	"os"
)

func main() {
	log.Info.Println("Initializing app...")

	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Warning.Println("Error trying to load environment variables from .env file:", err)
	}

	db := mongodb.MongoDB{
		URI:        os.Getenv("MONGO_URI"),
		DB:         os.Getenv("MONGO_DB_NAME"),
		Collection: os.Getenv("MONGO_ORDERS_COLLECTION_NAME"),
		Username:   os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
		Password:   os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
	}

	order.HandleHTTPRequests(db)
	product.HandleHTTPRequests(db)

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
