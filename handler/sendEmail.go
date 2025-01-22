package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type EmailWithTemplateRequestBody struct {
	ToAddr   string            `json:"to_addr"`
	Subject  string            `json:"subject"`
	Template string            `json:"template"`
	Vars     map[string]string `json:"vars"`
}

type Customer struct {
	Customer_id    int    `json:"Customer_id"`
	Customer_name  string `json:"Customer_name"`
	Customer_email string `json:"Customer_email"`
}

func CustomerDetails(customerID int) (*Customer, error) {
	resp, err := http.Get("http://localhost:8080/customers/" + strconv.Itoa(customerID))
	if err != nil {
		return nil, fmt.Errorf("error getting customer details: %v", err)
	}
	defer resp.Body.Close()

	var customer Customer
	err = json.NewDecoder(resp.Body).Decode(&customer)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &customer, nil
}

func SendEmails(order *OrderBulkRecieve) error {
	customer, err := CustomerDetails(order.CustomerID)
	if err != nil {
		return err
	}

	emailReq := EmailWithTemplateRequestBody{
		ToAddr:   customer.Customer_email,
		Subject:  "HTML Mail Test",
		Template: "helloEmail",
		Vars: map[string]string{
			"Name":    customer.Customer_name,
			"Product": strconv.Itoa(order.ProductID),
		},
	}

	payloadBytes, err := json.Marshal(emailReq)
	if err != nil {
		return fmt.Errorf("error marshaling email request: %v", err)
	}

	resp, err := http.Post("http://localhost:9001/html_email_template", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
