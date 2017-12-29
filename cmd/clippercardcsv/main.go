package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"4d63.com/clippercardcsv"
)

func main() {
	err := cmd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func cmd() error {
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

	err = clippercardcsv.PdfToCsv(os.Stdout, file)
	if err != nil {
		return err
	}

	return nil
}
