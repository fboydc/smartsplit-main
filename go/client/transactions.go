package client

import (
	"context"

	"github.com/fboydc/smartsplit-main/models/domain"
	plaid "github.com/plaid/plaid-go/v31/plaid"
)

// Type aliases for commonly used models in this package
type Transaction = domain.TransactionRaw

type SyncResult struct {
	Added      []Transaction
	Modified   []Transaction
	Removed    []string
	NextCursor string
}

func (c *PlaidClient) TransactionsSync(ctx context.Context, accessToken string, cursor string) (SyncResult, error) {

	client := c.client

	var (
		added    []Transaction
		modified []Transaction
		removed  []string
		next     = cursor
	)

	for {
		req := plaid.NewTransactionsSyncRequest(accessToken)

		if next != "" {
			req.SetCursor(next)
		}

		resp, _, err := client.PlaidApi.TransactionsSync(ctx).TransactionsSyncRequest(*req).Execute()
		if err != nil {
			return SyncResult{}, err
		}

		for _, tx := range resp.Added {
			//added = append(added
		}

	}

	//c.client.PlaidApi.TransactionsSync(ctx).TransactionsSyncRequest(*req).Execute()
}
