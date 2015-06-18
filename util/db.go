package util

import "database/sql"

/*DB is a container which holds mysql connection */
type DB struct {
	*sql.DB
}

/*Set sets the dbhandle into DB structure */
func (d *DB) Set(dbhandle *sql.DB) {
	d.DB = dbhandle
}
