# GoCryptoTrader Unified API

## Unified API

GoCryptoTrader supports a unified API for dealing with exchanges. Each exchange
has its own wrapper file which maps the exchanges own RESTful endpoints into a
standardised way for bot and standalone application usage.

A full breakdown of all the supported wrapper funcs can be found [here.](https://github.com/thrasher-corp/gocryptotrader/blob/master/exchanges/interfaces.go#L21)
Please note that these change on a regular basis as GoCryptoTrader is undergoing
rapid development.

Each exchange supports public API endpoints which don't require any authentication
(fetching ticker, orderbook, trade data) and also private API endpoints (which
require authentication). Some examples include submitting, cancelling and fetching
open orders). To use the authenticated API endpoints, you'll need to set your API
credentials in either the `config.json` file or when you initialise an exchange in
your application, and also have the appropriate key permissions set for the exchange.
Each exchange has a credentials validator which ensures that the API credentials
supplied meet the requirements to make an authenticated request.

## Public API Ticker Example

```go
    var b bitstamp.Bitstamp
    b.SetDefaults()
    ticker, err := b.FetchTicker(currency.NewPair(currency.BTC, currency.USD), asset.Spot)
    if err != nil {
        // Handle error
    }
    fmt.Println(ticker.Last)
```

## Private API Submit Order Example

```go
    var b bitstamp.Bitstamp
    b.SetDefaults()

    b.API.Credentials.Key = "your_key"
    b.API.Credentials.Secret = "your_secret"
    b.API.Credentials.ClientID = "your_clientid"

    o := &order.Submit{
        Pair:      currency.NewPair(currency.BTC, currency.USD),
        OrderSide: order.Sell,
        OrderType: order.Limit,
        Price:     1000000,
        Amount:    0.1,
        AssetType: asset.Spot,
    }
    resp, err := b.SubmitOrder(o)
    if err != nil {
        // Handle error
    }
    fmt.Println(resp.OrderID)
```
