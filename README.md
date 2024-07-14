# go-filterparams-gopg

[go-pg](https://github.com/go-pg/pg) query generator for use with [go-filterparams](https://github.com/cbrand/go-filterparams).

Currently in an experimental state. Do expect bugs!

## Usage

To use, simply create an instance of `Parser` and call the `AppendTo` function with your original query.

```go
package main

import (
    "github.com/cbrand/go-filterparams"
    "github.com/go-pg/pg/v10"
    "github.com/iancoleman/strcase"
    "github.com/ezraisw/go-filterparams-gopg"
)

type User struct {
    tableName interface{} `pg:"users"`
    ID        string      `pg:",notnull,use_zero"`
    Username  string      `pg:",notnull,use_zero"`
    Password  string      `pg:",notnull,use_zero"`
}

func main() {
    var db *pg.DB

    //... obtain your go-pg query
    var users []*User
    pgQuery := db.Model(&users)

    //... obtain your query data
    var queryData *filterparams.QueryData

    // create a new parser
    parser := fpgopg.NewParser(fpgopg.NamingFunc(strcase.ToSnake))

    // append the parsed filterparams.QueryData
    pgQuery = parser.AppendTo(pgQuery, queryData)

    // perform select operation
    pgQuery.Select()
}
```
