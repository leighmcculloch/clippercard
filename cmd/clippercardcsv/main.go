package main // import "4d63.com/clippercardtransactionhistory/cmd/clippercardcsv"

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"4d63.com/clippercardtransactionhistory/csv"
	"4d63.com/clippercardtransactionhistory/filters"
	"4d63.com/clippercardtransactionhistory/pdf"
)

func main() {
	err := cmd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func cmd() error {
	help := flag.Bool("help", false, "Print this help")
	headings := flag.Bool("headings", true, "Include headings on columns")
	filterWeekdaysStr := flag.String("filter-weekdays", allWeekdays, "Weekdays to filter by, only transactions occurring on these weekdays will be included in the CSV")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])

		fmt.Fprintln(os.Stderr, "Examples:")
		fmt.Fprintln(os.Stderr, "  clippercardcsv ridehistory.pdf")
		fmt.Fprintln(os.Stderr, "  clippercardcsv ridehistory.pdf > ridehistory.csv")
		fmt.Fprintln(os.Stderr, "  cat ridehistory.pdf | clippercardcsv > ridehistory.csv")
		fmt.Fprintln(os.Stderr)

		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *help {
		flag.Usage()
		return nil
	}

	filterWeekdays, err := parseWeekdays(*filterWeekdaysStr)
	if err != nil {
		return err
	}

	args := flag.Args()

	var in io.ReadSeeker

	switch len(args) {
	default:
		flag.Usage()
		return nil
	case 0:
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		in = bytes.NewReader(data)
	case 1:
		filename := args[0]
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		in = file
	}

	transactionHistory, err := pdf.Parse(in)
	if err != nil {
		return fmt.Errorf("error parsing pdf: %s", err)
	}

	transactions := filters.Weekday(transactionHistory.Transactions, filterWeekdays)

	err = csv.TransationsToCsv(os.Stdout, transactions, *headings)
	if err != nil {
		return err
	}

	return nil
}
