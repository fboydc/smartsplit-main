package client

import (
	"context"
	"time"

	plaid "github.com/plaid/plaid-go/v31/plaid"
)

// PlaidClient wraps the Plaid API client with custom methods
type PlaidClient struct {
	client      *plaid.APIClient
	environment string
}

// PlaidInterface defines the methods our Plaid client should implement
type PlaidInterface interface {
	// Link Token methods
	CreateLinkToken(paymentInitiation *plaid.LinkTokenCreateRequestPaymentInitiation, products []plaid.Products, countryCodes []plaid.CountryCode, redirectURI string, userToken string) (string, error)

	// Access Token methods
	ExchangePublicToken(publicToken string) (*plaid.ItemPublicTokenExchangeResponse, error)

	// Account methods
	GetAccounts(accessToken string) (*plaid.AccountsGetResponse, error)
	GetBalance(accessToken string) (interface{}, error)
	GetAuth(accessToken string) (*plaid.AuthGetResponse, error)

	// Transaction methods
	GetTransactions(accessToken string) ([]plaid.Transaction, error)

	// Identity methods
	GetIdentity(accessToken string) (*plaid.IdentityGetResponse, error)

	// Item methods
	GetItem(accessToken string) (*plaid.ItemGetResponse, error)
	GetInstitutionById(institutionId string, countryCodes []plaid.CountryCode) (*plaid.InstitutionsGetByIdResponse, error)

	// Categories
	GetCategories() (*plaid.CategoriesGetResponse, error)

	// Payment Initiation (UK/EU)
	CreatePaymentRecipient(name string, iban string, address plaid.PaymentInitiationAddress) (*plaid.PaymentInitiationRecipientCreateResponse, error)
	CreatePayment(recipientId string, reference string, amount plaid.PaymentAmount) (*plaid.PaymentInitiationPaymentCreateResponse, error)
	GetPayment(paymentId string) (*plaid.PaymentInitiationPaymentGetResponse, error)

	// Transfer (ACH)
	AuthorizeTransfer(accessToken string, accountId string, transferType plaid.TransferType, network plaid.TransferNetwork, amount string, user plaid.TransferAuthorizationUserInRequest, achClass plaid.ACHClass) (*plaid.TransferAuthorizationCreateResponse, error)
	CreateTransfer(accessToken string, accountId string, authorizationId string, description string) (*plaid.TransferCreateResponse, error)

	// Signal
	EvaluateSignal(accessToken string, accountId string, clientTransactionId string, amount float64) (*plaid.SignalEvaluateResponse, error)

	// Investments
	GetInvestmentTransactions(accessToken string, startDate string, endDate string) (*plaid.InvestmentsTransactionsGetResponse, error)
	GetHoldings(accessToken string) (*plaid.InvestmentsHoldingsGetResponse, error)

	// Assets
	CreateAssetReport(accessTokens []string, daysRequested int32) (*plaid.AssetReportCreateResponse, error)
	GetAssetReport(assetReportToken string) (*plaid.AssetReportGetResponse, error)
	GetAssetReportPDF(assetReportToken string) (*[]byte, error)

	// Statements
	GetStatements(accessToken string) (*plaid.StatementsListResponse, error)
	DownloadStatement(accessToken string, statementId string) (*[]byte, error)

	// CRA (Credit Reporting Agency)
	GetCRABaseReport(userToken string) (*plaid.CraCheckReportBaseReportGetResponse, error)
	GetCRAIncomeInsights(userToken string) (*plaid.CraCheckReportIncomeInsightsGetResponse, error)
	GetCRAPartnerInsights(userToken string) (*plaid.CraCheckReportPartnerInsightsGetResponse, error)
	GetCRAReportPDF(userToken string, addOns []plaid.CraPDFAddOns) (*[]byte, error)

	// User Token
	CreateUser(userId string, consumerReportUserIdentity *plaid.ConsumerReportUserIdentity) (*plaid.UserCreateResponse, error)

	// Public Token
	CreatePublicToken(accessToken string) (*plaid.ItemPublicTokenCreateResponse, error)
}

// NewPlaidClient creates a new Plaid client instance
func NewPlaidClient(clientID, secret, environment string) *PlaidClient {
	environments := map[string]plaid.Environment{
		"sandbox":    plaid.Sandbox,
		"production": plaid.Production,
	}

	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", clientID)
	configuration.AddDefaultHeader("PLAID-SECRET", secret)
	configuration.UseEnvironment(environments[environment])

	return &PlaidClient{
		client:      plaid.NewAPIClient(configuration),
		environment: environment,
	}
}

// CreateLinkToken creates a link token using the specified parameters
func (pc *PlaidClient) CreateLinkToken(
	paymentInitiation *plaid.LinkTokenCreateRequestPaymentInitiation,
	products []plaid.Products,
	countryCodes []plaid.CountryCode,
	redirectURI string,
	userToken string,
) (string, error) {
	ctx := context.Background()

	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: time.Now().String(),
	}

	request := plaid.NewLinkTokenCreateRequest(
		"Plaid Quickstart",
		"en",
		countryCodes,
		user,
	)

	if paymentInitiation != nil {
		request.SetPaymentInitiation(*paymentInitiation)
		request.SetProducts([]plaid.Products{plaid.PRODUCTS_PAYMENT_INITIATION})
	} else {
		request.SetProducts(products)
	}

	if pc.containsProduct(products, plaid.PRODUCTS_STATEMENTS) {
		statementConfig := plaid.NewLinkTokenCreateRequestStatements(
			time.Now().Local().Add(-30*24*time.Hour).Format("2006-01-02"),
			time.Now().Local().Format("2006-01-02"),
		)
		request.SetStatements(*statementConfig)
	}

	if pc.containsProduct(products, plaid.PRODUCTS_CRA_BASE_REPORT) ||
		pc.containsProduct(products, plaid.PRODUCTS_CRA_INCOME_INSIGHTS) ||
		pc.containsProduct(products, plaid.PRODUCTS_CRA_PARTNER_INSIGHTS) {
		if userToken != "" {
			request.SetUserToken(userToken)
		}
		request.SetConsumerReportPermissiblePurpose(plaid.CONSUMERREPORTPERMISSIBLEPURPOSE_ACCOUNT_REVIEW_CREDIT)
		request.SetCraOptions(*plaid.NewLinkTokenCreateRequestCraOptions(60))
	}

	if redirectURI != "" {
		request.SetRedirectUri(redirectURI)
	}

	linkTokenCreateResp, _, err := pc.client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", err
	}

	return linkTokenCreateResp.GetLinkToken(), nil
}

// ExchangePublicToken exchanges a public token for an access token
func (pc *PlaidClient) ExchangePublicToken(publicToken string) (*plaid.ItemPublicTokenExchangeResponse, error) {
	ctx := context.Background()

	response, _, err := pc.client.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(
		*plaid.NewItemPublicTokenExchangeRequest(publicToken),
	).Execute()

	return &response, err
}

// GetAccounts retrieves account information
func (pc *PlaidClient) GetAccounts(accessToken string) (*plaid.AccountsGetResponse, error) {
	ctx := context.Background()

	response, _, err := pc.client.PlaidApi.AccountsGet(ctx).AccountsGetRequest(
		*plaid.NewAccountsGetRequest(accessToken),
	).Execute()

	return &response, err
}

// GetBalance retrieves account balance information
func (pc *PlaidClient) GetBalance(accessToken string) (interface{}, error) {
	ctx := context.Background()

	response, _, err := pc.client.PlaidApi.AccountsBalanceGet(ctx).AccountsBalanceGetRequest(
		*plaid.NewAccountsBalanceGetRequest(accessToken),
	).Execute()

	return response, err
}

// GetAuth retrieves auth information
func (pc *PlaidClient) GetAuth(accessToken string) (*plaid.AuthGetResponse, error) {
	ctx := context.Background()

	response, _, err := pc.client.PlaidApi.AuthGet(ctx).AuthGetRequest(
		*plaid.NewAuthGetRequest(accessToken),
	).Execute()

	return &response, err
}

// GetTransactions retrieves transactions using sync endpoint
func (pc *PlaidClient) GetTransactions(accessToken string) ([]plaid.Transaction, error) {
	ctx := context.Background()

	var cursor *string
	var added []plaid.Transaction
	var modified []plaid.Transaction
	var removed []plaid.RemovedTransaction
	hasMore := true

	for hasMore {
		request := plaid.NewTransactionsSyncRequest(accessToken)
		if cursor != nil {
			request.SetCursor(*cursor)
		}

		resp, _, err := pc.client.PlaidApi.TransactionsSync(ctx).TransactionsSyncRequest(*request).Execute()
		if err != nil {
			return nil, err
		}

		nextCursor := resp.GetNextCursor()
		cursor = &nextCursor

		if *cursor == "" {
			time.Sleep(2 * time.Second)
			continue
		}

		added = append(added, resp.GetAdded()...)
		modified = append(modified, resp.GetModified()...)
		removed = append(removed, resp.GetRemoved()...)
		hasMore = resp.GetHasMore()
	}

	return added, nil
}

// GetIdentity retrieves identity information
func (pc *PlaidClient) GetIdentity(accessToken string) (*plaid.IdentityGetResponse, error) {
	ctx := context.Background()

	response, _, err := pc.client.PlaidApi.IdentityGet(ctx).IdentityGetRequest(
		*plaid.NewIdentityGetRequest(accessToken),
	).Execute()

	return &response, err
}

// GetItem retrieves item information
func (pc *PlaidClient) GetItem(accessToken string) (*plaid.ItemGetResponse, error) {
	ctx := context.Background()

	response, _, err := pc.client.PlaidApi.ItemGet(ctx).ItemGetRequest(
		*plaid.NewItemGetRequest(accessToken),
	).Execute()

	return &response, err
}

// GetInstitutionById retrieves institution information by ID
func (pc *PlaidClient) GetInstitutionById(institutionId string, countryCodes []plaid.CountryCode) (*plaid.InstitutionsGetByIdResponse, error) {
	ctx := context.Background()

	response, _, err := pc.client.PlaidApi.InstitutionsGetById(ctx).InstitutionsGetByIdRequest(
		*plaid.NewInstitutionsGetByIdRequest(institutionId, countryCodes),
	).Execute()

	return &response, err
}

// GetCategories retrieves Plaid categories
func (pc *PlaidClient) GetCategories() (*plaid.CategoriesGetResponse, error) {
	ctx := context.Background()

	response, _, err := pc.client.PlaidApi.CategoriesGet(ctx).Body(nil).Execute()

	return &response, err
}

// CreatePaymentRecipient creates a payment recipient (UK/EU Payment Initiation)
func (pc *PlaidClient) CreatePaymentRecipient(name string, iban string, address plaid.PaymentInitiationAddress) (*plaid.PaymentInitiationRecipientCreateResponse, error) {
	ctx := context.Background()

	request := plaid.NewPaymentInitiationRecipientCreateRequest(name)
	request.SetIban(iban)
	request.SetAddress(address)

	response, _, err := pc.client.PlaidApi.PaymentInitiationRecipientCreate(ctx).PaymentInitiationRecipientCreateRequest(*request).Execute()

	return &response, err
}

// CreatePayment creates a payment (UK/EU Payment Initiation)
func (pc *PlaidClient) CreatePayment(recipientId string, reference string, amount plaid.PaymentAmount) (*plaid.PaymentInitiationPaymentCreateResponse, error) {
	ctx := context.Background()

	request := plaid.NewPaymentInitiationPaymentCreateRequest(recipientId, reference, amount)

	response, _, err := pc.client.PlaidApi.PaymentInitiationPaymentCreate(ctx).PaymentInitiationPaymentCreateRequest(*request).Execute()

	return &response, err
}

// GetPayment retrieves payment information (UK/EU Payment Initiation)
func (pc *PlaidClient) GetPayment(paymentId string) (*plaid.PaymentInitiationPaymentGetResponse, error) {
	ctx := context.Background()

	response, _, err := pc.client.PlaidApi.PaymentInitiationPaymentGet(ctx).PaymentInitiationPaymentGetRequest(
		*plaid.NewPaymentInitiationPaymentGetRequest(paymentId),
	).Execute()

	return &response, err
}

// AuthorizeTransfer authorizes a transfer (ACH)
func (pc *PlaidClient) AuthorizeTransfer(accessToken string, accountId string, transferType plaid.TransferType, network plaid.TransferNetwork, amount string, user plaid.TransferAuthorizationUserInRequest, achClass plaid.ACHClass) (*plaid.TransferAuthorizationCreateResponse, error) {
	ctx := context.Background()

	request := plaid.NewTransferAuthorizationCreateRequest(accessToken, accountId, transferType, network, amount, user)
	request.SetAchClass(achClass)

	response, _, err := pc.client.PlaidApi.TransferAuthorizationCreate(ctx).TransferAuthorizationCreateRequest(*request).Execute()

	return &response, err
}

// CreateTransfer creates a transfer (ACH)
func (pc *PlaidClient) CreateTransfer(accessToken string, accountId string, authorizationId string, description string) (*plaid.TransferCreateResponse, error) {
	ctx := context.Background()

	request := plaid.NewTransferCreateRequest(accessToken, accountId, authorizationId, description)

	response, _, err := pc.client.PlaidApi.TransferCreate(ctx).TransferCreateRequest(*request).Execute()

	return &response, err
}

// EvaluateSignal evaluates a transaction signal
func (pc *PlaidClient) EvaluateSignal(accessToken string, accountId string, clientTransactionId string, amount float64) (*plaid.SignalEvaluateResponse, error) {
	ctx := context.Background()

	request := plaid.NewSignalEvaluateRequest(accessToken, accountId, clientTransactionId, amount)

	response, _, err := pc.client.PlaidApi.SignalEvaluate(ctx).SignalEvaluateRequest(*request).Execute()

	return &response, err
}

// GetInvestmentTransactions retrieves investment transactions
func (pc *PlaidClient) GetInvestmentTransactions(accessToken string, startDate string, endDate string) (*plaid.InvestmentsTransactionsGetResponse, error) {
	ctx := context.Background()

	request := plaid.NewInvestmentsTransactionsGetRequest(accessToken, startDate, endDate)

	response, _, err := pc.client.PlaidApi.InvestmentsTransactionsGet(ctx).InvestmentsTransactionsGetRequest(*request).Execute()

	return &response, err
}

// GetHoldings retrieves investment holdings
func (pc *PlaidClient) GetHoldings(accessToken string) (*plaid.InvestmentsHoldingsGetResponse, error) {
	ctx := context.Background()

	response, _, err := pc.client.PlaidApi.InvestmentsHoldingsGet(ctx).InvestmentsHoldingsGetRequest(
		*plaid.NewInvestmentsHoldingsGetRequest(accessToken),
	).Execute()

	return &response, err
}

// CreateAssetReport creates an asset report
func (pc *PlaidClient) CreateAssetReport(accessTokens []string, daysRequested int32) (*plaid.AssetReportCreateResponse, error) {
	ctx := context.Background()

	request := plaid.NewAssetReportCreateRequest(daysRequested)
	request.SetAccessTokens(accessTokens)

	response, _, err := pc.client.PlaidApi.AssetReportCreate(ctx).AssetReportCreateRequest(*request).Execute()

	return &response, err
}

// GetAssetReport retrieves an asset report
func (pc *PlaidClient) GetAssetReport(assetReportToken string) (*plaid.AssetReportGetResponse, error) {
	ctx := context.Background()

	request := plaid.NewAssetReportGetRequest()
	request.SetAssetReportToken(assetReportToken)

	response, _, err := pc.client.PlaidApi.AssetReportGet(ctx).AssetReportGetRequest(*request).Execute()

	return &response, err
}

// GetAssetReportPDF retrieves an asset report as PDF
func (pc *PlaidClient) GetAssetReportPDF(assetReportToken string) (*[]byte, error) {
	ctx := context.Background()

	request := plaid.NewAssetReportPDFGetRequest(assetReportToken)

	_, _, err := pc.client.PlaidApi.AssetReportPdfGet(ctx).AssetReportPDFGetRequest(*request).Execute()
	if err != nil {
		return nil, err
	}

	// Note: You'll need to handle the PDF reading in the handler
	return nil, nil
}

// GetStatements retrieves statements
func (pc *PlaidClient) GetStatements(accessToken string) (*plaid.StatementsListResponse, error) {
	ctx := context.Background()

	response, _, err := pc.client.PlaidApi.StatementsList(ctx).StatementsListRequest(
		*plaid.NewStatementsListRequest(accessToken),
	).Execute()

	return &response, err
}

// DownloadStatement downloads a statement
func (pc *PlaidClient) DownloadStatement(accessToken string, statementId string) (*[]byte, error) {
	ctx := context.Background()

	_, _, err := pc.client.PlaidApi.StatementsDownload(ctx).StatementsDownloadRequest(
		*plaid.NewStatementsDownloadRequest(accessToken, statementId),
	).Execute()

	if err != nil {
		return nil, err
	}

	// Note: You'll need to handle the PDF reading in the handler
	return nil, nil
}

// GetCRABaseReport retrieves CRA base report
func (pc *PlaidClient) GetCRABaseReport(userToken string) (*plaid.CraCheckReportBaseReportGetResponse, error) {
	ctx := context.Background()

	request := plaid.NewCraCheckReportBaseReportGetRequest()
	request.SetUserToken(userToken)

	response, _, err := pc.client.PlaidApi.CraCheckReportBaseReportGet(ctx).CraCheckReportBaseReportGetRequest(*request).Execute()

	return &response, err
}

// GetCRAIncomeInsights retrieves CRA income insights
func (pc *PlaidClient) GetCRAIncomeInsights(userToken string) (*plaid.CraCheckReportIncomeInsightsGetResponse, error) {
	ctx := context.Background()

	request := plaid.NewCraCheckReportIncomeInsightsGetRequest()
	request.SetUserToken(userToken)

	response, _, err := pc.client.PlaidApi.CraCheckReportIncomeInsightsGet(ctx).CraCheckReportIncomeInsightsGetRequest(*request).Execute()

	return &response, err
}

// GetCRAPartnerInsights retrieves CRA partner insights
func (pc *PlaidClient) GetCRAPartnerInsights(userToken string) (*plaid.CraCheckReportPartnerInsightsGetResponse, error) {
	ctx := context.Background()

	request := plaid.NewCraCheckReportPartnerInsightsGetRequest()
	request.SetUserToken(userToken)

	response, _, err := pc.client.PlaidApi.CraCheckReportPartnerInsightsGet(ctx).CraCheckReportPartnerInsightsGetRequest(*request).Execute()

	return &response, err
}

// GetCRAReportPDF retrieves CRA report as PDF
func (pc *PlaidClient) GetCRAReportPDF(userToken string, addOns []plaid.CraPDFAddOns) (*[]byte, error) {
	ctx := context.Background()

	request := plaid.NewCraCheckReportPDFGetRequest()
	request.SetUserToken(userToken)
	if len(addOns) > 0 {
		request.SetAddOns(addOns)
	}

	_, _, err := pc.client.PlaidApi.CraCheckReportPdfGet(ctx).CraCheckReportPDFGetRequest(*request).Execute()

	if err != nil {
		return nil, err
	}

	// Note: You'll need to handle the PDF reading in the handler
	return nil, nil
}

// CreateUser creates a user for CRA products
func (pc *PlaidClient) CreateUser(userId string, consumerReportUserIdentity *plaid.ConsumerReportUserIdentity) (*plaid.UserCreateResponse, error) {
	ctx := context.Background()

	request := plaid.NewUserCreateRequest(userId)
	if consumerReportUserIdentity != nil {
		request.SetConsumerReportUserIdentity(*consumerReportUserIdentity)
	}

	response, _, err := pc.client.PlaidApi.UserCreate(ctx).UserCreateRequest(*request).Execute()

	return &response, err
}

// CreatePublicToken creates a public token
func (pc *PlaidClient) CreatePublicToken(accessToken string) (*plaid.ItemPublicTokenCreateResponse, error) {
	ctx := context.Background()

	response, _, err := pc.client.PlaidApi.ItemCreatePublicToken(ctx).ItemPublicTokenCreateRequest(
		*plaid.NewItemPublicTokenCreateRequest(accessToken),
	).Execute()

	return &response, err
}

// Helper method to check if a product is contained in the products slice
func (pc *PlaidClient) containsProduct(products []plaid.Products, product plaid.Products) bool {
	for _, p := range products {
		if p == product {
			return true
		}
	}
	return false
}

// ConvertCountryCodes converts string array to CountryCode array
func ConvertCountryCodes(countryCodeStrs []string) []plaid.CountryCode {
	countryCodes := []plaid.CountryCode{}
	for _, countryCodeStr := range countryCodeStrs {
		countryCodes = append(countryCodes, plaid.CountryCode(countryCodeStr))
	}
	return countryCodes
}

// ConvertProducts converts string array to Products array
func ConvertProducts(productStrs []string) []plaid.Products {
	products := []plaid.Products{}
	for _, productStr := range productStrs {
		products = append(products, plaid.Products(productStr))
	}
	return products
}
