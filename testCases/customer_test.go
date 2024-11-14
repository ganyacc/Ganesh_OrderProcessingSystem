package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/ganyacc/Ganesh_OrderProcessingSystem/entities"
	"github.com/stretchr/testify/assert"
)

// Test case for GetAllCustomers endpoint
func TestGetAllCustomers(t *testing.T) {
	// Send GET request to the local server
	resp, err := http.Get("http://localhost:8080/api/customers")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Assert the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// You can also check the response body if needed
	var customers []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&customers); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Example assertion: Check if there is at least one customer
	assert.True(t, len(customers) > 0)
}

// Test case for GetCustomerByID endpoint
func TestGetCustomerByID(t *testing.T) {
	customerID := "10ac6f2c-18ae-46da-9cca-4f36c84ce342" // Replace with an actual customer ID in your database

	// Send GET request to the local server
	resp, err := http.Get("http://localhost:8080/api/customers/" + customerID)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Assert the response status code
	if resp.StatusCode == http.StatusOK {
		// Parse and check the customer data if status is OK
		var customer map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&customer); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		// Example assertion: Check if the customer name exists in the response
		assert.Contains(t, customer["name"], "ganesh")
	} else {
		// Assert not found if status is 404
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	}
}

// Test case for CreateOrder endpoint
func TestCreateOrder(t *testing.T) {
	// Prepare a sample order payload
	orderPayload := entities.OrderRequest{
		CustomerID: "10ac6f2c-18ae-46da-9cca-4f36c84ce381", // Replace with a valid customer ID
		ProductIDs: []string{"11ac5f2d-18ea-46ad-9cca-3f36c84ce123", "22ac5f2d-18ea-46ad-9cca-3f36c84ce103", "33ac5f2d-18ea-46ad-9cca-3f36c84ce103"},
	}

	orderPayloadJSON, _ := json.Marshal(orderPayload)

	// Send POST request to create order
	resp, err := http.Post("http://localhost:8080/api/orders", "application/json", bytes.NewBuffer(orderPayloadJSON))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Assert the response status code
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Print the raw response body to check its format
	t.Logf("Response Body: %s", responseBody)

	if resp.StatusCode >= 400 {
		var errorResponse string
		if err := json.Unmarshal(responseBody, &errorResponse); err != nil {
			t.Fatalf("Failed to parse error response: %v", err)
		}
		t.Fatalf("API Error: %v", errorResponse)
	}

	// Parse the response body (the created order)
	var order *entities.Order
	if err := json.Unmarshal(responseBody, &order); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Example assertion: Check if the order total price exists in the response
	assert.Equal(t, 77.35, order.TotalPrice)
}
