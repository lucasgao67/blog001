package util

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var SqlClient *sqlx.DB

func SqlInit() {
	SqlClient = sqlx.MustOpen("mysql", "root:@(localhost:4000)/test?parseTime=true")
}