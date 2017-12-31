# clippercard
Tools and packages relating to [ClipperCard](https://www.clippercard.com).

## Apps/Tools

### clippercardcsv (web)

[clippercardcsv.com](https://clippercardcsv.com)

### clippercardcsv (cli)

```
$ go get 4d63.com/clippercard/transactionhistory/apps/cmd/clippercardcsv
```

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

## Packages

### transactionhistory

```go
import "4d63.com/clippercard/transactionhistory/pdf"
```

```go
transactionHistory, err := pdf.Parse(file)
if err != nil {
	// error parsing
}

for _, t := range transactionHistory.Transactions {
	// do things with transactions
}
```
