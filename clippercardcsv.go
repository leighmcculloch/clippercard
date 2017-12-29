package clippercardcsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

// PdfToCsv converts a reader to CSV, writing the CSV to the writer. The reader
// should be a data stream containing the contents of a PDF file.
func PdfToCsv(w io.Writer, r io.ReadSeeker) error {
	csvWriter := csv.NewWriter(w)

	pdfReader, err := pdf.NewPdfReader(r)
	if err != nil {
		return fmt.Errorf("error loading pdf: %s", err)
	}

	encrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return fmt.Errorf("error loading encryption details from pdf: %s", err)
	}
	if encrypted {
		return fmt.Errorf("pdf is encrypted and cannot be read")
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return fmt.Errorf("error getting number of pages: %s", err)
	}

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return fmt.Errorf("error getting page %d of %d: %s", i+1, numPages, err)
		}

		err = pdfPageToCsv(csvWriter, page)
		if err != nil {
			return fmt.Errorf("error converting page %d of %d: %s", i+1, numPages, err)
		}
	}

	csvWriter.Flush()

	return nil
}

func pdfPageToCsv(csvWriter *csv.Writer, page *pdf.PdfPage) error {
	contentStreams, err := page.GetContentStreams()
	if err != nil {
		return fmt.Errorf("error getting content streams: %s", err)
	}

	contentStream := ""
	for _, cs := range contentStreams {
		contentStream += cs
	}

	parser := pdfcontent.NewContentStreamParser(contentStream)
	ops, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("error parsing content stream: %s", err)
	}

	return pdfOperationsToCsv(csvWriter, ops)
}

func pdfOperationsToCsv(csvWriter *csv.Writer, ops *pdfcontent.ContentStreamOperations) error {
	columHeadingsIndexes := map[string]int{
		"TRANSACTION TYPE": 1,
		"LOCATION":         2,
		"ROUTE":            3,
		"PRODUCT":          4,
		"DEBIT":            5,
		"CREDIT":           6,
		"BALANCE*":         7,
	}
	columnXs := make([]int, 8)
	columns := make([]string, 8)
	lastX := 0

	for _, op := range *ops {
		if op.Operand == "Tm" && len(op.Params) == 6 {
			switch x := op.Params[4].(type) {
			case *pdfcore.PdfObjectFloat:
				lastX = int(float64(*x))
			case *pdfcore.PdfObjectInteger:
				lastX = int(*x)
			default:
				return fmt.Errorf("invalid Tm parameters: %v", op.Params)
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
					if timeParseable("01/02/2006 15:04 PM", text) && !stringSliceBlank(columns) {
						csvWriter.Write(columns)
						for r := range columns {
							columns[r] = ""
						}
					}
					var i int
					for i = len(columnXs) - 1; i >= 0; i-- {
						if lastX >= columnXs[i] {
							break
						}
					}
					columns[i] = text
				}
			default:
				return fmt.Errorf("invalid Tj parameters: %v", op.Params)
			}
		}
	}

	if !stringSliceBlank(columns) {
		err := csvWriter.Write(columns)
		if err != nil {
			return fmt.Errorf("error writing csv: %s", err)
		}
	}

	return nil
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

func stringSliceBlank(slice []string) bool {
	for _, s := range slice {
		if s != "" {
			return false
		}
	}
	return true
}
