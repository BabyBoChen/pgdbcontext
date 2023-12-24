# PgDbContext
 A Postgresql Data Access Model Written in Golang

## Get started
~~~
import (
    ...
    "github.com/BabyBoChen/pgdbcontext"
)
...
//1. create connection
var db *pgdbcontext.DbContext
var err error
db, err = pgdbcontext.NewDbContext("host=db.server.com port=5432 dbname=postgres user=username password=secretpassword sslmode=require")

//2. select data
var dt *pgdbcontext.DataTable
dt, err = db.Query("SELECT * FROM table_name WHERE col=$1 OR col=$2", "param1", 9999)

//3. insert data
var repo *pgdbcontext.DbRepository
repo, err = db.GetRepository("table_name")
var lastInsertedId map[string]interface{}
lastInsertedId, err = repo.Insert(map[string]interface{}{
    "col": "some value",
    "col2": 9999,
    ...
})

//4. update data
err = repo.Update(map[string]interface{}{
    "pk_column": "must contains pk column",
    "col2": 10000,
    ...
})

//5. delete data
err = repo.Delete(map[string]interface{}{
    "pk_column": "must contains pk column",
})

//6. please don't forget to commit the transaction
db.Commit()
~~~
For more functionalities, please take a look at dbcontext.go.
