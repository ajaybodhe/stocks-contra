package db

import (
	//"database/sql"
	"fmt"
	"github.com/ajaybodhe/stocks-contra/coreStructures"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"time"
)

func WriteNSESecuritiesBuySignal(nlsd map[string]coreStructures.NseSecurityLongSignalData) error {
	sqlQ := "insert into NseSecurityLongSignalData (date, symbol, sector, pe, industry_pe, correction, closeness52weeklow, deliv_per, strategy_id) values "
	firstTime := true
	currentTime := time.Now().Local()
	today := currentTime.Format("2006-01-02")
	for k, v := range nlsd {
		if firstTime {
			sqlQ = sqlQ + fmt.Sprintf(" (\"%v\", \"%v\", \"%v\", %v, %v, %v, %v, %v, %v)", today, k, v.Sector, v.PE, v.IndustryPE, v.Correction, v.Closeness52WeekLow, v.DelivPer, v.Strategy)
			firstTime = false
		} else {
			sqlQ = sqlQ + fmt.Sprintf(", (\"%v\", \"%v\", \"%v\", %v, %v, %v, %v, %v, %v)", today, k, v.Sector, v.PE, v.IndustryPE, v.Correction, v.Closeness52WeekLow, v.DelivPer, v.Strategy)
		}
	}
	sqlQ = sqlQ + " ON DUPLICATE KEY UPDATE date=VALUES(date), symbol=VALUES(symbol), sector=VALUES(sector), pe=values(pe), industry_pe=VALUES(industry_pe), correction=VALUES(correction), closeness52weeklow=VALUES(closeness52weeklow), deliv_per=VALUES(deliv_per), strategy_id=VALUES(strategy_id);"
	fmt.Printf("the query is: %s", sqlQ)
	_, err := proddbhandle.Exec(sqlQ)
	if err != nil {
		return fmt.Errorf("WriteNSESecuritiesBuySignal: sql error:%s\n", err.Error())
	}

	err = deleteOldNSESecuritiesBuySignalData()
	if err != nil {
		return fmt.Errorf("deleteOldNSESecuritiesBuySignalData: sql error:%s\n", err.Error())
	}
	return nil
}

func deleteOldNSESecuritiesBuySignalData() error {
	deleteSql := "delete  from NseSecurityLongSignalData where date in (select minDate from (select min(date) as minDate from NseSecurityLongSignalData) as X) and exists (select count from (select if(count(distinct date)> 5, count(distinct date), 0) as count from NseSecurityLongSignalData) as Y where count > 0 );"
	_, err := proddbhandle.Exec(deleteSql)
	if err != nil {
		return fmt.Errorf("deleteOldNSESecuritiesBuySignalData: sql error:%s\n", err.Error())
	}
	return nil
}

func RetrieveAllSymbolsNStrategy() (map[string]int, error) {
	sqlQ := "select Nse.symbol, IFNULL(Trade.strategy_id,0) from (select distinct symbol from NSESecuritiesFullBhavData) AS Nse LEFT OUTER JOIN (select * from NseSecurityLongSignalData where date = (select date from (select max(date) from NseSecurityLongSignalData) AS temp)) AS Trade on Nse.symbol=Trade.symbol;"
	rows, err := proddbhandle.Query(sqlQ)
	if err != nil {
		return nil, fmt.Errorf("RetrieveAllSymbolsNStrategy() details err: sql error:%s\n", err.Error())
	}
	defer rows.Close()
	symbolStrategyMap := make(map[string]int)
	var symbol string
	var strategyId int
	found := false
	for rows.Next() {
		err = rows.Scan(&symbol, &strategyId)
		if err != nil {
			glog.Error("error: while reading security :error:%s", err.Error())
			return nil, err
		}
		found = true
		symbolStrategyMap[symbol] = strategyId
	}
	if found {
		return symbolStrategyMap, nil
	}
	return nil, fmt.Errorf("no data for security symbols")
}
