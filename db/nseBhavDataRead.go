package db

import (
	"database/sql"
	//"database/sql/driver"
	"fmt"
	"github.com/ajaybodhe/stocks-contra/coreStructures"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"time"
)

const allSecurityNseBhavDataSql = "select symbol, date, prev_close, open_price, high_price, low_price, last_price, close_price, avg_price, ttl_trd_qnty, deliv_qty, deliv_per from  NSESecuritiesFullBhavData order by symbol, date;"

func ReadNseBhavData(securitySymbols []string) (map[string][]coreStructures.NseBhavRecord, error) {
	var rows *sql.Rows
	var err error

	if len(securitySymbols) <= 0 {
		rows, err = proddbhandle.Query(allSecurityNseBhavDataSql)
		fmt.Println("Query NSE:", allSecurityNseBhavDataSql)
	} else {
		firstTime := true
		specificSecurityNseBhavDataSql := "select symbol, date, prev_close, open_price, high_price, low_price, last_price, close_price, avg_price, ttl_trd_qnty, deliv_qty, deliv_per from  NSESecuritiesFullBhavData where symbol in ("
		for _, securitySymbol := range securitySymbols {
			if firstTime == true {
				specificSecurityNseBhavDataSql = specificSecurityNseBhavDataSql + "\"" + securitySymbol + "\""
				firstTime = false
			} else {
				specificSecurityNseBhavDataSql = specificSecurityNseBhavDataSql + ", \"" + securitySymbol + "\""
			}
		}
		specificSecurityNseBhavDataSql = specificSecurityNseBhavDataSql + ") order by symbol, date;"
		rows, err = proddbhandle.Query(specificSecurityNseBhavDataSql)
		fmt.Println("query NSE:", specificSecurityNseBhavDataSql)
	}

	if err != nil {
		return nil, fmt.Errorf("fetch security symbol nse bhav data err: sql error:%s\n", err.Error())
	}
	defer rows.Close()

	var symbol string
	var symbol1 string
	var sDate string
	var dateFormat = "2006-01-02"
	var dDate time.Time
	var nbr coreStructures.NseBhavRecord
	nbra := make([]coreStructures.NseBhavRecord, 0)

	filled := false
	nbrm := make(map[string][]coreStructures.NseBhavRecord)

	for rows.Next() {
		err = rows.Scan(&symbol, &sDate, &nbr.PrevClosePrice, &nbr.OpenPrice, &nbr.HighPrice, &nbr.LowPrice, &nbr.LastPrice, &nbr.ClosePrice, &nbr.AvgPrice, &nbr.TtlTrdQnty, &nbr.DelivQty, &nbr.DelivPer)
		if err != nil {
			glog.Error("error: while reading security :error:%s", err.Error())
			return nil, err
		}
		if symbol != symbol1 && filled == true {
			nbrm[symbol1] = nbra
			nbra = make([]coreStructures.NseBhavRecord, 0)
		} //else {
		filled = true
		symbol1 = symbol

		if sDate != "" {
			dDate, err = time.Parse(dateFormat, sDate)
			if err != nil {
				glog.Info(err)
			}
		}

		nbr.RecordDate = dDate
		//fmt.Printf("symbol: %v, \nnbr: %v", symbol, nbr)
		nbra = append(nbra, nbr)
		//}
	}
	if filled == true {
		nbrm[symbol1] = nbra
	}

	//for k, v := range nbrm {
	//	for k1, v1 := range v {
	//		fmt.Printf("\n%v	%v	%v", k, k1, v1)
	//	}
	//}

	return nbrm, nil
}

/*
func readNseBhavData(db util.DB, securitySymbolSql []string) (map[string]*coreStructures.NseBhavData, error) {
	var rows *sql.Rows
	var err error

	if len(securitySymbolSql) <= 0 {
		rows, err = db.Query(allSecurityNseBhavDataSql)
	} else {
		specificSecurityNseBhavDataSql := ""
		rows, err = db.Query(specificSecurityNseBhavDataSql)
	}

	if err != nil {
		return nil, fmt.Errorf("fetch security symbol nse bhav data err: sql error:%s\n", err.Error())
	}
	defer rows.Close()

	var symbol string
	var sDate string
	var dateFormat = "2006-01-02"
	var dDate time.Time
	var nbr coreStructures.NseBhavRecord

	filled := false
	nbdm := make(map[string]*coreStructures.NseBhavData)
	nbd := coreStructures.NewNseBhavData(0)

	for rows.Next() {
		err = rows.Scan(&symbol, &sDate, &nbr.PrevClosePrice, &nbr.OpenPrice, &nbr.HighPrice, &nbr.LowPrice, &nbr.LastPrice, &nbr.ClosePrice, &nbr.AvgPrice, &nbr.TtlTrdQnty, &nbr.DelivQty, &nbr.DelivPer)
		if err != nil {
			glog.Error("error: while reading security :error:%s", err.Error())
			return nil, err
		}
		if symbol != nbd.Symbol && filled == true {
			nbdm[nbd.Symbol] = nbd
			nbd := coreStructures.NewNseBhavData(0)
			nbd.Symbol = symbol
		} else {
			filled = true
			nbd.Symbol = symbol

			if sDate != "" {
				dDate, err = time.Parse(dateFormat, sDate)
				if err != nil {
					glog.Info(err)
				}
			}

			nbr.RecordDate = dDate
			nbd.BhavRecord = append(nbd.BhavRecord, nbr)
		}
	}
	if filled == true {
		nbdm[nbd.Symbol] = nbd
	}

	return nbdm, nil
}
*/
