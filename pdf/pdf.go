package pdf

import (
	"fmt"
	"io"
	"strings"
	"time"

	"4d63.com/clippercardtransactionhistory"
	rscpdf "rsc.io/pdf"
)

// Parse converts a reader to a slice of Transactions. The reader should be a
// container the contents of a PDF file.
func Parse(r io.ReaderAt, size int64) (clippercardtransactionhistory.TransactionHistory, error) {
	history := clippercardtransactionhistory.TransactionHistory{
		Transactions: []clippercardtransactionhistory.Transaction{},
	}

	pdfReader, err := rscpdf.NewReader(r, size)
	if err != nil {
		return history, fmt.Errorf("error loading pdf: %s", err)
	}

	numPages := pdfReader.NumPage()

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page := pdfReader.Page(pageNum)

		transactions, err := parsePage(page)
		if err != nil {
			return history, fmt.Errorf("error converting page %d of %d: %s", pageNum, numPages, err)
		}

		history.Transactions = append(history.Transactions, transactions...)
	}

	return history, nil
}

func parsePage(page rscpdf.Page) ([]clippercardtransactionhistory.Transaction, error) {
	contents := page.V.Key("Contents")

	columHeadingsIndexes := map[string]int{
		"TRANSACTION TYPE": 1,
		"LOCATION":         2,
		"ROUTE":            3,
		"PRODUCT":          4,
		"DEBIT":            5,
		"CREDIT":           6,
		"BALANCE*":         7,
	}
	columnXs := [8]int{}
	columns := [8]string{}
	lastX := 0

	transactions := []clippercardtransactionhistory.Transaction{}

	rscpdf.Interpret(contents, func(stk *rscpdf.Stack, op string) {
		params := make([]rscpdf.Value, stk.Len())
		for i := stk.Len() - 1; i >= 0; i-- {
			params[i] = stk.Pop()
		}
		if op == "Tm" && len(params) == 6 {
			switch params[4].Kind() {
			case rscpdf.Real:
				lastX = int(params[4].Float64())
			case rscpdf.Integer:
				lastX = int(params[4].Int64())
			}
		} else if op == "Tj" && len(params) == 1 {
			switch params[0].Kind() {
			case rscpdf.String:
				text := params[0].Text()
				if ignoreText(text) {
					// ignore
				} else if columnIndex, ok := columHeadingsIndexes[text]; ok {
					columnXs[columnIndex] = lastX
				} else {
					var columnIndex int
					for columnIndex = len(columnXs) - 1; columnIndex >= 0; columnIndex-- {
						if lastX >= columnXs[columnIndex] {
							break
						}
					}
					if columnIndex == 0 && !stringSliceBlank(columns) {
						var err error
						transactions, err = appendTransaction(transactions, columns)
						if err != nil {
							//return nil, err
						}
						clearStringSlice(&columns)
					}
					columns[columnIndex] = text
				}
			}
		}
	})

	if !stringSliceBlank(columns) {
		var err error
		transactions, err = appendTransaction(transactions, columns)
		if err != nil {
			return nil, err
		}
	}

	return transactions, nil
}

func appendTransaction(transactions []clippercardtransactionhistory.Transaction, columns [8]string) ([]clippercardtransactionhistory.Transaction, error) {
	t, err := time.Parse("01/02/2006 15:04 PM", columns[0])
	if err != nil {
		return nil, err
	}
	transactions = append(transactions, clippercardtransactionhistory.Transaction{
		Timestamp:       t,
		TransactionType: columns[1],
		Location:        columns[2],
		Route:           columns[3],
		Product:         columns[4],
		Debit:           columns[5],
		Credit:          columns[6],
		Balance:         columns[7],
	})
	return transactions, nil
}

func timeParseable(layout, value string) bool {
	_, err := time.Parse(layout, value)
	return err == nil
}

func ignoreText(s string) bool {
	return strings.TrimSpace(s) == "" ||
		strings.HasPrefix(s, "Page") ||
		strings.HasPrefix(s, "*") ||
		strings.HasPrefix(s, "CARD ") ||
		strings.HasPrefix(s, "TRANSACTION HISTORY FOR") ||
		timeParseable("01/02/2006", s)
}

func stringSliceBlank(slice [8]string) bool {
	for _, s := range slice {
		if s != "" {
			return false
		}
	}
	return true
}

func clearStringSlice(slice *[8]string) {
	for i := range slice {
		slice[i] = ""
	}
}
