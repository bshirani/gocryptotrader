# GoCryptoTrader Backtester: Backtest package

This backtest package is part of the GoCryptoTrader codebase.

## Backtest package overview

The backtest package is the most important package of the GoCryptoTrader backtester. It is the engine which combines all elements.
It is responsible for the following functionality
- Loading settings from a provided config file
- Retrieving data
- Loading the data into assessable chunks
- Analysing the data via the `handleEvent` function
- Looping through all data
- Outputting results into a report


A flow of the application is as follows:
![workflow](https://i.imgur.com/Kup6IA9.png)



