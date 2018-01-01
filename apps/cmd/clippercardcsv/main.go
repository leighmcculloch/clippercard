package main // import "4d63.com/clippercard/apps/cmd/clippercardcsv"

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"4d63.com/clippercard/pkgs/transactionhistory/csv"
	"4d63.com/clippercard/pkgs/transactionhistory/filters"
	"4d63.com/clippercard/pkgs/transactionhistory/pdf"
)

var version = "from source"

func main() {
	err := cmd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func cmd() error {
	printHelp := flag.Bool("help", false, "Print this help")
	printVersion := flag.Bool("version", false, "Print version")
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

	if *printHelp {
		flag.Usage()
		return nil
	}

	if *printVersion {
		fmt.Fprintln(os.Stderr, "clippercardcsv", version)
		return nil
	}

	filterWeekdays, err := parseWeekdays(*filterWeekdaysStr)
	if err != nil {
		return err
	}

	args := flag.Args()

	var in io.ReaderAt
	var size int64

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
		size = int64(len(data))
	case 1:
		filename := args[0]
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		in = file

		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}
		size = fileInfo.Size()
	}

	transactionHistory, err := pdf.Parse(in, size)
	if err != nil {
		return fmt.Errorf("error parsing pdf: %s", err)
	}

	transactions := filters.ByWeekday(transactionHistory.Transactions, filterWeekdays)

	err = csv.TransationsToCsv(os.Stdout, transactions, *headings)
	if err != nil {
		return err
	}

	return nil
}
