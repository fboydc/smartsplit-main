package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/fboydc/smartsplit-main/models/api"
	"github.com/gin-gonic/gin"
)

// BudgetHandlers contains all budget-related HTTP handlers
type BudgetHandlers struct {
	db *sql.DB
}

// NewBudgetHandlers creates a new BudgetHandlers instance
func NewBudgetHandlers(db *sql.DB) *BudgetHandlers {
	return &BudgetHandlers{
		db: db,
	}
}

// SaveBudgetHandler saves budget information
func (h *BudgetHandlers) SaveBudgetHandler(c *gin.Context) {
	// TODO: Implement budget saving logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Budget save endpoint - implementation needed",
	})
}

// GetBudgetHandler retrieves budget information for a user
func (h *BudgetHandlers) GetBudgetHandler(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User ID is required",
		})
		return
	}

	// Query income data
	incomeQuery := fmt.Sprintf(`SELECT "income_id", "income_amount", "income_frequency" FROM "Income" WHERE "user_id" = '%s'`, userID)
	incomeRows, err := h.db.Query(incomeQuery)
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
	defer incomeRows.Close()

	// Query expenses data
	expensesQuery := fmt.Sprintf(`SELECT "expense_id", "expense_description", "expense_amount", "expense_category", "allocation_type" FROM "Expenses" WHERE "user_id" = '%s'`, userID)
	expensesRows, err := h.db.Query(expensesQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error: Could not execute expenses query",
		})
		return
	}
	defer expensesRows.Close()

	// Query allocations data
	allocationsQuery := fmt.Sprintf(`SELECT "allocation_type", "allocation_description", "allocation_factor" FROM "Allocations" WHERE "user_id" = '%s'`, userID)
	allocationsRows, err := h.db.Query(allocationsQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error: Could not execute allocations query",
		})
		return
	}
	defer allocationsRows.Close()

	// Process income data
	incomes := []api.Income{}
	for incomeRows.Next() {
		var income api.Income
		err = incomeRows.Scan(&income.Id, &income.Amount, &income.Frequency)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error: Could not scan income row",
			})
			return
		}
		incomes = append(incomes, income)
	}

	// Process expenses data
	expenses := []api.Expense{}
	for expensesRows.Next() {
		var expense api.Expense
		err = expensesRows.Scan(&expense.Id, &expense.Description, &expense.Amount, &expense.Category, &expense.AllocationType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error: Could not scan expenses row",
			})
			return
		}
		expenses = append(expenses, expense)
	}

	// Process allocations data
	allocations := []api.Allocation{}
	for allocationsRows.Next() {
		var allocation api.Allocation
		err = allocationsRows.Scan(&allocation.AllocationType, &allocation.AllocationDescription, &allocation.AllocationFactor)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error: Could not scan allocations row",
			})
			return
		}
		allocations = append(allocations, allocation)
	}

	// Create response
	budgetResponse := api.GetBudgetResponse{
		Incomes:     incomes,
		Allocations: allocations,
		Expenses:    expenses,
	}

	c.JSON(http.StatusOK, gin.H{
		"budget": budgetResponse,
	})
}

// GetCategoriesHandler retrieves all categories from the database
func (h *BudgetHandlers) GetCategoriesHandler(c *gin.Context) {
	rows, err := h.db.Query(`SELECT "category_id", "category_name" FROM "Category"`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	categories := []api.Category{}
	for rows.Next() {
		var category api.Category
		err = rows.Scan(&category.ID, &category.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}
