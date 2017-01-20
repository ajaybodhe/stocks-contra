package core

import (
//	"time"
	"fmt"
	"github.com/ajaybodhe/stocks-contra/algo"
	//"github.com/ajaybodhe/stocks-contra/api"
	_ "github.com/go-sql-driver/mysql"
	//"net/http"
		//"crypto/tls"
	"github.com/ajaybodhe/stocks-contra/api"
	"github.com/ajaybodhe/stocks-contra/util"
	"time"
	"sync"
)

//var client *http.Client
//var proddbhandle util.DB
//var testdbhandle util.DB
//
//func initDB() {
//	//initialize production db handle
//	fmt.Printf("\nconf.StocksContraConfig.DB.ConnID=%v",conf.StocksContraConfig.DB.ConnID)
//	proddb, err := sql.Open("mysql", conf.StocksContraConfig.DB.ConnID + "&parseTime=True")
//	if err != nil {
//		glog.Errorln("error: connecting to mysql:", conf.StocksContraConfig.DB.ConnID, ":error:", err.Error())
//		return
//	}
//	if err := proddb.Ping(); err != nil {
//		glog.Fatalln("fatal: unable to connect to db:", err.Error())
//		os.Exit(1)
//	}
//	proddbhandle.Set(proddb)
//}

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
// no stock by strat 1
// way to minimize output of algo, bloom or compare with NSE first 1k stocks
// concurrent news reading
// sentiment analysis of the news

func Serve() {
	/* TBD AJAY
	decide upon the structure of code,
	parallelise api calls
	few ratios missing: debt/equity, roe, roce, roa
	*/
	//var err error
	//initDB()
	//tr := &http.Transport{
     //   TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	//client := &http.Client{Transport: tr}

	///* Call to this function depends on passed argument */
	//api.GetNSESectoralIndexLists(client)
	//api.GetNSEBroadMarketIndexLists(client)
	//
	///*getNSEDeliveryPercentageData(5)*/
	//api.GetNSESecuritiesFullBhavData(client, false)
	//
	//err = api.FetchNStoreMoneyControlData(client)
	//if err != nil {
	//	fmt.Println("FetchNStoreMoneyControlData failed")
	//}
	//
	//err = algo.NSESecuritiesBuySignal()
	//if err != nil {
	//	fmt.Println("NSESecuritiesBuySignal failed")
	//}

	//fmt.Printf("%v", api.GetNSELiveQuote(client, "ABB"))
	var wg sync.WaitGroup
	wg.Add(1)
	go algo.NseOrderBookAnalyser(&wg)
	
	
	for {
		wg.Add(1)
		err = api.GetNseCorporateAnnouncements(&wg, client, proddbhandle, util.NSECorporateAnnounceMentLink)
		if err != nil {
			fmt.Println("GetNseCorporateAnnouncements failed")
		}
		time.Sleep(time.Second * 10)
	}
//	err := api.GetBseCorporateAnnouncements(client, proddbhandle, util.BSECorporateAnnounceMentLink)
//	if err != nil {
//		fmt.Println("GetBseCorporateAnnouncements failed")
//	}
}
