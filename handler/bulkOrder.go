package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	// "oms-service/handler"
)

type OrderBulkRecieve struct {
	ProductID  int `json:"Product_id"`
	CustomerID int `json:"Customer_id"`
}

var UniqueLogger *log.Logger

var bulkOrders []OrderBulkRecieve = []OrderBulkRecieve{
	{2, 1},
	{6, 7},
	{2, 2},
	{3, 2},
	{4, 3},
	{9, 4},
}

func InitializeRabbitMQ(database2 *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// bulkOrders = []byte(`[
		//     {
		//         "Product_id": 6,
		//         "Customer_id": 7
		//     },
		//     {
		//         "Product_id": 1,
		//         "Customer_id": 1
		//     },
		//     {
		//         "Product_id": 2,
		//         "Customer_id": 2
		//     }
		// ]`)

		go SendDataToMQ()
		go ReceiveFromMQ()
	}
}

func processMessage(body []byte) {
	var orders []OrderBulkRecieve
	err := json.Unmarshal(body, &orders)
	if err != nil {
		log.Printf("Failed to unmarshal JSON: %v", err)
		return
	}

	var workers []string = []string{"worker1", "worker2"}

	var tasksChann chan OrderBulkRecieve = make(chan OrderBulkRecieve, len(orders))
	var workersChann chan string = make(chan string, len(workers))

	fmt.Println("Worker Pool Started")

	var wg sync.WaitGroup

	go func() {
		for _, order := range orders {
			tasksChann <- order
		}
	}()

	go func() {
		for _, worker := range workers {
			workersChann <- worker
		}
	}()

	// for _, order := range orders {
	for i := 0; i < len(orders); i++ {
		order := <-tasksChann
		log.Printf("Processing Order -> ProductID: %d, CustomerID: %d \n", order.ProductID, order.CustomerID)
		worker := <-workersChann
		wg.Add(1)
		go processingSingleOrder(&order, worker, &wg, workersChann)
	}
}

func processingSingleOrder(order *OrderBulkRecieve, worker string, wg *sync.WaitGroup, workersChann chan string) {
	defer wg.Done()

	fmt.Println("Executing task", order, "by worker", worker)

	logFile, _ := os.OpenFile("bulkOrderFailures.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	UniqueLogger = log.New(logFile, "Log: ", log.Ldate|log.Ltime|log.Lshortfile)

	payload := map[string]int{
		"product_id":  order.ProductID,
		"customer_id": order.CustomerID,
	}

	newJSON, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		UniqueLogger.Printf("Request failed for Product id %v and Customer id %v", order.ProductID, order.CustomerID)
		return
	}

	res, err := http.Post("http://localhost:9000/order", "application/json", bytes.NewBuffer(newJSON))
	if err != nil {
		log.Printf("Failed to make POST request: %v", err)
		UniqueLogger.Printf("Request failed for Product id %v and Customer id %v", order.ProductID, order.CustomerID)
		return
	}

	defer res.Body.Close()

	var neww OrderBulkRecieve
	content, err := io.ReadAll(res.Body)
	errr := json.Unmarshal(content, &neww)
	if errr != nil {
		log.Printf("Failed to unmarshal response: %v", err)
		UniqueLogger.Printf("Request failed for Product id %v and Customer id %v", order.ProductID, order.CustomerID)
		return
	}

	fmt.Println(order, "completed by worker", worker)

	// order -> product id, customer id
	SendEmails(order)

	workersChann <- worker
	fmt.Printf("Response content: %+v\n", neww)
}
