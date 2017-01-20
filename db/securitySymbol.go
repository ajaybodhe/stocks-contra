package db

import (
	"fmt"
	"github.com/ajaybodhe/stocks-contra/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/ajaybodhe/stocks-contra/coreStructures"
)

const securitySymbolSql = "select distinct symbol from NSESecuritiesFullBhavData order by symbol"

func GetSecuritySymbols() ([]string, error) {
	var securitySymbols []string

	rows, err := proddbhandle.Query(securitySymbolSql)
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

const bloomSecuritySymbolSql = "select distinct symbol from %v order by symbol"
func GetInterestedSymbolsBloom() (*coreStructures.BloomFilter, error) {
	
	bloomFilter := coreStructures.NewBloomFilter(util.BloomBits, util.BloomHashCount)
	
	for k,_ := range util.NSEBroadMarketIndexList {
	
		query := fmt.Sprintf(bloomSecuritySymbolSql, k)
	
		fmt.Printf("\nquery=%v\n", query)
		
		rows, err := proddbhandle.Query(query)
		if err != nil {
			return nil, fmt.Errorf("fetch security symbols for bloom: sql error:%s\n", err.Error())
		}
		defer rows.Close()
	
		var securitySymbol string
		
		for rows.Next() {
			if err = rows.Scan(&securitySymbol); err != nil {
				glog.Error("error: while reading input in bloom security symbol fetch:error:%s", err.Error())
				return nil, err
			}
			bloomFilter.Add([]byte(securitySymbol))
		}
	}
	return bloomFilter, nil
}

const interestedSecuritySymbolsSql = "select distinct symbol from %v order by symbol"
func GetInterestedSymbols(tables []string) ([]string, error) {
	
	var symbols []string

	for k,_ := range tables {
		
		query := fmt.Sprintf(interestedSecuritySymbolsSql, tables[k])
		
		fmt.Printf("\nquery=%v\n", query)
		
		rows, err := proddbhandle.Query(query)
		if err != nil {
			return nil, fmt.Errorf("fetch interested security symbols: sql error:%s\n", err.Error())
		}
		defer rows.Close()
		
		var securitySymbol string
		
		for rows.Next() {
			if err = rows.Scan(&securitySymbol); err != nil {
				glog.Error("error: while reading input in interested security symbol fetch:error:%s", err.Error())
				return nil, err
			}
			symbols = append(symbols,securitySymbol)
		}
	}
	return symbols, nil
}