package clippercardtransactionhistory

import "time"

// Transaction is a single transaction in a ClipperCard transaction history.
type Transaction struct {
	Timestamp       time.Time
	TransactionType string
	Location        string
	Route           string
	Product         string
	Debit           string
	Credit          string
	Balance         string
}
