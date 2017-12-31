# clippercardtransactionhistory
Tool to convert ClipperCard transaction history PDFs to CSV.

## Web

[clippercardcsv.com](https://clippercardcsv.com)

## CLI

### Install

```
$ go get 4d63.com/clippercardtransactionhistory/apps/cmd/clippercardcsv
```

### Usage

```
Usage of clippercardcsv:

Examples:
  clippercardcsv ridehistory.pdf
  clippercardcsv ridehistory.pdf > ridehistory.csv
  cat ridehistory.pdf | clippercardcsv > ridehistory.csv

Flags:
  -help
        Print this help
```

## Package

```go
import "4d63.com/clippercardtransactionhistory"
```

```go
transactionHistory, err := clippercardtransactionhistory.Parse(file)
if err != nil {
	// error parsing
}

for _, t := range transactionHistory.Transactions {
	// do things with transactions
}
```
