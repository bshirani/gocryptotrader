# GoCryptoTrader package Withdraw manager

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">


## Current Features for Withdraw manager
+ The withdraw manager subsystem is responsible for the processing of withdrawal requests and submitting them to exchanges
+ The withdraw manager can be interacted with via GRPC commands such as `WithdrawFiatRequest` and `WithdrawCryptoRequest`
+ Supports caching of responses to allow for quick viewing of withdrawal events via GRPC
+ If the database is enabled, withdrawal events are stored to the database for later viewing
+ Will not process withdrawal events if `dryrun` is true
+ The withdraw manager subsystem is always enabled


