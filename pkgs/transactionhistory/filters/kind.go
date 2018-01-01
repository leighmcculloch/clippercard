package filters

import (
	"4d63.com/clippercard/pkgs/transactionhistory"
)

// Kind is a kind of Transaction.
type Kind int

const (
	// Informational transactions do not affect the balance.
	Informational Kind = iota
	// Credit transactions increase the balance.
	Credit
	// Debit transactions reduce the balance.
	Debit
)

// ByKind filters transactions by the given list of kinds, returning only
// transactions that are that kind of transaction.
func ByKind(transactions []transactionhistory.Transaction, kinds []Kind) []transactionhistory.Transaction {
	filtered := []transactionhistory.Transaction{}
	for _, t := range transactions {
		if (kindsContains(kinds, Credit) && t.Credit != "") ||
			(kindsContains(kinds, Debit) && t.Debit != "") ||
			(kindsContains(kinds, Informational) && t.Credit == "" && t.Debit == "") {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

func kindsContains(kinds []Kind, kind Kind) bool {
	for _, k := range kinds {
		if k == kind {
			return true
		}
	}
	return false
}
