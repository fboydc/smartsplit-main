package main

import (
	"log"

	plaidclient "github.com/fboydc/smartsplit-main/client"
	"github.com/fboydc/smartsplit-main/config"
	"github.com/fboydc/smartsplit-main/handlers"
	"github.com/fboydc/smartsplit-main/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.MustLoadConfig()

	// Initialize database
	db, err := services.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer services.CloseDB(db)

	// Initialize Plaid client
	plaidClient := plaidclient.NewPlaidClient(cfg.PlaidClientID, cfg.PlaidSecret, cfg.PlaidEnv)

	// Initialize handlers
	authHandlers := handlers.NewAuthHandlers(db)
	plaidHandlers := handlers.NewPlaidHandlers(plaidClient, cfg, db, authHandlers)
	budgetHandlers := handlers.NewBudgetHandlers(db)
	dummyHandlers := handlers.NewDummyHandlers()

	// Initialize Gin router
	r := gin.Default()

	// Public routes
	r.POST("/api/auth/login", authHandlers.LoginHandler)

	// Protected routes group
	protected := r.Group("/")
	protected.Use(authHandlers.AuthMiddleware())
	{
		// Info endpoint
		protected.POST("/api/info", plaidHandlers.InfoHandler)

		// Plaid Link Token endpoints
		protected.POST("/api/create_link_token", plaidHandlers.CreateLinkTokenHandler)
		protected.POST("/api/create_link_token_for_payment", plaidHandlers.CreateLinkTokenForPaymentHandler)
		protected.POST("/api/create_user_token", plaidHandlers.CreateUserTokenHandler)
		protected.GET("/api/create_public_token", plaidHandlers.CreatePublicTokenHandler)

		// Plaid Access Token endpoints
		protected.POST("/api/set_access_token", plaidHandlers.GetAccessTokenHandler)

		// Authentication endpoint
		protected.GET("/api/auth", plaidHandlers.AuthHandler)

		// Account endpoints
		protected.GET("/api/accounts", plaidHandlers.AccountsHandler)
		protected.GET("/api/balance", plaidHandlers.BalanceHandler)

		// Category endpoints
		protected.GET("/api/plaid_categories", plaidHandlers.GetPlaidCategoriesHandler)
		protected.GET("/api/categories", budgetHandlers.GetCategoriesHandler)

		// Item endpoints
		protected.GET("/api/item", plaidHandlers.ItemHandler)
		protected.POST("/api/item", plaidHandlers.ItemHandler)

		// Identity endpoint
		protected.GET("/api/identity", plaidHandlers.IdentityHandler)

		// Transaction endpoints
		protected.GET("/api/transactions", plaidHandlers.TransactionsHandler)
		protected.POST("/api/transactions", plaidHandlers.TransactionsHandler)

		// Payment endpoint (UK/EU Payment Initiation)
		protected.GET("/api/payment", plaidHandlers.PaymentHandler)

		// Investment endpoints
		protected.GET("/api/investments_transactions", handleInvestmentTransactions)
		protected.GET("/api/holdings", handleHoldings)

		// Asset endpoints
		protected.GET("/api/assets", handleAssets)

		// Transfer endpoints (ACH)
		protected.GET("/api/transfer_authorize", handleTransferAuthorize)
		protected.GET("/api/transfer_create", handleTransferCreate)

		// Signal endpoint
		protected.GET("/api/signal_evaluate", handleSignalEvaluate)

		// Statements endpoint
		protected.GET("/api/statements", handleStatements)

		// CRA endpoints
		protected.GET("/api/cra/get_base_report", handleCRABaseReport)
		protected.GET("/api/cra/get_income_insights", handleCRAIncomeInsights)
		protected.GET("/api/cra/get_partner_insights", handleCRAPartnerInsights)

		// Budget endpoints
		protected.POST("/api/save_budget", budgetHandlers.SaveBudgetHandler)
		protected.GET("/api/budget", budgetHandlers.GetBudgetHandler)

		// Dummy data endpoints
		protected.GET("/api/dummy/transactions", dummyHandlers.GetDummyTransactionsHandler)
	}

	// Start server
	err = r.Run(":" + cfg.AppPort)
	if err != nil {
		log.Fatal("Unable to start server:", err)
	}
}

// TODO: These handlers need to be implemented or moved to appropriate handler packages
// They are placeholders for now to maintain API compatibility

func handleInvestmentTransactions(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Investment transactions endpoint - implementation needed"})
}

func handleHoldings(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Holdings endpoint - implementation needed"})
}

func handleAssets(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Assets endpoint - implementation needed"})
}

func handleTransferAuthorize(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Transfer authorize endpoint - implementation needed"})
}

func handleTransferCreate(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Transfer create endpoint - implementation needed"})
}

func handleSignalEvaluate(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Signal evaluate endpoint - implementation needed"})
}

func handleStatements(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Statements endpoint - implementation needed"})
}

func handleCRABaseReport(c *gin.Context) {
	c.JSON(200, gin.H{"message": "CRA base report endpoint - implementation needed"})
}

func handleCRAIncomeInsights(c *gin.Context) {
	c.JSON(200, gin.H{"message": "CRA income insights endpoint - implementation needed"})
}

func handleCRAPartnerInsights(c *gin.Context) {
	c.JSON(200, gin.H{"message": "CRA partner insights endpoint - implementation needed"})
}
