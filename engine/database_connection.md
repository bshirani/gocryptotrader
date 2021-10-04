# GoCryptoTrader package Database connection

<img src="/common/gctlogo.png?raw=true" width="350px" height="350px" hspace="70">


## Current Features for Database connection
+ The database connection manager subsystem is used to periodically check whether the application is connected to the database and will provide alerts of any changes
+ In order to modify the behaviour of the database connection manager subsystem, you can edit the following inside your config file under `database`:

### database

| Config | Description | Example |
| ------ | ----------- | ------- |
| enabled | Enabled or disables the database connection subsystem |  `true` |
| verbose | Displays more information to the logger which can be helpful for debugging | `false` |
| driver | The SQL driver to use. Can be `postgres` or `sqlite` | `sqlite` |
| connectionDetails | See below |  |

### connectionDetails

| Config | Description | Example |
| ------ | ----------- | ------- |
| host | The host address of the database |  `localhost` |
| port |  The port used to connect to the database |  `5432` |
| username | An optional username to connect to the database | `username` |
| password | An optional password to connect to the database | `password` |
| database | The name of the database | `database.db` |
| sslmode | The connection type of the database for Postgres databases only | `disable` |

### Please click GoDocs chevron above to view current GoDoc information for this package

## Contribution

Please feel free to submit any pull requests or suggest any desired features to be added.

When submitting a PR, please abide by our coding guidelines:

+ Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting) guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
+ Code must be documented adhering to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary) guidelines.
+ Code must adhere to our [coding style](https://github.com/thrasher-corp/gocryptotrader/blob/master/doc/coding_style.md).
+ Pull requests need to be based on and opened against the `master` branch.

## Donations

<img src="https://github.com/thrasher-corp/gocryptotrader/blob/master/web/src/assets/donate.png?raw=true" hspace="70">

If this framework helped you in any way, or you would like to support the developers working on it, please donate Bitcoin to:

***bc1qk0jareu4jytc0cfrhr5wgshsq8282awpavfahc***
