package pdf

import (
	"fmt"
	"io"
	"strings"
	"time"

	"4d63.com/clippercardtransactionhistory"

	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdfmodel "github.com/unidoc/unidoc/pdf/model"
)

// Parse converts a reader to a slice of Transactions. The reader should be a
// container the contents of a PDF file.
func Parse(r io.ReadSeeker) (clippercardtransactionhistory.TransactionHistory, error) {
	history := clippercardtransactionhistory.TransactionHistory{
		Transactions: []clippercardtransactionhistory.Transaction{},
	}

	pdfReader, err := pdfmodel.NewPdfReader(r)
	if err != nil {
		return history, fmt.Errorf("error loading pdf: %s", err)
	}

	encrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return history, fmt.Errorf("error loading encryption details from pdf: %s", err)
	}
	if encrypted {
		return history, fmt.Errorf("pdf is encrypted and cannot be read")
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return history, fmt.Errorf("error getting number of pages: %s", err)
	}

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return history, fmt.Errorf("error getting page %d of %d: %s", pageNum, numPages, err)
		}

		transactions, err := parsePage(page)
		if err != nil {
			return history, fmt.Errorf("error converting page %d of %d: %s", pageNum, numPages, err)
		}

		history.Transactions = append(history.Transactions, transactions...)
	}

	return history, nil
}

func parsePage(page *pdfmodel.PdfPage) ([]clippercardtransactionhistory.Transaction, error) {
	contentStreams, err := page.GetContentStreams()
	if err != nil {
		return nil, fmt.Errorf("error getting content streams: %s", err)
	}

	contentStream := ""
	for _, cs := range contentStreams {
		contentStream += cs
	}

	parser := pdfcontent.NewContentStreamParser(contentStream)
	ops, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing content stream: %s", err)
	}

	return parseOperations(ops)
}

func parseOperations(ops *pdfcontent.ContentStreamOperations) ([]clippercardtransactionhistory.Transaction, error) {
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
	storeTransaction := func(columns [8]string) {
		transactions = append(transactions, clippercardtransactionhistory.Transaction{
			Timestamp:       columns[0],
			TransactionType: columns[1],
			Location:        columns[2],
			Route:           columns[3],
			Product:         columns[4],
			Debit:           columns[5],
			Credit:          columns[6],
			Balance:         columns[7],
		})
	}

	for _, op := range *ops {
		if op.Operand == "Tm" && len(op.Params) == 6 {
			switch x := op.Params[4].(type) {
			case *pdfcore.PdfObjectFloat:
				lastX = int(float64(*x))
			case *pdfcore.PdfObjectInteger:
				lastX = int(*x)
			default:
				return nil, fmt.Errorf("invalid Tm parameters: %v", op.Params)
			}
		} else if op.Operand == "Tj" && len(op.Params) == 1 {
			switch param := op.Params[0].(type) {
			case *pdfcore.PdfObjectString:
				text := param.String()
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
						storeTransaction(columns)
						clearStringSlice(&columns)
					}
					columns[columnIndex] = text
				}
			default:
				return nil, fmt.Errorf("invalid Tj parameters: %v", op.Params)
			}
		}
	}

	if !stringSliceBlank(columns) {
		storeTransaction(columns)
	}

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
