package util

const (
	ForwardSlashChar   = "/"
	FileDownloadPath   = "/tmp/"
	TruncateTableQuery = "truncate table %s"
	LoadFileQuery      = "LOAD DATA LOCAL INFILE '%s' INTO TABLE %s FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n' IGNORE %d ROWS;"

	LoadFileQueryNSEFBD           = "LOAD DATA LOCAL INFILE '%s' INTO TABLE NSESecuritiesFullBhavData FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n' IGNORE 1 ROWS (symbol, security_type, @date, prev_close, open_price, high_price, low_price, last_price, close_price, avg_price, ttl_trd_qnty, turnover_lacs, no_of_trades, deliv_qty, deliv_per) set date=STR_TO_DATE(@date, '%s');"
	NSEFBDDateFormat              = "%d-%M-%Y"
	NSEDeliveryPercentageDataLink = "http://www.nseindia.com/archives/equities/mto/MTO_%02d%02d%04d.DAT"
	NSESecuritiesFullBhavDataLink = "http://www.nseindia.com/products/content/sec_bhavdata_full.csv"
	CreateTableQueryNSESFBD       = "CREATE TABLE IF NOT EXISTS `NSESecuritiesFullBhavData` ( `symbol` varchar (200), `security_type` varchar(10), `date` Date, `prev_close` DOUBLE, `open_price` DOUBLE, `high_price` DOUBLE,`low_price` DOUBLE, `last_price` DOUBLE, `close_price` DOUBLE, `avg_price` DOUBLE, `ttl_trd_qnty` INTEGER, `turnover_lacs` DOUBLE, `no_of_trades` INTEGER, `deliv_qty` INTEGER, `deliv_per` double, PRIMARY KEY (`symbol`, `date`)) ENGINE=InnoDB DEFAULT CHARSET=utf8 ;"
	DeleteTableQueryNSEFBD        = "delete from NSESecuritiesFullBhavData where date in (select date from (select min(date) date from NSESecuritiesFullBhavData) D);"

	NSEGetLiveQuoteURL     = "http://nseindia.com/live_market/dynaContent/live_watch/get_quote/GetQuote.jsp?symbol=%s&illiquid=0&smeFlag=0&itpFlag=0"
	MoneyControlURLFetcher = "http://www.moneycontrol.com/mccode/common/autosuggesion.php?query=%s%sp&type=1&section=mc_home"

	/* Moneycontrol Constants */
	EPS             = "EPS (TTM)"
	IndustryPE      = "INDUSTRY P/E"
	PE              = "P/E"
	MarketCap       = "MARKET CAP (Rs Cr)"
	BookValue       = "BOOK VALUE (Rs)"
	Dividend        = "DIV (%)"
	PB              = "PRICE/BOOK"
	PC              = "P/C"
	FV              = "FACE VALUE (Rs)"
	DivYield        = "DIV YIELD.(%)"
	FIIHolding      = "FII"
	DIIHolding      = "DII"
	PromoterHolding = "Promoter"
	OtherHolding    = "Others"
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
