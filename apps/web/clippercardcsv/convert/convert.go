package app

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"4d63.com/clippercard/pkgs/transactionhistory/csv"
	"4d63.com/clippercard/pkgs/transactionhistory/filters"
	"4d63.com/clippercard/pkgs/transactionhistory/pdf"
)

var weekdays = map[string]time.Weekday{
	"sunday":    time.Sunday,
	"monday":    time.Monday,
	"tuesday":   time.Tuesday,
	"wednesday": time.Wednesday,
	"thursday":  time.Thursday,
	"friday":    time.Friday,
	"saturday":  time.Saturday,
}

var kinds = map[string]filters.Kind{
	"informational": filters.Informational,
	"credit":        filters.Credit,
	"debit":         filters.Debit,
}

func convertHandler(c context.Context, w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	headings, _ := strconv.ParseBool(r.FormValue("headings"))
	filterWeekdays := []time.Weekday{}
	for s, wd := range weekdays {
		value, _ := strconv.ParseBool(r.FormValue(s))
		if value {
			filterWeekdays = append(filterWeekdays, wd)
		}
	}
	filterKinds := []filters.Kind{}
	for s, k := range kinds {
		value, _ := strconv.ParseBool(r.FormValue(s))
		if value {
			filterKinds = append(filterKinds, k)
		}
	}

	uploadedPdf, _, err := r.FormFile("pdf")
	if err != nil {
		return err
	}
	uploadedPdfBytes, err := ioutil.ReadAll(uploadedPdf)
	if err != nil {
		return err
	}
	uploadedPdfReader := bytes.NewReader(uploadedPdfBytes)

	transactionHistory, err := pdf.Parse(uploadedPdfReader, uploadedPdfReader.Size())
	if err != nil {
		return err
	}

	transactions := transactionHistory.Transactions
	transactions = filters.ByKind(transactions, filterKinds)
	transactions = filters.ByWeekday(transactions, filterWeekdays)

	w.Header().Add("Content-Type", "text/csv")

	return csv.TransationsToCsv(w, transactions, headings)
}
