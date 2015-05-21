package db

import (
	"fmt"
	"github.com/ajaybodhe/stocks-contra/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

const securitySymbolSql = "select distinct symbol from NSESecuritiesFullBhavData order by symbol"

func GetSecuritySymbols(db util.DB) ([]string, error) {
	var securitySymbols []string

	rows, err := db.Query(securitySymbolSql)
	if err != nil {
		return nil, fmt.Errorf("fetch security symbols: sql error:%s\n", err.Error())
	}
	defer rows.Close()

	var securitySymbol string
	found := false

	for rows.Next() {
		if err = rows.Scan(&securitySymbol); err != nil {
			glog.Error("error: while reading input in securit symbol fetch:error:%s", err.Error())
			return nil, err
		}
		securitySymbols = append(securitySymbols, securitySymbol)
		found = true

	}
	if !found {
		return nil, fmt.Errorf("No security symbol information found")
	}
	return securitySymbols, nil
}
