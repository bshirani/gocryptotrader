# GoCryptoTrader File Hierarchy

## Default data directory

By default, GoCryptoTrader uses the following data directores:

Operating System | Path | Translated
--- | --- | ----
| Linux | ~/.gocryptotrader | /home/user/.gocryptotrader
| macOS | ~/.gocryptotrader | /Users/User/.gocryptotrader

This can be overridden by running GoCryptoTrader with the `-datadir` command line
parameter.

## Subdirectories

Depending on the features enabled, you'll see the following directories created
inside the data directory:

Directory | Reason
--- | ---
| database | Used to store the database file (if using SQLite3) and sqlboiler config files
| logs | Used to store the debug log file (`log.txt` by default), if file output and logging is enabled
| tls | Used to store the generated self-signed certificate and key for gRPC authentication

## Files

File | Reason
--- | ---
config.json or config.dat (encrypted config) | Config file which GoCryptoTrader loads from (can be overridden by the `-config` command line parameter).
currency.json | Cached list of fiat and digital currencies
