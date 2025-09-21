package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/fboydc/smartsplit-main/models"
	"github.com/gin-gonic/gin"
)

// DummyHandlers contains handlers for dummy/test data
type DummyHandlers struct{}

// NewDummyHandlers creates a new DummyHandlers instance
func NewDummyHandlers() *DummyHandlers {
	return &DummyHandlers{}
}

// GetDummyTransactionsHandler returns dummy transaction data from a JSON file
func (h *DummyHandlers) GetDummyTransactionsHandler(c *gin.Context) {
	// Open the JSON file
	file, err := os.Open("data/test_transactions(first period payment).json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open JSON file"})
		return
	}
	defer file.Close()

	var response models.GetDummyTransactionsResponse
	if err := json.NewDecoder(file).Decode(&response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON file"})
		return
	}

	// Send the parsed data as a JSON response
	c.JSON(http.StatusOK, gin.H{"latest_transactions": response.LatestTransactions})
}
