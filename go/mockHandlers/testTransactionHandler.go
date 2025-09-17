package mockHandlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Transaction struct {
	TransactionID           string   `json:"transaction_id"`
	AccountID               string   `json:"account_id"`
	Amount                  float64  `json:"amount"`
	IsoCurrencyCode         string   `json:"iso_currency_code"`
	Date                    string   `json:"date"`
	AuthorizedDate          string   `json:"authorized_date"`
	Name                    string   `json:"name"`
	MerchantName            *string  `json:"merchant_name"`
	PaymentChannel          string   `json:"payment_channel"`
	Pending                 bool     `json:"pending"`
	TransactionType         string   `json:"transaction_type"`
	Category                []string `json:"category"`
	CategoryID              string   `json:"category_id"`
	PersonalFinanceCategory struct {
		Primary         string `json:"primary"`
		Detailed        string `json:"detailed"`
		ConfidenceLevel string `json:"confidence_level"`
	} `json:"personal_finance_category"`
	Location struct {
		City    string `json:"city,omitempty"`
		Region  string `json:"region,omitempty"`
		Country string `json:"country"`
	} `json:"location"`
}

type TransactionTestResponse struct {
	Transactions []Transaction `json:"transactions"`
}

// TestTransactionHandler is a mock handler for testing transaction processing
func TestTransactionHandler(c *gin.Context) {
	data, err := os.ReadFile("mockResponses/30-06-2025/transactions.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var res TransactionTestResponse
	if err := json.Unmarshal(data, &res); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process the transactions in res.Transactions
	// ...\

	c.JSON(http.StatusOK, gin.H{"message": "Transactions processed successfully", "data": res.Transactions})

}
