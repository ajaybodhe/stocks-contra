package core

import (
	"database/sql"
	"net/http"
	"os"
	"strings"

	"github.com/ajaybodhe/stocks-contra/conf"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/pubmatic/pub-phoenix/cfiller/util"
)

const (
	NSEAutoIndex     = "http://www.nse-india.com/content/indices/ind_cnxautolist.csv"
	NSECNX500        = "http://www.nseindia.com/products/content/equities/indices/cnx_500.htm"
	forwardSlashChar = "/"
	fileDownloadPath = "/tmp/"
)

var NSEIndexList = map[string]string{
	"AUTO":     "http://www.nse-india.com/content/indices/ind_cnxautolist.csv",
	"BANK":     "http://www.nse-india.com/content/indices/ind_cnxbanklist.csv",
	"ENERGY":   "http://www.nse-india.com/content/indices/ind_cnxenergylist.csv",
	"FINANCE":  "http://www.nse-india.com/content/indices/ind_cnxfinancelist.csv",
	"FMCG":     "http://www.nse-india.com/content/indices/ind_cnxfmcglist.csv",
	"IT":       "http://www.nse-india.com/content/indices/ind_cnxitlist.csv",
	"MEDIA":    "http://www.nse-india.com/content/indices/ind_cnxmedialist.csv",
	"METAL":    "http://www.nse-india.com/content/indices/ind_cnxmetallist.csv",
	"PHARMA":   "http://www.nse-india.com/content/indices/ind_cnxpharmalist.csv",
	"PSU_BANK": "http://www.nse-india.com/content/indices/ind_cnxpsubanklist.csv",
	"REALTY":   "http://www.nse-india.com/content/indices/ind_cnxrealtylist.csv",
	"INDUSTRY": "http://www.nse-india.com/content/indices/ind_cnxindustrylist.csv",
}

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

func getNSEIndexLists() {
	glog.Println("============Getting NSE Indices along with Listed Comapnies==============")
	client = &http.Client{}
	for key, value := range NSEIndexList {
		glog.Println(key, value)
		resp, err := client.Get(value)
		if err != nil {
			glog.Println(":Result:Fail:Error:", err.Error())
			continue
		}
		defer resp.Body.Close()
		filePath := strings.Split(key, forwardSlahCHar)
		path := fileDownloadPath + filePath[len(filePath)-1]
		glog.Println(resp.Status)
		//log.Println(resp)
	}
}

func Serve() {
	initDB()
	/* Call to this function depends on passed argument */
	getNSEIndexLists()
}
