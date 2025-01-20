package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"oms-service/handler"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var database2 *sql.DB

func init() {
	dbURL := "host=localhost port=5432 user=postgres password=Pyari@123 sslmode=disable dbname=orderDB"
	fmt.Println("Database URL: ", dbURL)

	var err error
	database2, err = sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		panic(err)
	} else {
		fmt.Println("Database connected successfully")
	}

	_, _ = database2.Exec("CREATE TABLE IF NOT EXISTS orders (order_id SERIAL PRIMARY KEY, product_id TEXT, customer_id TEXT)")
}

func main() {

	defer database2.Close()

	// handler.InitializeRabbitMQ(database2)

	router := mux.NewRouter()
	router.HandleFunc("/order", handler.CreateOrder(database2)).Methods("POST")
	router.HandleFunc("/bulkorder", handler.InitializeRabbitMQ(database2)).Methods("POST")
	http.ListenAndServe(":9000", router)

}
