package clippercardtransactionhistory

// Transaction is a single transaction in a ClipperCard transaction history.
type Transaction struct {
	Timestamp       string
	TransactionType string
	Location        string
	Route           string
	Product         string
	Debit           string
	Credit          string
	Balance         string
}
