package db

import (
	//"database/sql"
	"fmt"
	"github.com/ajaybodhe/stocks-contra/coreStructures"
	"github.com/ajaybodhe/stocks-contra/util"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func WriteNSESecuritiesBuySignal(db util.DB, nlsd map[string]coreStructures.NseSecurityLongSignalData) error {
	sqlQ := "insert into NseSecurityLongSignalData (date, symbol, pe, industry_pe, correction, closeness52weeklow, deliv_per, sector) values "
	firstTime := true
	currentTime := time.Now().Local()
	today := currentTime.Format("2006-01-02")
	for k, v := range nlsd {
		if firstTime {
			sqlQ = sqlQ + fmt.Sprintf(" (%v, \"%v\", \"%v\", %v, %v, %v, %v, %v, %v)", today, k, v.Sector, v.PE, v.IndustryPE, v.Correction, v.Closeness52WeekLow, v.DelivPer)
			firstTime = false
		} else {
			sqlQ = sqlQ + fmt.Sprintf(", (%v, \"%v\", \"%v\", %v, %v, %v, %v, %v, %v)", today, k, v.Sector, v.PE, v.IndustryPE, v.Correction, v.Closeness52WeekLow, v.DelivPer)
		}
	}
	sqlQ = sqlQ + " ON DUPLICATE KEY UPDATE date=VALUES(date), symbol=VALUES(symbol), sector=VALUES(sector), pe=values(pe), industry_pe=VALUES(industry_pe), correction=VALUES(correction), closeness52WeekLow=VALUES(closeness52WeekLow), delivPer=VALUES(delivPer)"
	fmt.Printf("the query is: %s", sqlQ)
	_, err := db.Exec(sqlQ)
	if err != nil {
		return fmt.Errorf("WriteNSESecuritiesBuySignal: sql error:%s\n", err.Error())
	}
	return nil
}
