package domain

import "time"

type TransactionRaw struct {
	TransactionID   string // plaid_transaction_id
	UserID          string
	ItemID          string
	AccountID       string
	Name            string
	MerchantName    string
	Amount          float64 // normalize convention; outflow positive
	ISOCurrencyCode string
	Date            time.Time
	AuthorizedDate  *time.Time
	Pending         bool
	PaymentChannel  string
	TransactionType string
	Category        []string
	CategoryID      string
	PFCPrimary      string
	PFCDetailed     string
	PFCConfidence   string
	LocationCity    string
	LocationRegion  string
	LocationCountry string
	RawJSON         []byte // full plaid txn for future use
	RemovedAt       *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
