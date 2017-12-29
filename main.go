package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func main() {
	err := app()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func app() error {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		return errors.New("Usage: clippercsv ridehistory.pdf")
	}

	filename := args[0]

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = pdfToCsv(os.Stdout, file)
	if err != nil {
		return err
	}

	return nil
}

func pdfToCsv(w io.Writer, r io.ReadSeeker) error {
	csvWriter := csv.NewWriter(w)

	pdfReader, err := pdf.NewPdfReader(r)
	if err != nil {
		return err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		err = pdfPageToCsv(csvWriter, page)
		if err != nil {
			return err
		}
	}

	csvWriter.Flush()

	return nil
}

func pdfPageToCsv(csvWriter *csv.Writer, page *pdf.PdfPage) error {
	contentStreams, err := page.GetContentStreams()
	if err != nil {
		return err
	}

	lastX := 0
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

	for _, contentStream := range contentStreams {
		parser := pdfcontent.NewContentStreamParser(contentStream)
		ops, err := parser.Parse()
		if err != nil {
			return err
		}
		for _, op := range *ops {
			if op.Operand == "Tm" && len(op.Params) == 6 {
				switch x := op.Params[4].(type) {
				case *pdfcore.PdfObjectFloat:
					lastX = int(float64(*x))
				case *pdfcore.PdfObjectInteger:
					lastX = int(*x)
				}
			} else if op.Operand == "Tj" && len(op.Params) == 1 {
				text := op.Params[0].(*pdfcore.PdfObjectString).String()
				if ignoreContent(text) {
					// ignore
				} else if columnIndex, ok := columHeadingsIndexes[text]; ok {
					columnXs[columnIndex] = lastX
				} else {
					if timeParseable("01/02/2006 15:04 PM", text) && !stringSliceEmpty(columns) {
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
			}
		}
	}

	if len(columns) > 0 {
		csvWriter.Write(columns)
	}

	return nil
}

func timeParseable(layout, value string) bool {
	_, err := time.Parse(layout, value)
	return err == nil
}

func ignoreContent(s string) bool {
	return strings.TrimSpace(s) == "" ||
		strings.HasPrefix(s, "Page") ||
		strings.HasPrefix(s, "*") ||
		strings.HasPrefix(s, "CARD ") ||
		strings.HasPrefix(s, "TRANSACTION HISTORY FOR") ||
		timeParseable("01/02/2006", s)
}

func stringSliceEmpty(slice []string) bool {
	for _, s := range slice {
		if s != "" {
			return false
		}
	}
	return true
}
