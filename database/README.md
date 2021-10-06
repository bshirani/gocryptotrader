# GoCryptoTrader package Database

## Current Features for database package

+ Establishes & Maintains database connection across program life cycle
+ Migration handed by [Goose](https://github.com/thrasher-corp/goose)
+ Model generation handled by [SQLBoiler](https://github.com/thrasher-corp/sqlboiler)

## How to use

##### Prerequisites

[SQLBoiler](https://github.com/thrasher-corp/sqlboiler)
```shell script
go get -u github.com/thrasher-corp/sqlboiler
```

[Postgres Driver](https://github.com/thrasher-corp/sqlboiler/drivers/sqlboiler-psql)
```shell script
go get -u github.com/thrasher-corp/sqlboiler/drivers/sqlboiler-psql
```

[SQLite Driver](https://github.com/thrasher-corp/sqlboiler-sqlite3)
```shell script
go get -u github.com/thrasher-corp/sqlboiler-sqlite3
```

##### Configuration

The database configuration struct is currently:
```shell script
type Config struct {
	Enabled                   bool   `json:"enabled"`
	Verbose                   bool   `json:"verbose"`
	Driver                    string `json:"driver"`
	drivers.ConnectionDetails `json:"connectionDetails"`
}
```
And Connection Details:
```sh
type ConnectionDetails struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSLMode  string `json:"sslmode"`
}
```

With an example configuration being:

```sh
 "database": {
  "enabled": true,
  "verbose": true,
  "driver": "postgres",
  "connectionDetails": {
   "host": "localhost",
   "port": 5432,
   "username": "gct-dev",
   "password": "gct-dev",
   "database": "gct-dev",
   "sslmode": "disable"
  }
 },
```

##### Create and Run migrations
 Migrations are created using a modified version of [Goose](https://github.com/thrasher-corp/goose)

 A helper tool sits in the ./cmd/dbmigrate folder that includes the following features:

+ Check current database version with the "status" command
```shell script
dbmigrate -command status
```

+ Create a new migration
```sh
dbmigrate -command "create" -args "model"
```
_This will create a folder in the ./database/migration folder that contains postgres.sql and sqlite.sql files_
 + Run dbmigrate command with -command up
```shell script
dbmigrate -command "up"
```

dbmigrate provides a -migrationdir flag override to tell it what path to look in for migrations

###### Note: its highly recommended to backup any data before running migrations against a production database especially if you are running SQLite due to alter table limitations


##### Adding a new model
Model's are generated using [SQLBoiler](https://github.com/thrasher-corp/sqlboiler)
A helper tool has been made located in gen_sqlboiler_config that will parse your GoCryptoTrader config and output a SQLBoiler config

```sh
gen_sqlboiler_config
```

By default this will look in your gocryptotrader data folder and default config, these can be overwritten
along with the location of the sqlboiler generated config

```shell script
-config "configname.json"
-datadir "~/.gocryptotrader/"
-outdir "~/.gocryptotrader/"
```

Generate a new model that gets placed in ./database/models/<databasetype> folder

Linux:
```shell script
sqlboiler -o database/models/postgres -p postgres --no-auto-timestamps --wipe psql
```

Helpers have been provided in the Makefile for linux users
```
make gen_db_models
```

##### Adding a Repository
+ Create Repository directory in github.com/thrasher-corp/gocryptotrader/database/repository/


##### DBSeed helper
A helper tool [cmd/dbseed](../cmd/dbseed/README.md) has been created for assisting with data migration
