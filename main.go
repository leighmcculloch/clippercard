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

		record := []string{}

		for _, contentStream := range contentStreams {
			parser := pdfcontent.NewContentStreamParser(contentStream)
			ops, err := parser.Parse()
			if err != nil {
				return err
			}
			for _, op := range *ops {
				for _, pdfObject := range op.Params {
					switch o := pdfObject.(type) {
					case *pdfcore.PdfObjectString:
						s := o.String()
						switch {
						case strings.TrimSpace(s) == "":
						case strings.HasPrefix(s, "Page"):
						case strings.HasPrefix(s, "*"):
						case strings.HasPrefix(s, "CARD "):
						case strings.HasPrefix(s, "TRANSACTION HISTORY FOR"):
						case timeParseErrNil(time.Parse("01/02/2006", s)):
						case s == "TRANSACTION TYPE":
						case s == "LOCATION":
						case s == "ROUTE":
						case s == "PRODUCT":
						case s == "DEBIT":
						case s == "CREDIT":
						case s == "BALANCE*":
						default:
							_, err := time.Parse("01/02/2006 15:04 PM", s)
							if err == nil && len(record) > 0 {
								csvWriter.Write(record)
								record = record[:0]
							}
							record = append(record, s)
						}
					}
				}
			}
		}

		if len(record) > 0 {
			csvWriter.Write(record)
		}
	}

	csvWriter.Flush()

	return nil
}

func timeParseErrNil(t time.Time, err error) bool {
	return err == nil
}
