package db

import (
	"database/sql"
	"fmt"
	"github.com/ajaybodhe/stocks-contra/coreStructures"
	"github.com/ajaybodhe/stocks-contra/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

const allSecurityDetailsSql = "select symbol, sector, high52, low52, eps, pe, industry_pe, market_cap, book_value, dividend, pb, pc, face_value, div_yield, promoter_holding, fii_holding, dii_holding, other_holding from MoneyControlSecurityDetails"

func ReadSecurityDetails(db util.DB, securitySymbol string) (*coreStructures.MoneyControlSecurityStructure, error) {

	securityDetailsSql := "select symbol, sector, high52, low52, eps, pe, industry_pe, market_cap, book_value, dividend, pb, pc, face_value, div_yield, promoter_holding, fii_holding, dii_holding, other_holding from MoneyControlSecurityDetails where symbol=" + "\"" + securitySymbol + "\""
	rows, err := db.Query(securityDetailsSql)
	if err != nil {
		return nil, fmt.Errorf("fetch security symbol=%s details err: sql error:%s\n", securitySymbol, err.Error())
	}
	defer rows.Close()

	var mcss coreStructures.MoneyControlSecurityStructure
	var symbol string

	if rows.Next() {
		err = rows.Scan(&symbol, &mcss.Sector, &mcss.High52, &mcss.Low52, &mcss.EPS, &mcss.PE, &mcss.IndustryPE, &mcss.MarketCap, &mcss.BookValue, &mcss.Dividend, &mcss.PB, &mcss.PC, &mcss.FaceValue, &mcss.DivYield, &mcss.PromoterHolding, &mcss.DIIHolding, &mcss.FIIHolding, &mcss.OtherHolding)
		if err != nil {
			glog.Error("error: while reading security :error:%s", err.Error())
			return nil, err
		}
		return &mcss, nil
	}

	return nil, fmt.Errorf("no data for security symbol &s", securitySymbol)
}

func ReadAllSecurityDetails(db util.DB, securitySymbols []string) (map[string]coreStructures.MoneyControlSecurityStructure, error) {
	var rows *sql.Rows
	var err error

	if len(securitySymbols) <= 0 {
		rows, err = db.Query(allSecurityDetailsSql)
	} else {
		firstTime := true
		specificSecurityDetailsSql := "select symbol, sector, high52, low52, eps, pe, industry_pe, market_cap, book_value, dividend, pb, pc, face_value, div_yield, promoter_holding, fii_holding, dii_holding, other_holding from MoneyControlSecurityDetails where symbol in ("
		for _, securitySymbol := range securitySymbols {
			if firstTime == true {
				specificSecurityDetailsSql = specificSecurityDetailsSql + "\"" + securitySymbol + "\""
				firstTime = false
			} else {
				specificSecurityDetailsSql = specificSecurityDetailsSql + ", \"" + securitySymbol + "\""
			}
		}
		specificSecurityDetailsSql = specificSecurityDetailsSql + ");"
		rows, err = db.Query(specificSecurityDetailsSql)
	}
	if err != nil {
		return nil, fmt.Errorf("fetch security symbols details err: sql error:%s\n", err.Error())
	}
	defer rows.Close()

	var mcss coreStructures.MoneyControlSecurityStructure
	var symbol string
	mcssCollection := make(map[string]coreStructures.MoneyControlSecurityStructure)
	var found bool

	for rows.Next() {
		err = rows.Scan(&symbol, &mcss.Sector, &mcss.High52, &mcss.Low52, &mcss.EPS, &mcss.PE, &mcss.IndustryPE, &mcss.MarketCap, &mcss.BookValue, &mcss.Dividend, &mcss.PB, &mcss.PC, &mcss.FaceValue, &mcss.DivYield, &mcss.PromoterHolding, &mcss.DIIHolding, &mcss.FIIHolding, &mcss.OtherHolding)
		if err != nil {
			glog.Error("error: while reading security :error:%s", err.Error())
			return nil, err
		}
		mcssCollection[symbol] = mcss
		found = true
	}
	if found {
		return mcssCollection, nil
	}

	return nil, fmt.Errorf("no data for security symbols")

}
