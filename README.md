# clippercard
Tools and packages relating to [ClipperCard](https://www.clippercard.com).

## Apps/Tools

### clippercardcsv (web)

[clippercardcsv.com](https://clippercardcsv.com)

### clippercardcsv (cli)

#### Install

##### Linux, macOS, Windows

Download and install the binary from the [releases](https://github.com/leighmcculloch/clippercard/releases) page.

##### macOS

```
brew install 4d63/clippercard/clippercardcsv
```

##### Source
```
go get 4d63.com/clippercard/transactionhistory/apps/cmd/clippercardcsv
```

#### Usage

```
Usage of clippercardcsv:

Examples:
  clippercardcsv ridehistory.pdf
  clippercardcsv ridehistory.pdf > ridehistory.csv
  cat ridehistory.pdf | clippercardcsv > ridehistory.csv

Flags:
  -filter-weekdays string
        Weekdays to filter by, only transactions occurring on these weekdays will be included in the CSV (default "monday,tuesday,wednesday,thursday,friday,saturday,sunday")
  -headings
        Include headings on columns (default true)
  -help
        Print this help
  -version
        Print version
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
