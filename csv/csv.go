package csv

import (
	"encoding/csv"
	"fmt"
	"io"

	"4d63.com/clippercardtransactionhistory"
)

var headings = [8]string{
	"Timestamp",
	"Type",
	"Location",
	"Route",
	"Product",
	"Debit",
	"Credit",
	"Balance",
}

const timestampFormat = "01/02/2006 15:04 PM"

// TransationsToCsv converts slice of Transactions into CSV, writing the CSV to
// the writer.
func TransationsToCsv(w io.Writer, transactions []clippercardtransactionhistory.Transaction, includeHeadings bool) error {
	csvWriter := csv.NewWriter(w)

	if includeHeadings {
		err := csvWriter.Write(headings[:])
		if err != nil {
			return fmt.Errorf("error writing csv: %s", err)
		}
	}

	for _, t := range transactions {
		columns := transactionColumns(t)
		err := csvWriter.Write(columns[:])
		if err != nil {
			return fmt.Errorf("error writing csv: %s", err)
		}
	}

	csvWriter.Flush()

	return nil
}

func transactionColumns(t clippercardtransactionhistory.Transaction) [8]string {
	return [8]string{
		t.Timestamp.Format(timestampFormat),
		t.TransactionType,
		t.Location,
		t.Route,
		t.Product,
		t.Debit,
		t.Credit,
		t.Balance,
	}
}
