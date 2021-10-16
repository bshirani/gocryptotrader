# GoCryptoTrader package Currency state manager

## Current Features for Currency state manager
+ The state manager keeps currency states up to date, which include:
* Withdrawal - Determines if the currency is allowed to be withdrawn from the exchange.
* Deposit - Determines if the currency is allowed to be deposited to an exchange.
* Trading - Determines if the currency is allowed to be traded on the exchange.

+ This allows for an internal state check to compliment internal and external
strategies.
