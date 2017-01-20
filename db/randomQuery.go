package db

import "database/sql"

func RandomQuery(q string) (*sql.Rows, error){
	return proddbhandle.Query(q)
}
