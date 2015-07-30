package db

import (
	"fmt"
	"github.com/ajaybodhe/stocks-contra/coreStructures"
	"github.com/ajaybodhe/stocks-contra/util"
	_ "github.com/go-sql-driver/mysql"
)

/* TBD AJAY try second option from
http://stackoverflow.com/questions/21108084/golang-mysql-insert-multiple-data-at-once*/

//const securityDetailsWriteSql = "insert into MoneyControlSecurityDetails (symbol, eps, pe, industry_pe, market_cap, face_value, book_value, dividend, pb, pc, fv, div_yeild, promoter_holding, fii_holding, dii_holding, other_holding) values "

/* what is difference betn db.Exec n db.Query */
func WriteSecurityDetails(db util.DB, mcss map[string]*coreStructures.MoneyControlSecurityStructure) error {
	securityDetailsWriteSql := "insert into MoneyControlSecurityDetails (symbol, sector, high52, low52, eps, pe, industry_pe, market_cap, book_value, dividend, pb, pc, face_value, div_yield, promoter_holding, fii_holding, dii_holding, other_holding) values "
	firstTime := true
	for k, v := range mcss {
		if firstTime {
			securityDetailsWriteSql = securityDetailsWriteSql + fmt.Sprintf("(\"%v\", \"%v\", %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v)", k, v.Sector, v.High52, v.Low52, v.EPS, v.PE, v.IndustryPE, v.MarketCap, v.BookValue, v.Dividend, v.PB, v.PC, v.FaceValue, v.DivYield, v.PromoterHolding, v.DIIHolding, v.FIIHolding, v.OtherHolding)
			firstTime = false
		} else {
			securityDetailsWriteSql = securityDetailsWriteSql + fmt.Sprintf(",(\"%v\", \"%v\", %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v)", k, v.Sector, v.High52, v.Low52, v.EPS, v.PE, v.IndustryPE, v.MarketCap, v.BookValue, v.Dividend, v.PB, v.PC, v.FaceValue, v.DivYield, v.PromoterHolding, v.DIIHolding, v.FIIHolding, v.OtherHolding)
		}
	}
	securityDetailsWriteSql = securityDetailsWriteSql + " ON DUPLICATE KEY UPDATE symbol=VALUES(symbol), sector=VALUES(sector), high52=VALUES(high52), low52=VALUES(low52), eps=VALUES(eps), pe=values(pe), industry_pe=VALUES(industry_pe), market_cap=VALUES(market_cap), book_value=VALUES(book_value), dividend=VALUES(dividend), pb=VALUES(pb), pc=VALUES(pc), face_value=VALUES(face_value), div_yield=VALUES(div_yield), promoter_holding=VALUES(promoter_holding), fii_holding=VALUES(fii_holding), dii_holding=VALUES(dii_holding), other_holding=VALUES(other_holding)"
	fmt.Printf("the query is: %s", securityDetailsWriteSql)
	_, err := db.Exec(securityDetailsWriteSql)
	if err != nil {
		return fmt.Errorf("fetch security symbols: sql error:%s\n", err.Error())
	}
	return nil
}
