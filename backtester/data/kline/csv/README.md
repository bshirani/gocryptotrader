# GoCryptoTrader Backtester: Csv package

<img src="/backtester/common/backtester.png?raw=true" width="350px" height="350px" hspace="70">



This csv package is part of the GoCryptoTrader codebase.

## Csv package overview

This package is responsible for the loading of kline data via a CSV file. It can retrieve candle data or trade data which is converted into candle data.

### CSV Format
#### Candle based CSV

| Field | Example |
| ----- | -------- |
| Timestamp | 1546300800 |
| Volume | 3 |
| Open | 1335 |
| High | 1338 |
| Low | 1336 |
| Close | 1337 |

Additionally, you can view an example under `./testdata/binance_BTCUSDT_24h_2019_01_01_2020_01_01.csv`

#### Trade based CSV

| Field | Example |
| ----- | -------- |
| Timestamp | 1546300800 |
| Price | 1337 |
| Amount | 420.69 |

Additionally, you can view an example under `./testdata/binance_BTCUSDT_24h-trades_2020_11_16.csv`



