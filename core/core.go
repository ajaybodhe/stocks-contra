package core

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ajaybodhe/stocks-contra/conf"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/pubmatic/pub-phoenix/cfiller/util"
)

const (
	NSECNX500        = "http://www.nseindia.com/products/content/equities/indices/cnx_500.htm"
	forwardSlashChar = "/"
	fileDownloadPath = "/tmp/"
)

var NSESectoralIndexList = map[string]string{
	"AUTO":     "http://www.nseindia.com/content/indices/ind_cnxautolist.csv",
	"BANK":     "http://www.nseindia.com/content/indices/ind_cnxbanklist.csv",
	"ENERGY":   "http://www.nseindia.com/content/indices/ind_cnxenergylist.csv",
	"FINANCE":  "http://www.nseindia.com/content/indices/ind_cnxfinancelist.csv",
	"FMCG":     "http://www.nseindia.com/content/indices/ind_cnxfmcglist.csv",
	"IT":       "http://www.nseindia.com/content/indices/ind_cnxitlist.csv",
	"MEDIA":    "http://www.nseindia.com/content/indices/ind_cnxmedialist.csv",
	"METAL":    "http://www.nseindia.com/content/indices/ind_cnxmetallist.csv",
	"PHARMA":   "http://www.nseindia.com/content/indices/ind_cnxpharmalist.csv",
	"PSU_BANK": "http://www.nseindia.com/content/indices/ind_cnxpsubanklist.csv",
	"REALTY":   "http://www.nseindia.com/content/indices/ind_cnxrealtylist.csv",
	"INDUSTRY": "http://www.nseindia.com/content/indices/ind_cnxindustrylist.csv",
}

var NSEBroadMarketIndexList = map[string]string{
	"CNX_NIFTY":       "http://www.nseindia.com/content/indices/ind_niftylist.csv",
	"CNX_NIFT_JUNIOR": "http://www.nseindia.com/content/indices/ind_jrniftylist.csv",
	"CNX_100":         "http://www.nseindia.com/content/indices/ind_cnx100list.csv",
	"CNX_200":         "http://www.nseindia.com/content/indices/ind_cnx200list.csv",
	"CNX_500":         "http://www.nseindia.com/content/indices/ind_cnx500list.csv",
	"NIFT_MIDCAP_50":  "http://www.nseindia.com/content/indices/ind_niftymidcap50list.csv",
	"CNX_MIDCAP":      "http://www.nseindia.com/content/indices/ind_cnxmidcaplist.csv",
	"CNX_SMALLCAP":    "http://www.nseindia.com/content/indices/ind_cnxsmallcap.csv",
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

func getNSEBroadMarketIndexList() {
	glog.Infoln("============Getting NSE Broad Market Indices along with Listed Comapnies==============")
	/* TBD AJAY req/resp/client which objects should be created outside loop?*/

	for key, value := range NSEBroadMarketIndexList {
		glog.Infoln(key, value)

		req, err := http.NewRequest("GET", value, nil)
		if err != nil {
			glog.Fatalln(err)
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:28.0) Gecko/20100101 Firefox/28.0")
		req.Header.Set("Host", "www.nseindia.com")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Accept-Encoding", "gzip, deflate")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

		resp, err := client.Do(req)
		if err != nil {
			glog.Errorln(":Result:Fail:Error:", err.Error())
			continue
		}
		//defer resp.Body.Close()

		filePath := strings.Split(value, forwardSlashChar)
		path := fileDownloadPath + filePath[len(filePath)-1]
		glog.Infoln(resp.Status)

		file, err := os.Create(path)
		if err != nil {
			glog.Errorln(err)
		}
		//defer file.Close()

		size, err := io.Copy(file, resp.Body)
		if err != nil {
			glog.Errorln(err)
		}
		glog.Infoln("%s with %v bytes downloaded", path, size)
		//log.Println(resp)
		resp.Body.Close()
		file.Close()
		req.Close
	}
}

func getNSESectoralIndexLists() {
	glog.Infoln("============Getting NSE Sectoral Indices along with Listed Comapnies==============")
	/* TBD AJAY req/resp/client which objects should be created outside loop?*/

	for key, value := range NSESectoralIndexList {
		glog.Infoln(key, value)

		req, err := http.NewRequest("GET", value, nil)
		if err != nil {
			glog.Fatalln(err)
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:28.0) Gecko/20100101 Firefox/28.0")
		req.Header.Set("Host", "www.nseindia.com")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Accept-Encoding", "gzip, deflate")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

		resp, err := client.Do(req)
		if err != nil {
			glog.Errorln(":Result:Fail:Error:", err.Error())
			continue
		}
		//defer resp.Body.Close()

		filePath := strings.Split(value, forwardSlashChar)
		path := fileDownloadPath + filePath[len(filePath)-1]
		glog.Infoln(resp.Status)

		file, err := os.Create(path)
		if err != nil {
			glog.Errorln(err)
		}
		//defer file.Close()

		size, err := io.Copy(file, resp.Body)
		if err != nil {
			glog.Errorln(err)
		}
		glog.Infoln("%s with %v bytes downloaded", path, size)
		//log.Println(resp)
		resp.Body.Close()
		file.Close()
		req.Close
	}
}

func Serve() {
	initDB()
	client = &http.Client{}
	/* Call to this function depends on passed argument */
	getNSESectoralIndexLists()
	getNSEBroadMarketIndexList()
}
