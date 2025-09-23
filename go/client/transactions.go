package client

import (
	"context"
	"encoding/json"
	"time"

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
			added = append(added, mapTxn(tx))
		}

		for _, tx := range resp.Modified {
			modified = append(modified, mapTxn(tx))
		}

		for _, r := range resp.Removed {
			removed = append(removed, r.GetTransactionId())
		}

		next = resp.GetNextCursor()
		if !resp.GetHasMore() {
			break
		}

	}

	return SyncResult{
		Added:      added,
		Modified:   modified,
		Removed:    removed,
		NextCursor: next,
	}, nil

	//c.client.PlaidApi.TransactionsSync(ctx).TransactionsSyncRequest(*req).Execute()
}

func mapTxn(p plaid.Transaction) Transaction {

	amt := p.GetAmount()

	raw, _ := json.Marshal(p)

	var cat []string
	if p.Category != nil {
		cat = append(cat, p.Category...)
	}

	var pfcPrimary, pfcDetailed string
	if p.PersonalFinanceCategory.IsSet() {
		pfc := p.GetPersonalFinanceCategory()
		pfcPrimary = pfc.GetPrimary()
		pfcDetailed = pfc.GetDetailed()
	}

	date, _ := time.Parse("2006-01-02", p.GetDate())

	return Transaction{
		TransactionID:   p.GetTransactionId(),
		AccountID:       p.GetAccountId(),
		Name:            p.GetName(),
		Amount:          amt,
		ISOCurrencyCode: p.GetIsoCurrencyCode(),
		Date:            date,
		Pending:         p.GetPending(),
		Category:        cat,
		PFCPrimary:      pfcPrimary,
		PFCDetailed:     pfcDetailed,
		RawJSON:         raw,
	}

}
