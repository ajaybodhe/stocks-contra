package core

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ajaybodhe/stocks-contra/conf"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/pubmatic/pub-phoenix/cfiller/util"
	/*
		"golang.org/x/text/encoding"
		"golang.org/x/text/encoding/charmap"
		"golang.org/x/text/transform"
	*/)

const (
	forwardSlashChar   = "/"
	fileDownloadPath   = "/tmp/"
	truncateTableQuery = "truncate table %s"
	loadFileQuery      = "LOAD DATA LOCAL INFILE '%s' INTO TABLE %s FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n' IGNORE %d ROWS;"

	createTableQuery           = "CREATE TABLE IF NOT EXISTS `%s` ( `record_type` INTEGER(4), `sr_no` INTEGER(4), `symbol` varchar (200), `security_type` varchar(10), `traded_quantity` INTEGER(20), `deliverable_quantity` INTEGER(20), delivery_percentage double, PRIMARY KEY (`symbol`)) ENGINE=InnoDB DEFAULT CHARSET=utf8 ;"
	NSEDeliveryPercentageTable = "NSEDeliveryData_%02d%02d%04d"

	loadFileQueryNSEFBD           = "LOAD DATA LOCAL INFILE '%s' INTO TABLE NSESecuritiesFullBhavData FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n' IGNORE 1 ROWS (symbol, security_type, @date, prev_close, open_price, high_price, low_price, last_price, close_price, avg_price, ttl_trd_qnty, turnover_lacs, no_of_trades, deliv_qty, deliv_per) set date=STR_TO_DATE(@date, '%s');"
	NSEFBDDateFormat              = "%d-%M-%Y"
	NSEDeliveryPercentageDataLink = "http://www.nseindia.com/archives/equities/mto/MTO_%02d%02d%04d.DAT"
	NSESecuritiesFullBhavDataLink = "http://www.nseindia.com/products/content/sec_bhavdata_full.csv"
	createTableQueryNSESFBD       = "CREATE TABLE IF NOT EXISTS `NSESecuritiesFullBhavData` ( `symbol` varchar (200), `security_type` varchar(10), `date` Date, `prev_close` DOUBLE, `open_price` DOUBLE, `high_price` DOUBLE,`low_price` DOUBLE, `last_price` DOUBLE, `close_price` DOUBLE, `avg_price` DOUBLE, `ttl_trd_qnty` INTEGER, `turnover_lacs` DOUBLE, `no_of_trades` INTEGER, `deliv_qty` INTEGER, `deliv_per` double, PRIMARY KEY (`symbol`, `date`)) ENGINE=InnoDB DEFAULT CHARSET=utf8 ;"
	deleteTableQueryNSEFBD        = "delete from NSESecuritiesFullBhavData where date in (select date from (select min(date) date from NSESecuritiesFullBhavData) D);"

	NSEGetLiveQuoteURL = "http://nseindia.com/live_market/dynaContent/live_watch/get_quote/GetQuote.jsp?symbol=%s&illiquid=0&smeFlag=0&itpFlag=0"
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
	//"INDUSTRY": "http://www.nseindia.com/content/indices/ind_cnxindustrylist.csv",
}

var NSEBroadMarketIndexList = map[string]string{
	"CNX_NIFTY":        "http://www.nseindia.com/content/indices/ind_niftylist.csv",
	"CNX_NIFTY_JUNIOR": "http://www.nseindia.com/content/indices/ind_jrniftylist.csv",
	"CNX_100":          "http://www.nseindia.com/content/indices/ind_cnx100list.csv",
	"CNX_200":          "http://www.nseindia.com/content/indices/ind_cnx200list.csv",
	"CNX_500":          "http://www.nseindia.com/content/indices/ind_cnx500list.csv",
	"NIFTY_MIDCAP_50":  "http://www.nseindia.com/content/indices/ind_niftymidcap50list.csv",
	"CNX_MIDCAP":       "http://www.nseindia.com/content/indices/ind_cnxmidcaplist.csv",
	"CNX_SMALLCAP":     "http://www.nseindia.com/content/indices/ind_cnxsmallcap.csv",
}

type Outer struct {
	Futlink        string   `json:"futLink"`
	OtherSeries    []string `json:"otherSeries"`
	LastUpdateTime string   `json:"lastUpdateTime"`
	TradedDate     string   `json:"tradedDate"`
	Data           []Data   `json:"data"`
	OptLink        string   `json:"optLink"`
}

type Data struct {
	ExtremeLossMargin        string `json:"extremeLossMargin"`
	Cm_ffm                   string `json:"cm_ffm"`
	BcStartDate              string `json:"bcStartDate"`
	Change                   string `json:"change"`
	BuyQuantity1             string `json:"buyQuantity1"`
	BuyQuantity2             string `json:"buyQuantity2"`
	BuyQuantity3             string `json:"buyQuantity3"`
	BuyQuantity4             string `json:"buyQuantity4"`
	BuyQuantity5             string `json:"buyQuantity5"`
	SellPrice1               string `json:"sellPrice1"`
	SellPrice2               string `json:"sellPrice2"`
	SellPrice3               string `json:"sellPrice3"`
	SellPrice4               string `json:"sellPrice4"`
	SellPrice5               string `json:"sellPrice5"`
	PriceBand                string `json:"priceBand"`
	DeliveryQuantity         string `json:"deliveryQuantity"`
	QuantityTraded           string `json:"quantityTraded"`
	Open                     string `json:"open"`
	Low52                    string `json:"Low52"`
	SecurityVar              string `json:"securityVar"`
	MarketType               string `json:"marketType"`
	TotalTradedValue         string `json:"totalTradedValue"`
	Pricebandupper           string `json:"pricebandupper"`
	FaceValue                string `json:"faceValue"`
	NdStartDate              string `json:"ndStartDate"`
	PreviousClose            string `json:"previousClose"`
	Symbol                   string `json:"symbol"`
	VarMargin                string `json:"varMargin"`
	LastPrice                string `json:"lastPrice"`
	PChange                  string `json:"pChange"`
	AdhocMargin              string `json:"adhocMargin"`
	CompanyName              string `json:"companyName"`
	averagePrice             string `json:"averagePrice"`
	SecDate                  string `json:"secDate"`
	Series                   string `json:"series"`
	IsinCode                 string `json:"isinCode"`
	IndexVar                 string `json:"indexVar"`
	Pricebandlower           string `json:"pricebandlower"`
	TotalBuyQuantity         string `json:"totalBuyQuantity"`
	High52                   string `json:"high52"`
	Purpose                  string `json:"purpose"`
	Cm_adj_low_dt            string `json:"cm_adj_low_dt"`
	ClosePrice               string `json:"closePrice"`
	RecordDate               string `json:"recordDate"`
	Cm_adj_high_dt           string `json:"cm_adj_high_dt"`
	TotalSellQuantity        string `json:"totalSellQuantity"`
	DayHigh                  string `json:"dayHigh"`
	ExDate                   string `json:"exDate"`
	SellQuantity1            string `json:"sellQuantity1"`
	SellQuantity2            string `json:"sellQuantity2"`
	SellQuantity3            string `json:"sellQuantity3"`
	SellQuantity4            string `json:"sellQuantity4"`
	SellQuantity5            string `json:"sellQuantity5"`
	BcEndDate                string `json:"bcEndDate"`
	Css_status_desc          string `json:"css_status_desc"`
	NdEndDate                string `json:"ndEndDate"`
	BuyPrice1                string `json:"buyPrice1"`
	BuyPrice2                string `json:"buyPrice2"`
	BuyPrice3                string `json:"buyPrice3"`
	BuyPrice4                string `json:"buyPrice4"`
	BuyPrice5                string `json:"buyPrice5"`
	ApplicableMargin         string `json:"applicableMargin"`
	DayLow                   string `json:"dayLow"`
	DeliveryToTradedQuantity string `json:"deliveryToTradedQuantity"`
	TotalTradedVolume        string `json:"totalTradedVolume"`
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

func getNSEIndexList(list map[string]string) {
	for key, value := range list {
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

		sqlQueryTruncateTable := fmt.Sprintf(truncateTableQuery, key)
		rows, err := proddbhandle.Query(sqlQueryTruncateTable)
		if err != nil {
			glog.Errorln(err)
		}
		if rows != nil {
			rows.Close()
		}

		sqlQueryLoadFIle := fmt.Sprintf(loadFileQuery, path, key, 1)
		rows, err = proddbhandle.Query(sqlQueryLoadFIle)
		if err != nil {
			glog.Errorln(err)
		}
		if rows != nil {
			rows.Close()
		}

		cmd := exec.Command("rm", "-f", path)
		err = cmd.Run()
		if err != nil {

			glog.Error("Error removing the file: ", path)
		}

		//log.Println(resp)
		resp.Body.Close()
		file.Close()
		//req.Close
	}

}

func getNSEBroadMarketIndexLists() {
	glog.Infoln("============Getting NSE Broad Market Indices along with Listed Comapnies==============")
	/* TBD AJAY req/resp/client which objects should be created outside loop?*/
	getNSEIndexList(NSEBroadMarketIndexList)
}

func getNSESectoralIndexLists() {
	glog.Infoln("============Getting NSE Sectoral Indices along with Listed Comapnies==============")
	/* TBD AJAY req/resp/client which objects should be created outside loop?*/
	getNSEIndexList(NSESectoralIndexList)
}

func getNSEDeliveryPercentageData(noOfDays int) {
	count := 0
	for i := 0; count < noOfDays; i++ {
		/* get date for day ith */
		today := time.Now().Add(time.Duration(-86400*i) * time.Second)
		if today.Weekday() == time.Saturday || today.Weekday() == time.Sunday {
			continue
		}
		day := today.Day()
		month := int(today.Month())
		year := today.Year()
		NSEDeliveryPercentageDataUrl := fmt.Sprintf(NSEDeliveryPercentageDataLink, day, month, year)
		fmt.Println("Delivery daya for ", day, month, year, NSEDeliveryPercentageDataUrl)

		/* preapre the http get req */
		req, err := http.NewRequest("GET", NSEDeliveryPercentageDataUrl, nil)
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

		/* get the response */
		resp, err := client.Do(req)
		if err != nil {
			glog.Errorln(":Result:Fail:Error:", err.Error())
			continue
		}
		//defer resp.Body.Close()

		/* file path where we need to store delivery data */
		filePath := strings.Split(NSEDeliveryPercentageDataUrl, forwardSlashChar)
		path := fileDownloadPath + filePath[len(filePath)-1]
		glog.Infoln(resp.Status)

		/* fetch n store file */
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

		/* create table query */
		NSEDeliveryPercentageTableName := fmt.Sprintf(NSEDeliveryPercentageTable, day, month, year)
		sqlQueryCreateTable := fmt.Sprintf(createTableQuery, NSEDeliveryPercentageTableName)
		rows, err := proddbhandle.Query(sqlQueryCreateTable)
		if err != nil {
			glog.Errorln(err)
		}
		if rows != nil {
			rows.Close()
		}

		/* truncate tables if already exist */
		sqlQueryTruncateTable := fmt.Sprintf(truncateTableQuery, NSEDeliveryPercentageTableName)
		rows, err = proddbhandle.Query(sqlQueryTruncateTable)
		if err != nil {
			glog.Errorln(err)
		}
		if rows != nil {
			rows.Close()
		}

		/* Load csv into mysql */
		sqlQueryLoadFIle := fmt.Sprintf(loadFileQuery, path, NSEDeliveryPercentageTableName, 4)
		rows, err = proddbhandle.Query(sqlQueryLoadFIle)
		if err != nil {
			glog.Errorln(err)
		}
		if rows != nil {
			rows.Close()
		}

		/* delete the downloaded file */
		cmd := exec.Command("rm", "-f", path)
		err = cmd.Run()
		if err != nil {

			glog.Error("Error removing the file: ", path)
		}

		/* free the stuff */
		//log.Println(resp)
		resp.Body.Close()
		file.Close()
		//req.Close
		count++
	}
}

func getNSESecuritiesFullBhavData(deleteFromTable bool) { //noOfDays int) {
	/* TBD ajay fetch data for today, this one is for yesterday */
	today := time.Now().Add(time.Duration(-86400*1) * time.Second)
	if today.Weekday() == time.Saturday || today.Weekday() == time.Sunday {
		return
	}

	/* preapre the http get req */
	req, err := http.NewRequest("GET", NSESecuritiesFullBhavDataLink, nil)
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

	/* get the response */
	resp, err := client.Do(req)
	if err != nil {
		glog.Errorln(":Result:Fail:Error:", err.Error())
		return
	}

	/* file path where we need to store delivery data */
	filePath := strings.Split(NSESecuritiesFullBhavDataLink, forwardSlashChar)
	path := fileDownloadPath + filePath[len(filePath)-1]
	glog.Infoln(resp.Status)

	/* fetch n store file */
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

	/* Load csv into mysql */
	sqlQueryLoadFIle := fmt.Sprintf(loadFileQueryNSEFBD, path, NSEFBDDateFormat)
	fmt.Println("bhav data file ", path)
	fmt.Println("bhav query", sqlQueryLoadFIle)

	rows, err := proddbhandle.Query(sqlQueryLoadFIle)
	if err != nil {
		glog.Errorln(err)
		fmt.Println(err)
	}
	if rows != nil {
		rows.Close()
	}

	/* delete from table the data for oldest day */
	if deleteFromTable == true {
		rows, err = proddbhandle.Query(deleteTableQueryNSEFBD)
		if err != nil {
			glog.Errorln(err)
			fmt.Println(err)
		}
		if rows != nil {
			rows.Close()
		}
	}

	/* delete the downloaded file */
	cmd := exec.Command("rm", "-f", path)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		glog.Error("Error removing the file: ", path)
	}

	/* free the stuff */
	resp.Body.Close()
	file.Close()
}

func retriveNSESecuritiesTradeSignals() {
	// BUY SIGNAL ALGORITHM
	/* join full bhav copy daya, close price, % delivery data */
	/* get eps, pe, pe-industry, 52 week high low */
	/* poll current NSE order book */

	// SELL SIGNAL ALGORITHM
}

func getFiftyTwoWeekHighLow() {

	/* get quote for each script, update the 52 week high low
	read quote for each actively traded script n read 52 week high low
	insert or update into TradedCompanyInfo table */

	/* TBD AJAY we may have to convert iso-8859-1 to utf-8 */

	/* preapre the http get req */
	symbol := "ABB"
	reqURL := fmt.Sprintf(NSEGetLiveQuoteURL, symbol)
	fmt.Println("reqURL", reqURL)
	req, err := http.NewRequest("GET", reqURL, nil)
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

	/* get the response */
	resp, err := client.Do(req)
	if err != nil {
		glog.Errorln(":Result:Fail:Error:", err.Error())
		return
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	quoteDataStr := buf.String()

	/* TBD AJAY remove extra comma, from values */
	strs := strings.Split(quoteDataStr, "\n")
	for i := range strs {
		if strings.Contains(strs[i], "futLink") {
			var o Outer
			if err = json.Unmarshal([]byte(strs[i]), &o); err != nil {
				panic(err)

			}
			//fmt.Printf("%+v", o)
			break
		}
	}

	resp.Body.Close()

}

func Serve() {
	/* TBD AJAY
	decide upon the structure of code,
	write seperate files
	convert each function to an api
	*/
	initDB()
	client = &http.Client{}
	/* Call to this function depends on passed argument */
	//getNSESectoralIndexLists()
	//getNSEBroadMarketIndexLists()
	//getNSEDeliveryPercentageData(5)
	//getNSESecuritiesFullBhavData(false)
	getFiftyTwoWeekHighLow()
	retriveNSESecuritiesTradeSignals()
}
