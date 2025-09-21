package api

import (
	"github.com/google/uuid"
)

// Allocation represents budget allocation information
type Allocation struct {
	Id                    uuid.UUID `json:"id"`
	AllocationType        string    `json:"allocation_type"`
	AllocationDescription string    `json:"allocation_description"`
	AllocationFactor      float64   `json:"allocation_factor"`
}

// Category represents a category in the database
type Category struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// Expense represents an expense record
type Expense struct {
	Id             uuid.UUID `json:"id"`
	Description    string    `json:"description"`
	Amount         float64   `json:"amount"`
	Category       string    `json:"category"`
	AllocationType string    `json:"allocation_type"`
}

// Income represents an income record
type Income struct {
	Id        uuid.UUID `json:"id"`
	Amount    float64   `json:"amount"`
	Frequency string    `json:"frequency"`
}

// LoginRequest represents the login payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// GetBudgetResponse represents the response structure for budget data
type GetBudgetResponse struct {
	Expenses    []Expense    `json:"expenses"`
	Incomes     []Income     `json:"incomes"`
	Allocations []Allocation `json:"allocations"`
}

// Transaction represents a financial transaction from Plaid
type Transaction struct {
	TransactionID           string   `json:"transaction_id"`
	AccountID               string   `json:"account_id"`
	Amount                  float64  `json:"amount"`
	ISOCurrencyCode         string   `json:"iso_currency_code"`
	Date                    string   `json:"date"`
	AuthorizedDate          string   `json:"authorized_date"`
	Name                    string   `json:"name"`
	MerchantName            string   `json:"merchant_name"`
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
		City    string `json:"city"`
		Region  string `json:"region"`
		Country string `json:"country"`
	} `json:"location"`
}

// GetDummyTransactionsResponse represents the response structure for dummy transactions
type GetDummyTransactionsResponse struct {
	LatestTransactions []Transaction `json:"latest_transactions"`
}

// APIResponse represents a generic API response structure
type APIResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Message    string `json:"message"`
	Token      string `json:"token"`
	PlaidToken string `json:"plaidToken"`
	UserID     string `json:"user_id"`
	Username   string `json:"username"`
}
