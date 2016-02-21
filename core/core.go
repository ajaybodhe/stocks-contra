package core

import (
	"database/sql"
	"fmt"
	//"github.com/ajaybodhe/stocks-contra/algo"
	"github.com/ajaybodhe/stocks-contra/api"
	"github.com/ajaybodhe/stocks-contra/conf"
	//"github.com/ajaybodhe/stocks-contra/coreStructures"
	"github.com/ajaybodhe/stocks-contra/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"net/http"
	"os"
	/*
		"golang.org/x/text/encoding"
		"golang.org/x/text/encoding/charmap"
		"golang.org/x/text/transform"
	*/)

var client *http.Client
var proddbhandle util.DB
var testdbhandle util.DB

func initDB() {
	//initialize production db handle
	proddb, err := sql.Open("mysql", conf.StocksContraConfig.DB.ConnID)
	if err != nil {
		glog.Errorln("error: connecting to mysql:", conf.StocksContraConfig.DB.ConnID, ":error:", err.Error())
		return
	}
	if err := proddb.Ping(); err != nil {
		glog.Fatalln("fatal: unable to connect to db:", err.Error())
		os.Exit(1)
	}
	proddbhandle.Set(proddb)
}

/* TBD AJAY */
// RATIOOOOOS & NEWWWWWWS - edelwiess
// Mutual fund activity, promoters increasing stakes
// fut n options - strategies
//twitter feeds, news feeds
// rakesh jhunjhunwalla site - HNI holding, how order book works
// correction in fav stocks
// measure performace of algorithm - are suggested stocks good as well, should feature in suggestion n if suggested how they performed
// future market/ oi at put call option
// simple moving averages, other ratios,
// bulk deals, block deals, short selling
// instead of 3/5/8 day avg calculate when the delivery is going high

func Serve() {
	/* TBD AJAY
	decide upon the structure of code,
	parallelise api calls
	few ratios missing: debt/equity, roe, roce, roa
	*/
	var err error
	initDB()
	client = &http.Client{}

	/* Call to this function depends on passed argument */
//	api.GetNSESectoralIndexLists(client, proddbhandle)
//	api.GetNSEBroadMarketIndexLists(client, proddbhandle)

//	err = api.FetchNStoreMoneyControlData(client, proddbhandle)
//	if err != nil {
//		fmt.Println("FetchNStoreMoneyControlData failed")
//	}

	/*getNSEDeliveryPercentageData(5)*/
//	api.GetNSESecuritiesFullBhavData(client, proddbhandle, false)

//	err = algo.NSESecuritiesBuySignal(proddbhandle)
//	if err != nil {
//		fmt.Println("NSESecuritiesBuySignal failed")
//	}

	//fmt.Printf("%v", api.GetNSELiveQuote(client, "ABB"))
//	err = algo.NseOrderBookAnalyser(client, proddbhandle)
//	if err != nil {
//		fmt.Println("NseOrderBookAnalyser failed")
//	}
	
	err = api.GetNseCorporateAnnouncements(client, proddbhandle, util.NSECorporateAnnounceMentLink)
	if err != nil {
		fmt.Println("GetNseCorporateAnnouncements failed")
	}
}
