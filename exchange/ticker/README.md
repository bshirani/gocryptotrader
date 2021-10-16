## Current Features for ticker

+ ticker generation.
+ Gets a loaded ticker by exchange, asset type and currency pair.
+ Attaches methods to an ticker
	- Returns a string of a value

+ This package is primarily used in conjunction with but not limited to the
exchange interface system set by exchange wrapper orderbook functions in
"exchange"_wrapper.go.

Examples below:

```go
tick, err := yobitExchange.FetchTicker()
if err != nil {
	// Handle error
}

// Converts ticker value to string
tickerValString := tick.PriceToString(...)
```

+ or if you have a routine setting an exchange orderbook you can access it via
the package itself.

```go
tick, err := ticker.GetTicker(...)
if err != nil {
	// Handle error
}
```


