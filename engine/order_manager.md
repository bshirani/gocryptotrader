# GoCryptoTrader package Order manager

## Current Features for Order manager
+ The order manager subsystem stores and monitors all orders from enabled exchanges with API keys and `authenticatedSupport` enabled
+ It can be enabled or disabled via runtime command `-ordermanager=false` and defaults to true
+ All orders placed via GoCryptoTrader will be added to the order manager store

