package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/fboydc/smartsplit-main/client"
	"github.com/fboydc/smartsplit-main/config"
	"github.com/gin-gonic/gin"
	plaid "github.com/plaid/plaid-go/v31/plaid"
)

// PlaidHandlers contains all Plaid-related HTTP handlers
type PlaidHandlers struct {
	client      *client.PlaidClient
	config      *config.Config
	db          *sql.DB
	authHandler *AuthHandlers
}

// Global variables for token storage - in production, use a secure persistent store
var (
	accessToken     string
	userToken       string
	itemID          string
	paymentID       string
	authorizationID string
	accountID       string
)

// NewPlaidHandlers creates a new PlaidHandlers instance
func NewPlaidHandlers(plaidClient *client.PlaidClient, cfg *config.Config, db *sql.DB, authHandler *AuthHandlers) *PlaidHandlers {
	return &PlaidHandlers{
		client:      plaidClient,
		config:      cfg,
		db:          db,
		authHandler: authHandler,
	}
}

// InfoHandler returns basic info about the Plaid connection
func (h *PlaidHandlers) InfoHandler(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"item_id":      itemID,
		"access_token": accessToken,
		"products":     strings.Split(h.config.PlaidProducts, ","),
	})
}

// GetAccessTokenHandler exchanges a public token for an access token
func (h *PlaidHandlers) GetAccessTokenHandler(c *gin.Context) {
	publicToken := c.PostForm("public_token")
	user := c.PostForm("user")

	exchangeResp, err := h.client.ExchangePublicToken(publicToken)
	if err != nil {
		renderError(c, err)
		return
	}

	accessToken = exchangeResp.GetAccessToken()
	itemID = exchangeResp.GetItemId()

	ok, err := h.authHandler.SaveAccessToken(accessToken, user)
	if err != nil {
		renderError(c, err)
		return
	}

	if ok {
		fmt.Println("public token: " + publicToken)
		fmt.Println("access token: " + accessToken)
		fmt.Println("item ID: " + itemID)

		c.JSON(http.StatusOK, gin.H{
			"access_token": accessToken,
			"item_id":      itemID,
		})
	}
}

// CreateLinkTokenHandler creates a link token
func (h *PlaidHandlers) CreateLinkTokenHandler(c *gin.Context) {
	products := client.ConvertProducts(strings.Split(h.config.PlaidProducts, ","))
	countryCodes := client.ConvertCountryCodes(strings.Split(h.config.PlaidCountryCodes, ","))

	linkToken, err := h.client.CreateLinkToken(nil, products, countryCodes, h.config.PlaidRedirectURI, userToken)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"link_token": linkToken})
}

// CreateLinkTokenForPaymentHandler creates a link token for payment initiation
func (h *PlaidHandlers) CreateLinkTokenForPaymentHandler(c *gin.Context) {
	// Create payment recipient
	address := *plaid.NewPaymentInitiationAddress(
		[]string{"4 Privet Drive"},
		"Little Whinging",
		"11111",
		"GB",
	)

	recipientResp, err := h.client.CreatePaymentRecipient("Harry Potter", "GB33BUKB20201555555555", address)
	if err != nil {
		renderError(c, err)
		return
	}

	// Create payment
	amount := *plaid.NewPaymentAmount("GBP", 1.34)
	paymentResp, err := h.client.CreatePayment(recipientResp.GetRecipientId(), "paymentRef", amount)
	if err != nil {
		renderError(c, err)
		return
	}

	paymentID = paymentResp.GetPaymentId()
	fmt.Println("payment id: " + paymentID)

	// Create the link_token with payment initiation
	linkTokenCreateReqPaymentInitiation := plaid.NewLinkTokenCreateRequestPaymentInitiation()
	linkTokenCreateReqPaymentInitiation.SetPaymentId(paymentID)

	products := client.ConvertProducts(strings.Split(h.config.PlaidProducts, ","))
	countryCodes := client.ConvertCountryCodes(strings.Split(h.config.PlaidCountryCodes, ","))

	linkToken, err := h.client.CreateLinkToken(linkTokenCreateReqPaymentInitiation, products, countryCodes, h.config.PlaidRedirectURI, userToken)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"link_token": linkToken,
	})
}

// CreateUserTokenHandler creates a user token
func (h *PlaidHandlers) CreateUserTokenHandler(c *gin.Context) {
	products := client.ConvertProducts(strings.Split(h.config.PlaidProducts, ","))

	var consumerReportUserIdentity *plaid.ConsumerReportUserIdentity
	if containsProduct(products, plaid.PRODUCTS_CRA_BASE_REPORT) ||
		containsProduct(products, plaid.PRODUCTS_CRA_INCOME_INSIGHTS) ||
		containsProduct(products, plaid.PRODUCTS_CRA_PARTNER_INSIGHTS) {

		city := "New York"
		region := "NY"
		street := "4 Privet Drive"
		postalCode := "11111"
		country := "US"
		addressData := plaid.AddressData{
			City:       *plaid.NewNullableString(&city),
			Region:     *plaid.NewNullableString(&region),
			Street:     street,
			PostalCode: *plaid.NewNullableString(&postalCode),
			Country:    *plaid.NewNullableString(&country),
		}

		consumerReportUserIdentity = plaid.NewConsumerReportUserIdentity(
			"Harry",
			"Potter",
			[]string{"+16174567890"},
			[]string{"harrypotter@example.com"},
			addressData,
		)
	}

	userResp, err := h.client.CreateUser(time.Now().String(), consumerReportUserIdentity)
	if err != nil {
		renderError(c, err)
		return
	}

	userToken = userResp.GetUserToken()
	c.JSON(http.StatusOK, gin.H{"user_token": userToken})
}

// AuthHandler retrieves auth information
func (h *PlaidHandlers) AuthHandler(c *gin.Context) {
	authResp, err := h.client.GetAuth(accessToken)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": authResp.GetAccounts(),
		"numbers":  authResp.GetNumbers(),
	})
}

// AccountsHandler retrieves account information
func (h *PlaidHandlers) AccountsHandler(c *gin.Context) {
	accountsResp, err := h.client.GetAccounts(accessToken)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accountsResp.GetAccounts(),
	})
}

// BalanceHandler retrieves account balance information
func (h *PlaidHandlers) BalanceHandler(c *gin.Context) {
	balanceResp, err := h.client.GetBalance(accessToken)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": balanceResp,
	})
}

// GetPlaidCategoriesHandler retrieves Plaid categories
func (h *PlaidHandlers) GetPlaidCategoriesHandler(c *gin.Context) {
	categoriesResp, err := h.client.GetCategories()
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categoriesResp.GetCategories(),
	})
}

// ItemHandler retrieves item information
func (h *PlaidHandlers) ItemHandler(c *gin.Context) {
	itemResp, err := h.client.GetItem(accessToken)
	if err != nil {
		renderError(c, err)
		return
	}

	countryCodes := client.ConvertCountryCodes(strings.Split(h.config.PlaidCountryCodes, ","))
	institutionResp, err := h.client.GetInstitutionById(*itemResp.GetItem().InstitutionId.Get(), countryCodes)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"item":        itemResp.GetItem(),
		"institution": institutionResp.GetInstitution(),
	})
}

// IdentityHandler retrieves identity information
func (h *PlaidHandlers) IdentityHandler(c *gin.Context) {
	identityResp, err := h.client.GetIdentity(accessToken)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"identity": identityResp.GetAccounts(),
	})
}

// TransactionsHandler retrieves transactions
func (h *PlaidHandlers) TransactionsHandler(c *gin.Context) {

	accessToken, exists := c.Get("AccessToken")
	if !exists || accessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "AccessToken not provided"})
		return
	}
	transactions, err := h.client.GetTransactions(accessToken.(string))
	if err != nil {
		renderError(c, err)
		return
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].GetDate() < transactions[j].GetDate()
	})

	var latestTransactions []plaid.Transaction
	if len(transactions) >= 9 {
		latestTransactions = transactions[len(transactions)-9:]
	} else {
		latestTransactions = transactions
	}

	c.JSON(http.StatusOK, gin.H{
		"latest_transactions": latestTransactions,
	})
}

// CreatePublicTokenHandler creates a public token
func (h *PlaidHandlers) CreatePublicTokenHandler(c *gin.Context) {
	publicTokenResp, err := h.client.CreatePublicToken(accessToken)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"public_token": publicTokenResp.GetPublicToken(),
	})
}

// PaymentHandler retrieves payment information (UK/EU Payment Initiation)
func (h *PlaidHandlers) PaymentHandler(c *gin.Context) {
	paymentResp, err := h.client.GetPayment(paymentID)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment": paymentResp,
	})
}

// Helper functions

// renderError handles error responses consistently
func renderError(c *gin.Context, originalErr error) {
	if plaidError, err := plaid.ToPlaidError(originalErr); err == nil {
		// Return 200 and allow the front end to render the error.
		c.JSON(http.StatusOK, gin.H{"error": plaidError})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": originalErr.Error()})
}

// containsProduct checks if a product is in the products slice
func containsProduct(products []plaid.Products, product plaid.Products) bool {
	for _, p := range products {
		if p == product {
			return true
		}
	}
	return false
}
