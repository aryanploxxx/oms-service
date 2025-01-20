package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"oms-service/models"
	"strconv"
)

func CreateOrder(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Create Orders Function called")

		orderToCheck := models.Order{}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &orderToCheck)
		fmt.Println(orderToCheck)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		//
		// LOGIC TO VALIDATE PRODUCT AND CONSUMER
		//

		checkError := make(chan error, 2)
		var validationErrors int

		go func() {
			resp, err := http.Get("http://localhost:8080/products/" + strconv.Itoa(orderToCheck.Product_id))
			if err != nil || resp.StatusCode == http.StatusNotFound {
				checkError <- fmt.Errorf("product ID not found")
				return
			}
			checkError <- nil
		}()

		go func() {
			resp, err := http.Get("http://localhost:8080/customers/" + strconv.Itoa(orderToCheck.Customer_id))
			if err != nil || resp.StatusCode == http.StatusNotFound {
				checkError <- fmt.Errorf("customer ID not found")
				return
			}
			checkError <- nil
		}()

		for i := 0; i < 2; i++ {
			err := <-checkError
			if err != nil {
				validationErrors++
			}
		}

		if validationErrors > 0 {
			http.Error(w, "Invalid product ID or customer ID", http.StatusNotFound)
			return
		}

		err = database.QueryRow("INSERT INTO orders (product_id, customer_id) VALUES ($1, $2) RETURNING order_id", orderToCheck.Product_id, orderToCheck.Customer_id).Scan(&orderToCheck.Order_id)
		if err != nil {
			http.Error(w, "Error saving order", http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(orderToCheck)
		if err != nil {
			http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)

	}
}
