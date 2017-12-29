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

		contentStreams, err := page.GetContentStreams()
		if err != nil {
			return err
		}

		lastX := 0
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
					s := op.Params[0].(*pdfcore.PdfObjectString).String()
					switch {
					case ignore(s):
					case s == "TRANSACTION TYPE":
						columnXs[1] = lastX
					case s == "LOCATION":
						columnXs[2] = lastX
					case s == "ROUTE":
						columnXs[3] = lastX
					case s == "PRODUCT":
						columnXs[4] = lastX
					case s == "DEBIT":
						columnXs[5] = lastX
					case s == "CREDIT":
						columnXs[6] = lastX
					case s == "BALANCE*":
						columnXs[7] = lastX
					default:
						if timeParseable("01/02/2006 15:04 PM", s) && !stringSliceEmpty(columns) {
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
						columns[i] = s
					}
				}
			}
		}

		if len(columns) > 0 {
			csvWriter.Write(columns)
		}
	}

	csvWriter.Flush()

	return nil
}

func timeParseable(layout, value string) bool {
	_, err := time.Parse(layout, value)
	return err == nil
}

func ignore(s string) bool {
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
