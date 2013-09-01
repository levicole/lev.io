package models

import (
    "log"
    "os"
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/coopernurse/gorp"
)

var (
    db      *sql.DB
    dbmap   *gorp.DbMap
    Links   *gorp.TableMap
    Clients *gorp.TableMap
    NewBase60Vocab [60]string
)

func setupDbConn() {
    NewBase60Vocab = [60]string{ "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "_", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z" }
    var mysqlerr error
    db, mysqlerr = sql.Open("mysql", "root@unix(/opt/boxen/data/mysql/socket)/levio")
    if mysqlerr != nil {
        fmt.Printf("some error:%s\n", mysqlerr.Error())
        panic(mysqlerr)
    }

    dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
    dbmap.AddTableWithName(Link{},   "links").SetKeys(true, "Id")

    dbmap.CreateTablesIfNotExists()
    dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
}

func init() {
    setupDbConn()
}
