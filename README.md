## Exchange Support Table

| Exchange | REST API | Streaming API | FIX API |
|----------|------|-----------|-----|
| Alphapoint | Yes  | Yes        | NA  |
| Binance| Yes  | Yes        | NA  |
| Bitfinex | Yes  | Yes        | NA  |
| Bitflyer | Yes  | No      | NA  |
| Bithumb | Yes  | Yes       | NA  |
| BitMEX | Yes | Yes | NA |
| Bitstamp | Yes  | Yes       | No  |
| Bittrex | Yes | Yes | NA |
| BTCMarkets | Yes | Yes       | NA  |
| BTSE | Yes | Yes | NA |
| CoinbasePro | Yes | Yes | No|
| Coinbene | Yes | Yes | No |
| COINUT | Yes | Yes | NA |
| Exmo | Yes | NA | NA |
| FTX | Yes | Yes | No |
| GateIO | Yes | Yes | NA |
| Gemini | Yes | Yes | No |
| HitBTC | Yes | Yes | No |
| Huobi.Pro | Yes | Yes | NA |
| ItBit | Yes | NA | No |
| Kraken | Yes | Yes | NA |
| Lbank | Yes | No | NA |
| LocalBitcoins | Yes | NA | NA |
| OKCoin International | Yes | Yes | No |
| OKEX | Yes | Yes | No |
| Poloniex | Yes | Yes | NA |
| Yobit | Yes | NA | NA |
| ZB.COM | Yes | Yes | NA |

We are aiming to support the top 30 exchanges sorted by average liquidity as [ranked by CoinMarketCap](https://coinmarketcap.com/rankings/exchanges/).
However, we welcome pull requests for any exchange which does not match this criterion. If you need help with this, please join us on [Slack](https://join.slack.com/t/gocryptotrader/shared_invite/enQtNTQ5NDAxMjA2Mjc5LTc5ZDE1ZTNiOGM3ZGMyMmY1NTAxYWZhODE0MWM5N2JlZDk1NDU0YTViYzk4NTk3OTRiMDQzNGQ1YTc4YmRlMTk).

** NA means not applicable as the exchange does not support the feature.

## Current Features

+ Support for all exchange fiat and digital currencies, with the ability to individually toggle them on/off.
+ AES256 encrypted config file.
+ REST API support for all exchanges.
+ Websocket support for applicable exchanges.
+ Ability to turn off/on certain exchanges.
+ Communication packages (Slack, SMS via SMSGlobal, Telegram and SMTP).
+ HTTP rate limiter package.
+ Unified API for exchange usage.
+ Customisation of HTTP client features including setting a proxy, user agent and adjusting transport settings.
+ NTP client package.
+ Database support (Postgres and SQLite3). See [database](/database/README.md).
+ OTP generation tool. See [gen otp](/cmd/gen_otp).
+ Connection monitor package.
+ gRPC service and JSON RPC proxy. See [gRPC service](/gctrpc/README.md).
+ gRPC client. See [gctcli](/cmd/gctcli/README.md).
+ Forex currency converter packages (CurrencyConverterAPI, CurrencyLayer, Fixer.io, OpenExchangeRates).
+ Packages for handling currency pairs, tickers and orderbooks.
+ Portfolio management tool; fetches balances from supported exchanges and allows for custom address tracking.
+ Basic event trigger system.
+ OHLCV/Candle retrieval support. See [OHLCV](/docs/OHLCV.md).
+ Scripting support. See [gctscript](/gctscript/README.md).
+ Recent and historic trade processing. See [trades](/exchanges/trade/README.md).
+ Backtesting application. An event-driven backtesting tool to test and iterate trading strategies using historical or custom data. See [backtester](/backtester/README.md).
+ WebGUI (discontinued).

## Compiling instructions

Download and install Go from [Go Downloads](https://golang.org/dl/) for your
platform.

### Linux/OSX

GoCryptoTrader is built using [Go Modules](https://github.com/golang/go/wiki/Modules) and requires Go 1.11 or above
Using Go Modules you now clone this repository **outside** your GOPATH

```bash
git clone https://github.com/thrasher-corp/gocryptotrader.git
cd gocryptotrader
go build
mkdir ~/.gocryptotrader
cp config_example.json ~/.gocryptotrader/config.json
```

### Windows

```bash
git clone https://github.com/thrasher-corp/gocryptotrader.git
cd gocryptotrader
go build
copy config_example.json %APPDATA%\GoCryptoTrader\config.json
```

+ Make any neccessary changes to the `config.json` file.
+ Run the `gocryptotrader` binary file inside your GOPATH bin folder.

## Binaries

Binaries will be published once the codebase reaches a stable condition.
