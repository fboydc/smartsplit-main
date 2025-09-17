package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	plaid "github.com/plaid/plaid-go/v31/plaid"
)

var (
	PLAID_CLIENT_ID                      = ""
	PLAID_SECRET                         = ""
	PLAID_ENV                            = ""
	PLAID_PRODUCTS                       = ""
	PLAID_COUNTRY_CODES                  = ""
	PLAID_REDIRECT_URI                   = ""
	APP_PORT                             = ""
	client              *plaid.APIClient = nil
	DB                  *sql.DB          = nil
)

var environments = map[string]plaid.Environment{
	"sandbox":    plaid.Sandbox,
	"production": plaid.Production,
}

// Category represents a category in the database

type Allocation struct {
	Id                    uuid.UUID `json:"id"`
	AllocationType        string    `json:"allocation_type"`
	AllocationDescription string    `json:"allocation_description"`
	AllocationFactor      float64   `json:"allocation_factor"`
}

type Category struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Expense struct {
	Id             uuid.UUID `json:"id"`
	Description    string    `json:"description"`
	Amount         float64   `json:"amount"`
	Category       string    `json:"category"`
	AllocationType string    `json:"allocation_type"`
}

type Income struct {
	Id        uuid.UUID `json:"id"`
	Amount    float64   `json:"amount"`
	Frequency string    `json:"frequency"`
}

// LoginRequest represents the login payload
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type getBudgetResponse struct {
	Expenses    []Expense    `json:"expenses"`
	Incomes     []Income     `json:"incomes"`
	Allocations []Allocation `json:"allocations"`
}

// Define the structure of the JSON data
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

type getDummyTransactionsResponse struct {
	LatestTransactions []Transaction `json:"latest_transactions"`
}

func init() {
	// load env vars from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error when loading environment variables from .env file %w", err)
	}

	// set constants from env
	PLAID_CLIENT_ID = os.Getenv("PLAID_CLIENT_ID")
	//PLAID_SECRET = os.Getenv("PLAID_SECRET")
	PLAID_SECRET := os.Getenv("PLAID_SECRET")

	if PLAID_CLIENT_ID == "" || PLAID_SECRET == "" {
		log.Fatal("Error: PLAID_SECRET or PLAID_CLIENT_ID is not set. Did you copy .env.example to .env and fill it out?")
	}

	PLAID_ENV = os.Getenv("PLAID_ENV")
	PLAID_PRODUCTS = os.Getenv("PLAID_PRODUCTS")
	PLAID_COUNTRY_CODES = os.Getenv("PLAID_COUNTRY_CODES")
	PLAID_REDIRECT_URI = os.Getenv("PLAID_REDIRECT_URI")
	APP_PORT = os.Getenv("APP_PORT")

	// set defaults
	if PLAID_PRODUCTS == "" {
		PLAID_PRODUCTS = "transactions"
	}
	if PLAID_COUNTRY_CODES == "" {
		PLAID_COUNTRY_CODES = "US"
	}
	if PLAID_ENV == "" {
		PLAID_ENV = "sandbox"
	}
	if APP_PORT == "" {
		APP_PORT = "8000"
	}
	if PLAID_CLIENT_ID == "" {
		log.Fatal("PLAID_CLIENT_ID is not set. Make sure to fill out the .env file")
	}
	if PLAID_SECRET == "" {
		log.Fatal("PLAID_SECRET is not set. Make sure to fill out the .env file")
	}

	// create Plaid client
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", PLAID_CLIENT_ID)
	configuration.AddDefaultHeader("PLAID-SECRET", PLAID_SECRET)
	configuration.UseEnvironment(environments[PLAID_ENV])
	client = plaid.NewAPIClient(configuration)
}

func main() {
	r := gin.Default()

	DB, err := InitDB()
	if err != nil {
		log.Fatal(err)
	}

	defer CloseDB(DB)

	r.POST("/api/auth/login", loginHandler)

	protected := r.Group("/")
	protected.Use(AuthMiddleware())
	{
		r.POST("/api/info", info)

		// For OAuth flows, the process looks as follows.
		// 1. Create a link token with the redirectURI (as white listed at https://dashboard.plaid.com/team/api).
		// 2. Once the flow succeeds, Plaid Link will redirect to redirectURI with
		// additional parameters (as required by OAuth standards and Plaid).
		// 3. Re-initialize with the link token (from step 1) and the full received redirect URI
		// from step 2.

		protected.POST("/api/set_access_token", getAccessToken)
		protected.POST("/api/create_link_token_for_payment", createLinkTokenForPayment)
		protected.GET("/api/auth", auth)
		protected.GET("/api/accounts", accounts)
		protected.GET("/api/balance", balance)
		protected.GET("/api/plaid_categories", getPlaidCategories)
		protected.GET("/api/categories", getCategories)
		protected.GET("/api/item", item)
		protected.POST("/api/item", item)
		protected.GET("/api/identity", identity)
		protected.GET("/api/transactions", transactions)
		protected.POST("/api/transactions", transactions)
		protected.GET("/api/payment", payment)
		protected.GET("/api/create_public_token", createPublicToken)
		protected.POST("/api/create_link_token", createLinkToken)
		protected.POST("/api/create_user_token", createUserToken)
		protected.GET("/api/investments_transactions", investmentTransactions)
		protected.GET("/api/holdings", holdings)
		protected.GET("/api/assets", assets)
		protected.GET("/api/transfer_authorize", transferAuthorize)
		protected.GET("/api/transfer_create", transferCreate)
		protected.GET("/api/signal_evaluate", signalEvaluate)
		protected.GET("/api/statements", statements)
		protected.GET("/api/cra/get_base_report", getCraBaseReportHandler)
		protected.GET("/api/cra/get_income_insights", getCraIncomeInsightsHandler)
		protected.GET("/api/cra/get_partner_insights", getCraPartnerInsightsHandler)
		protected.POST("/api/save_budget", saveBudgetHandler)
		protected.GET("/api/budget", getBudgetHandler)
		protected.GET("/api/dummy/transactions", getDummyTransactions)
	}

	err = r.Run(":" + APP_PORT)
	if err != nil {
		panic("unable to start server")
	}
}

// We store the access_token and user_token in memory - in production, store it in a secure
// persistent data store.
var accessToken string
var userToken string
var itemID string

var paymentID string

// The authorizationID is only relevant for the Transfer ACH product.
// We store the authorizationID in memory - in production, store it in a secure
// persistent data store
var authorizationID string
var accountID string

func loginHandler(c *gin.Context) {

	//var requestBody loginRequest
	uname := c.PostForm("user")
	passwd := c.PostForm("password")

	auth, userid, plaidToken, err := AuthenthicateUser(uname, passwd, DB)
	if err != nil || !auth {

		if err == sql.ErrNoRows || !auth {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid credentials",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error: Could not authenticate user",
			})
		}

		return
	}

	token, err := GenerateJWT(userid, DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error: Could not generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Login successful",
		"token":      token,
		"plaidToken": plaidToken,
		"user_id":    userid,
		"username":   uname,
	})

}

func getPlaidCategories(c *gin.Context) {
	ctx := context.Background()
	categoriesResp, _, err := client.PlaidApi.CategoriesGet(ctx).Body(nil).Execute()

	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categoriesResp.GetCategories(),
	})
}

func getCategories(c *gin.Context) {
	rows, err := executeQuery(`SELECT "category_id", "category_name" FROM "Category"`, DB)
	if err != nil {
		renderError(c, err)
	}

	returnCategories := []Category{}
	for rows.Next() {
		var category Category
		err = rows.Scan(&category.ID, &category.Name)
		if err != nil {
			renderError(c, err)
		}
		returnCategories = append(returnCategories, category)
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": returnCategories,
	})

}

func renderError(c *gin.Context, originalErr error) {
	if plaidError, err := plaid.ToPlaidError(originalErr); err == nil {
		// Return 200 and allow the front end to render the error.
		c.JSON(http.StatusOK, gin.H{"error": plaidError})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": originalErr.Error()})
}

func getAccessToken(c *gin.Context) {
	publicToken := c.PostForm("public_token")
	user := c.PostForm("user")
	ctx := context.Background()

	// exchange the public_token for an access_token
	exchangePublicTokenResp, _, err := client.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(
		*plaid.NewItemPublicTokenExchangeRequest(publicToken),
	).Execute()
	if err != nil {
		renderError(c, err)
		return
	}

	accessToken = exchangePublicTokenResp.GetAccessToken()
	itemID = exchangePublicTokenResp.GetItemId()

	ok, err := saveAccessToken(accessToken, user, DB)
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

// This functionality is only relevant for the UK/EU Payment Initiation product.
// Creates a link token configured for payment initiation. The payment
// information will be associated with the link token, and will not have to be
// passed in again when we initialize Plaid Link.
// See:
// - https://plaid.com/docs/payment-initiation/
// - https://plaid.com/docs/#payment-initiation-create-link-token-request
func createLinkTokenForPayment(c *gin.Context) {
	ctx := context.Background()

	// Create payment recipient
	paymentRecipientRequest := plaid.NewPaymentInitiationRecipientCreateRequest("Harry Potter")
	paymentRecipientRequest.SetIban("GB33BUKB20201555555555")
	paymentRecipientRequest.SetAddress(*plaid.NewPaymentInitiationAddress(
		[]string{"4 Privet Drive"},
		"Little Whinging",
		"11111",
		"GB",
	))
	paymentRecipientCreateResp, _, err := client.PlaidApi.PaymentInitiationRecipientCreate(ctx).PaymentInitiationRecipientCreateRequest(*paymentRecipientRequest).Execute()
	if err != nil {
		renderError(c, err)
		return
	}

	// Create payment
	paymentCreateRequest := plaid.NewPaymentInitiationPaymentCreateRequest(
		paymentRecipientCreateResp.GetRecipientId(),
		"paymentRef",
		*plaid.NewPaymentAmount("GBP", 1.34),
	)
	paymentCreateResp, _, err := client.PlaidApi.PaymentInitiationPaymentCreate(ctx).PaymentInitiationPaymentCreateRequest(*paymentCreateRequest).Execute()
	if err != nil {
		renderError(c, err)
		return
	}

	// We store the payment_id in memory for demo purposes - in production, store it in a secure
	// persistent data store along with the Payment metadata, such as userId.
	paymentID = paymentCreateResp.GetPaymentId()
	fmt.Println("payment id: " + paymentID)

	// Create the link_token
	linkTokenCreateReqPaymentInitiation := plaid.NewLinkTokenCreateRequestPaymentInitiation()
	linkTokenCreateReqPaymentInitiation.SetPaymentId(paymentID)
	linkToken, err := linkTokenCreate(linkTokenCreateReqPaymentInitiation)
	if err != nil {
		renderError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"link_token": linkToken,
	})
}

func auth(c *gin.Context) {
	ctx := context.Background()

	authGetResp, _, err := client.PlaidApi.AuthGet(ctx).AuthGetRequest(
		*plaid.NewAuthGetRequest(accessToken),
	).Execute()

	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": authGetResp.GetAccounts(),
		"numbers":  authGetResp.GetNumbers(),
	})
}

func accounts(c *gin.Context) {
	ctx := context.Background()

	accountsGetResp, _, err := client.PlaidApi.AccountsGet(ctx).AccountsGetRequest(
		*plaid.NewAccountsGetRequest(accessToken),
	).Execute()

	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accountsGetResp.GetAccounts(),
	})
}

func balance(c *gin.Context) {
	ctx := context.Background()

	balancesGetResp, _, err := client.PlaidApi.AccountsBalanceGet(ctx).AccountsBalanceGetRequest(
		*plaid.NewAccountsBalanceGetRequest(accessToken),
	).Execute()

	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": balancesGetResp.GetAccounts(),
	})
}

func item(c *gin.Context) {
	ctx := context.Background()

	itemGetResp, _, err := client.PlaidApi.ItemGet(ctx).ItemGetRequest(
		*plaid.NewItemGetRequest(accessToken),
	).Execute()

	if err != nil {
		renderError(c, err)
		return
	}

	institutionGetByIdResp, _, err := client.PlaidApi.InstitutionsGetById(ctx).InstitutionsGetByIdRequest(
		*plaid.NewInstitutionsGetByIdRequest(
			*itemGetResp.GetItem().InstitutionId.Get(),
			convertCountryCodes(strings.Split(PLAID_COUNTRY_CODES, ",")),
		),
	).Execute()

	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"item":        itemGetResp.GetItem(),
		"institution": institutionGetByIdResp.GetInstitution(),
	})
}

func identity(c *gin.Context) {
	ctx := context.Background()

	identityGetResp, _, err := client.PlaidApi.IdentityGet(ctx).IdentityGetRequest(
		*plaid.NewIdentityGetRequest(accessToken),
	).Execute()
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"identity": identityGetResp.GetAccounts(),
	})
}

func transactions(c *gin.Context) {
	ctx := context.Background()

	// Set cursor to empty to receive all historical updates
	var cursor *string

	// New transaction updates since "cursor"
	var added []plaid.Transaction
	var modified []plaid.Transaction
	var removed []plaid.RemovedTransaction // Removed transaction ids
	hasMore := true
	// Iterate through each page of new transaction updates for item
	for hasMore {
		request := plaid.NewTransactionsSyncRequest(accessToken)
		if cursor != nil {
			request.SetCursor(*cursor)
		}
		resp, _, err := client.PlaidApi.TransactionsSync(
			ctx,
		).TransactionsSyncRequest(*request).Execute()
		if err != nil {
			renderError(c, err)
			return
		}

		// Update cursor to the next cursor
		nextCursor := resp.GetNextCursor()
		cursor = &nextCursor

		// If no transactions are available yet, wait and poll the endpoint.
		// Normally, we would listen for a webhook, but the Quickstart doesn't
		// support webhooks. For a webhook example, see
		// https://github.com/plaid/tutorial-resources or
		// https://github.com/plaid/pattern

		if *cursor == "" {
			time.Sleep(2 * time.Second)
			continue
		}

		// Add this page of results
		added = append(added, resp.GetAdded()...)
		modified = append(modified, resp.GetModified()...)
		removed = append(removed, resp.GetRemoved()...)
		hasMore = resp.GetHasMore()
	}

	sort.Slice(added, func(i, j int) bool {
		return added[i].GetDate() < added[j].GetDate()
	})
	latestTransactions := added[len(added)-9:]

	c.JSON(http.StatusOK, gin.H{
		"latest_transactions": latestTransactions,
	})
}

// This functionality is only relevant for the UK Payment Initiation product.
// Retrieve Payment for a specified Payment ID
func payment(c *gin.Context) {
	ctx := context.Background()

	paymentGetResp, _, err := client.PlaidApi.PaymentInitiationPaymentGet(ctx).PaymentInitiationPaymentGetRequest(
		*plaid.NewPaymentInitiationPaymentGetRequest(paymentID),
	).Execute()

	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payment": paymentGetResp,
	})
}

// This functionality is only relevant for the ACH Transfer product.
// Create Transfer for a specified Authorization ID

func transferAuthorize(c *gin.Context) {
	ctx := context.Background()
	accountsGetResp, _, err := client.PlaidApi.AccountsGet(ctx).AccountsGetRequest(
		*plaid.NewAccountsGetRequest(accessToken),
	).Execute()

	if err != nil {
		renderError(c, err)
		return
	}

	accountID = accountsGetResp.GetAccounts()[0].AccountId
	transferType, err := plaid.NewTransferTypeFromValue("debit")
	transferNetwork, err := plaid.NewTransferNetworkFromValue("ach")
	ACHClass, err := plaid.NewACHClassFromValue("ppd")

	transferAuthorizationCreateUser := plaid.NewTransferAuthorizationUserInRequest("FirstName LastName")
	transferAuthorizationCreateRequest := plaid.NewTransferAuthorizationCreateRequest(
		accessToken,
		accountID,
		*transferType,
		*transferNetwork,
		"1.00",
		*transferAuthorizationCreateUser)

	transferAuthorizationCreateRequest.SetAchClass(*ACHClass)
	transferAuthorizationCreateResp, _, err := client.PlaidApi.TransferAuthorizationCreate(ctx).TransferAuthorizationCreateRequest(*transferAuthorizationCreateRequest).Execute()

	if err != nil {
		renderError(c, err)
		return
	}

	authorizationID = transferAuthorizationCreateResp.GetAuthorization().Id

	c.JSON(http.StatusOK, transferAuthorizationCreateResp)
}

func transferCreate(c *gin.Context) {
	ctx := context.Background()

	transferCreateRequest := plaid.NewTransferCreateRequest(
		accessToken,
		accountID,
		authorizationID,
		"Debit",
	)

	transferCreateResp, _, err := client.PlaidApi.TransferCreate(ctx).TransferCreateRequest(*transferCreateRequest).Execute()

	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, transferCreateResp)
}

func signalEvaluate(c *gin.Context) {
	ctx := context.Background()
	accountsGetResp, _, err := client.PlaidApi.AccountsGet(ctx).AccountsGetRequest(
		*plaid.NewAccountsGetRequest(accessToken),
	).Execute()

	if err != nil {
		renderError(c, err)
		return
	}

	accountID = accountsGetResp.GetAccounts()[0].AccountId

	signalEvaluateRequest := plaid.NewSignalEvaluateRequest(
		accessToken,
		accountID,
		"txn1234",
		100.00)

	signalEvaluateResp, _, err := client.PlaidApi.SignalEvaluate(ctx).SignalEvaluateRequest(*signalEvaluateRequest).Execute()

	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, signalEvaluateResp)
}

func investmentTransactions(c *gin.Context) {
	ctx := context.Background()

	endDate := time.Now().Local().Format("2006-01-02")
	startDate := time.Now().Local().Add(-30 * 24 * time.Hour).Format("2006-01-02")

	request := plaid.NewInvestmentsTransactionsGetRequest(accessToken, startDate, endDate)
	invTxResp, _, err := client.PlaidApi.InvestmentsTransactionsGet(ctx).InvestmentsTransactionsGetRequest(*request).Execute()

	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"investments_transactions": invTxResp,
	})
}

func holdings(c *gin.Context) {
	ctx := context.Background()

	holdingsGetResp, _, err := client.PlaidApi.InvestmentsHoldingsGet(ctx).InvestmentsHoldingsGetRequest(
		*plaid.NewInvestmentsHoldingsGetRequest(accessToken),
	).Execute()
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"holdings": holdingsGetResp,
	})
}

func info(context *gin.Context) {
	context.JSON(http.StatusOK, map[string]interface{}{
		"item_id":      itemID,
		"access_token": accessToken,
		"products":     strings.Split(PLAID_PRODUCTS, ","),
	})
}

func createPublicToken(c *gin.Context) {
	ctx := context.Background()

	// Create a one-time use public_token for the Item.
	// This public_token can be used to initialize Link in update mode for a user
	publicTokenCreateResp, _, err := client.PlaidApi.ItemCreatePublicToken(ctx).ItemPublicTokenCreateRequest(
		*plaid.NewItemPublicTokenCreateRequest(accessToken),
	).Execute()

	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"public_token": publicTokenCreateResp.GetPublicToken(),
	})
}

func createLinkToken(c *gin.Context) {
	linkToken, err := linkTokenCreate(nil)
	if err != nil {
		renderError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"link_token": linkToken})
}

func createUserToken(c *gin.Context) {
	userToken, err := userTokenCreate()
	if err != nil {
		renderError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_token": userToken})
}

func convertCountryCodes(countryCodeStrs []string) []plaid.CountryCode {
	countryCodes := []plaid.CountryCode{}

	for _, countryCodeStr := range countryCodeStrs {
		countryCodes = append(countryCodes, plaid.CountryCode(countryCodeStr))
	}

	return countryCodes
}

func convertProducts(productStrs []string) []plaid.Products {
	products := []plaid.Products{}

	for _, productStr := range productStrs {
		products = append(products, plaid.Products(productStr))
	}

	return products
}

func containsProduct(products []plaid.Products, product plaid.Products) bool {
	for _, p := range products {
		if p == product {
			return true
		}
	}
	return false
}

// linkTokenCreate creates a link token using the specified parameters
func linkTokenCreate(
	paymentInitiation *plaid.LinkTokenCreateRequestPaymentInitiation,
) (string, error) {
	ctx := context.Background()

	// Institutions from all listed countries will be shown.
	countryCodes := convertCountryCodes(strings.Split(PLAID_COUNTRY_CODES, ","))
	redirectURI := PLAID_REDIRECT_URI

	// This should correspond to a unique id for the current user.
	// Typically, this will be a user ID number from your application.
	// Personally identifiable information, such as an email address or phone number, should not be used here.
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: time.Now().String(),
	}

	request := plaid.NewLinkTokenCreateRequest(
		"Plaid Quickstart",
		"en",
		countryCodes,
		user,
	)

	products := convertProducts(strings.Split(PLAID_PRODUCTS, ","))
	if paymentInitiation != nil {
		request.SetPaymentInitiation(*paymentInitiation)
		// The 'payment_initiation' product has to be the only element in the 'products' list.
		request.SetProducts([]plaid.Products{plaid.PRODUCTS_PAYMENT_INITIATION})
	} else {
		request.SetProducts(products)
	}

	if containsProduct(products, plaid.PRODUCTS_STATEMENTS) {
		statementConfig := plaid.NewLinkTokenCreateRequestStatements(
			time.Now().Local().Add(-30*24*time.Hour).Format("2006-01-02"),
			time.Now().Local().Format("2006-01-02"),
		)
		request.SetStatements(*statementConfig)
	}

	if containsProduct(products, plaid.PRODUCTS_CRA_BASE_REPORT) ||
		containsProduct(products, plaid.PRODUCTS_CRA_INCOME_INSIGHTS) ||
		containsProduct(products, plaid.PRODUCTS_CRA_PARTNER_INSIGHTS) {
		request.SetUserToken(userToken)
		request.SetConsumerReportPermissiblePurpose(plaid.CONSUMERREPORTPERMISSIBLEPURPOSE_ACCOUNT_REVIEW_CREDIT)
		request.SetCraOptions(*plaid.NewLinkTokenCreateRequestCraOptions(60))
	}

	if redirectURI != "" {
		request.SetRedirectUri(redirectURI)
	}

	linkTokenCreateResp, _, err := client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()

	if err != nil {
		return "", err
	}

	return linkTokenCreateResp.GetLinkToken(), nil
}

// Create a user token which can be used for Plaid Check, Income, or Multi-Item link flows
// https://plaid.com/docs/api/users/#usercreate
func userTokenCreate() (string, error) {
	ctx := context.Background()

	request := plaid.NewUserCreateRequest(
		// Typically this will be a user ID number from your application.
		time.Now().String(),
	)

	products := convertProducts(strings.Split(PLAID_PRODUCTS, ","))
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

		request.SetConsumerReportUserIdentity(*plaid.NewConsumerReportUserIdentity(
			"Harry",
			"Potter",
			[]string{"+16174567890"},
			[]string{"harrypotter@example.com"},
			addressData,
		))
	}

	userCreateResp, _, err := client.PlaidApi.UserCreate(ctx).UserCreateRequest(*request).Execute()

	if err != nil {
		return "", err
	}

	userToken = userCreateResp.GetUserToken()

	return userCreateResp.GetUserToken(), nil
}

func statements(c *gin.Context) {
	ctx := context.Background()
	statementsListResp, _, err := client.PlaidApi.StatementsList(ctx).StatementsListRequest(
		*plaid.NewStatementsListRequest(accessToken),
	).Execute()
	statementId := statementsListResp.GetAccounts()[0].GetStatements()[0].StatementId

	statementsDownloadResp, _, err := client.PlaidApi.StatementsDownload(ctx).StatementsDownloadRequest(
		*plaid.NewStatementsDownloadRequest(accessToken, statementId),
	).Execute()
	if err != nil {
		renderError(c, err)
		return
	}

	reader := bufio.NewReader(statementsDownloadResp)
	content, err := io.ReadAll(reader)
	if err != nil {
		renderError(c, err)
		return
	}

	// convert pdf to base64
	encodedPdf := base64.StdEncoding.EncodeToString(content)

	c.JSON(http.StatusOK, gin.H{
		"json": statementsListResp,
		"pdf":  encodedPdf,
	})
}

func assets(c *gin.Context) {
	ctx := context.Background()

	createRequest := plaid.NewAssetReportCreateRequest(10)
	createRequest.SetAccessTokens([]string{accessToken})

	// create the asset report
	assetReportCreateResp, _, err := client.PlaidApi.AssetReportCreate(ctx).AssetReportCreateRequest(
		*createRequest,
	).Execute()
	if err != nil {
		renderError(c, err)
		return
	}

	assetReportToken := assetReportCreateResp.GetAssetReportToken()

	// get the asset report
	assetReportGetResp, err := pollForAssetReport(ctx, client, assetReportToken)
	if err != nil {
		renderError(c, err)
		return
	}

	// get it as a pdf
	pdfRequest := plaid.NewAssetReportPDFGetRequest(assetReportToken)
	pdfFile, _, err := client.PlaidApi.AssetReportPdfGet(ctx).AssetReportPDFGetRequest(*pdfRequest).Execute()
	if err != nil {
		renderError(c, err)
		return
	}

	reader := bufio.NewReader(pdfFile)
	content, err := io.ReadAll(reader)
	if err != nil {
		renderError(c, err)
		return
	}

	// convert pdf to base64
	encodedPdf := base64.StdEncoding.EncodeToString(content)

	c.JSON(http.StatusOK, gin.H{
		"json": assetReportGetResp.GetReport(),
		"pdf":  encodedPdf,
	})
}

func pollForAssetReport(ctx context.Context, client *plaid.APIClient, assetReportToken string) (*plaid.AssetReportGetResponse, error) {
	return pollWithRetries(func() (*plaid.AssetReportGetResponse, error) {
		request := plaid.NewAssetReportGetRequest()
		request.SetAssetReportToken(assetReportToken)
		response, _, err := client.PlaidApi.AssetReportGet(ctx).AssetReportGetRequest(*request).Execute()
		return &response, err
	}, 1000, 20)
}

// Retrieve CRA Base Report and PDF
// Base report: https://plaid.com/docs/check/api/#cracheck_reportbase_reportget
// PDF: https://plaid.com/docs/check/api/#cracheck_reportpdfget
func getCraBaseReportHandler(c *gin.Context) {
	ctx := context.Background()
	getResponse, err := getCraBaseReportWithRetries(ctx, userToken)
	if err != nil {
		renderError(c, err)
		return
	}

	pdfRequest := plaid.NewCraCheckReportPDFGetRequest()
	pdfRequest.SetUserToken(userToken)
	pdfResponse, _, err := client.PlaidApi.CraCheckReportPdfGet(ctx).CraCheckReportPDFGetRequest(*pdfRequest).Execute()
	if err != nil {
		renderError(c, err)
		return
	}

	reader := bufio.NewReader(pdfResponse)
	content, err := io.ReadAll(reader)
	if err != nil {
		renderError(c, err)
		return
	}

	// convert pdf to base64
	encodedPdf := base64.StdEncoding.EncodeToString(content)

	c.JSON(http.StatusOK, gin.H{
		"report": getResponse.Report,
		"pdf":    encodedPdf,
	})
}

func getCraBaseReportWithRetries(ctx context.Context, userToken string) (*plaid.CraCheckReportBaseReportGetResponse, error) {
	return pollWithRetries(func() (*plaid.CraCheckReportBaseReportGetResponse, error) {
		request := plaid.NewCraCheckReportBaseReportGetRequest()
		request.SetUserToken(userToken)
		response, _, err := client.PlaidApi.CraCheckReportBaseReportGet(ctx).CraCheckReportBaseReportGetRequest(*request).Execute()
		return &response, err
	}, 1000, 20)
}

// Retrieve CRA Income Insights and PDF with Insights
// Income insights: https://plaid.com/docs/check/api/#cracheck_reportincome_insightsget
// PDF w/ income insights: https://plaid.com/docs/check/api/#cracheck_reportpdfget
func getCraIncomeInsightsHandler(c *gin.Context) {
	ctx := context.Background()
	getResponse, err := getCraIncomeInsightsWithRetries(ctx, userToken)
	if err != nil {
		renderError(c, err)
		return
	}

	pdfRequest := plaid.NewCraCheckReportPDFGetRequest()
	pdfRequest.SetUserToken(userToken)
	pdfRequest.SetAddOns([]plaid.CraPDFAddOns{plaid.CRAPDFADDONS_CRA_INCOME_INSIGHTS})
	pdfResponse, _, err := client.PlaidApi.CraCheckReportPdfGet(ctx).CraCheckReportPDFGetRequest(*pdfRequest).Execute()
	if err != nil {
		renderError(c, err)
		return
	}

	reader := bufio.NewReader(pdfResponse)
	content, err := io.ReadAll(reader)
	if err != nil {
		renderError(c, err)
		return
	}

	// convert pdf to base64
	encodedPdf := base64.StdEncoding.EncodeToString(content)

	c.JSON(http.StatusOK, gin.H{
		"report": getResponse.Report,
		"pdf":    encodedPdf,
	})
}

func getCraIncomeInsightsWithRetries(ctx context.Context, userToken string) (*plaid.CraCheckReportIncomeInsightsGetResponse, error) {
	return pollWithRetries(func() (*plaid.CraCheckReportIncomeInsightsGetResponse, error) {
		request := plaid.NewCraCheckReportIncomeInsightsGetRequest()
		request.SetUserToken(userToken)
		response, _, err := client.PlaidApi.CraCheckReportIncomeInsightsGet(ctx).CraCheckReportIncomeInsightsGetRequest(*request).Execute()
		return &response, err
	}, 1000, 20)
}

// Retrieve CRA Partner Insights
// https://plaid.com/docs/check/api/#cracheck_reportpartner_insightsget
func getCraPartnerInsightsHandler(c *gin.Context) {
	ctx := context.Background()
	getResponse, err := getCraPartnerInsightsWithRetries(ctx, userToken)
	if err != nil {
		renderError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"report": getResponse.Report,
	})
}

func getCraPartnerInsightsWithRetries(ctx context.Context, userToken string) (*plaid.CraCheckReportPartnerInsightsGetResponse, error) {
	return pollWithRetries(func() (*plaid.CraCheckReportPartnerInsightsGetResponse, error) {
		request := plaid.NewCraCheckReportPartnerInsightsGetRequest()
		request.SetUserToken(userToken)
		response, _, err := client.PlaidApi.CraCheckReportPartnerInsightsGet(ctx).CraCheckReportPartnerInsightsGetRequest(*request).Execute()
		return &response, err
	}, 1000, 20)
}

// Since this quickstart does not support webhooks, this function can be used to poll
// an API that would otherwise be triggered by a webhook.
// For a webhook example, see
// https://github.com/plaid/tutorial-resources or
// https://github.com/plaid/pattern
func pollWithRetries[T any](requestCallback func() (T, error), ms int, retriesLeft int) (T, error) {
	var zero T
	if retriesLeft == 0 {
		return zero, fmt.Errorf("ran out of retries while polling")
	}
	response, err := requestCallback()
	if err != nil {
		plaidErr, err := plaid.ToPlaidError(err)
		if plaidErr.ErrorCode != "PRODUCT_NOT_READY" {
			return zero, err
		}
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return pollWithRetries[T](requestCallback, ms, retriesLeft-1)
	}
	return response, nil
}

func saveBudgetHandler(c *gin.Context) {

}

func getBudgetHandler(c *gin.Context) {

	user_id := c.Query("user_id")
	if user_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID is required",
		})
		return
	}

	income_query := fmt.Sprintf(`SELECT "income_id", "income_amount", "income_frequency" FROM "Income" WHERE "user_id" = '%s'`, user_id)

	income_row, err := DB.Query(income_query)
	if err != nil {

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "No income found for the given user ID",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error: Could not execute income query",
			})

		}

		return
	}

	expenses_query := fmt.Sprintf(`SELECT "expense_id", "expense_description", "expense_amount", "expense_category", "allocation_type" FROM "Expenses" WHERE "user_id" = '%s'`, user_id)

	expenses_row, err := DB.Query(expenses_query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error: Could not execute expenses query",
		})
		return
	}

	allocations_query := fmt.Sprintf(`SELECT "allocation_type", "allocation_description", "allocation_factor" FROM "Allocations" WHERE "user_id" = '%s'`, user_id)

	allocations_row, err := DB.Query(allocations_query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error: Could not execute allocations query",
		})
		return
	}

	defer expenses_row.Close()
	defer income_row.Close()
	defer allocations_row.Close()

	incomes := []Income{}
	expenses := []Expense{}
	allocations := []Allocation{}

	for income_row.Next() {
		var income Income
		err = income_row.Scan(&income.Id, &income.Amount, &income.Frequency)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error: Could not scan income row",
			})
			return
		}
		incomes = append(incomes, income)

	}

	for expenses_row.Next() {
		var expense Expense
		err = expenses_row.Scan(&expense.Id, &expense.Description, &expense.Amount, &expense.Category, &expense.AllocationType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error: Could not scan expenses row",
			})
			return
		}
		expenses = append(expenses, expense)

	}

	for allocations_row.Next() {
		var allocation Allocation
		err = allocations_row.Scan(&allocation.AllocationType, &allocation.AllocationDescription, &allocation.AllocationFactor)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error: Could not scan allocations row",
			})
			return
		}
		allocations = append(allocations, allocation)

	}

	getBudgetResponse := getBudgetResponse{
		Incomes:     incomes,
		Allocations: allocations,
		Expenses:    expenses,
	}

	c.JSON(http.StatusOK, gin.H{
		"budget": getBudgetResponse,
	})

}

func getDummyTransactions(c *gin.Context) {

	// Open the JSON file
	file, err := os.Open("data/test_transactions(first period payment).json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open JSON file"})
		return
	}
	defer file.Close()

	var response getDummyTransactionsResponse
	if err := json.NewDecoder(file).Decode(&response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON file"})
		return
	}

	// Send the parsed data as a JSON respons	e
	c.JSON(http.StatusOK, gin.H{"latest_transactions": response.LatestTransactions})

}
