package util

const (
	NewLineChar           = "\n"
	SpaceChar             = " "
	PercentageChar        = "%"
	CommaChar             = ","
	ForwardSlashChar      = "/"
	EmptyString           = ""
	StartBreakTag         = "<b>"
	EndBreakTag           = "</b>"
	MoneycontrolNSESymbol = "NSE: "
	FileDownloadPath      = "/tmp/"
	TruncateTableQuery    = "truncate table %s"
	LoadFileQuery         = "LOAD DATA LOCAL INFILE '%s' INTO TABLE %s FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n' IGNORE %d ROWS;"

	LoadFileQueryNSEFBD           = "LOAD DATA LOCAL INFILE '%s' INTO TABLE NSESecuritiesFullBhavData FIELDS TERMINATED BY ',' ENCLOSED BY '\"' LINES TERMINATED BY '\n' IGNORE 1 ROWS (symbol, security_type, @date, prev_close, open_price, high_price, low_price, last_price, close_price, avg_price, ttl_trd_qnty, turnover_lacs, no_of_trades, deliv_qty, deliv_per) set date=STR_TO_DATE(@date, '%s');"
	NSEFBDDateFormat              = "%d-%M-%Y"
	NseBaseURL                    = "https://www.nseindia.com/"
	NSEDeliveryPercentageDataLink = "https://www.nseindia.com/archives/equities/mto/MTO_%02d%02d%04d.DAT"
	NSESecuritiesFullBhavDataLink = "https://www.nseindia.com/products/content/sec_bhavdata_full.csv"
	NSECorporateAnnounceMentLink  = "https://www.nseindia.com/corporates/directLink/latestAnnouncementsCorpHome.jsp"
	BSECorporateAnnounceMentLink  = "http://www.bseindia.com/corporates/ann.aspx"
	CreateTableQueryNSESFBD       = "CREATE TABLE IF NOT EXISTS `NSESecuritiesFullBhavData` ( `symbol` varchar (200), `security_type` varchar(10), `date` Date, `prev_close` DOUBLE, `open_price` DOUBLE, `high_price` DOUBLE,`low_price` DOUBLE, `last_price` DOUBLE, `close_price` DOUBLE, `avg_price` DOUBLE, `ttl_trd_qnty` INTEGER, `turnover_lacs` DOUBLE, `no_of_trades` INTEGER, `deliv_qty` INTEGER, `deliv_per` double, PRIMARY KEY (`symbol`, `date`)) ENGINE=InnoDB DEFAULT CHARSET=utf8 ;"
	DeleteTableQueryNSEFBD        = "delete  from NSESecuritiesFullBhavData where date in (select minDate from (select min(date) as minDate from NSESecuritiesFullBhavData) as X) and exists (select count from (select if(count(distinct date)> 10, count(distinct date), 0) as count from NSESecuritiesFullBhavData) as Y where count > 0 );"
	//DeleteTableQueryNSEFBD        = "delete from NSESecuritiesFullBhavData where date in (select date from (select min(date) date from NSESecuritiesFullBhavData) D);"

	NSEGetLiveQuoteURL = "https://nseindia.com/live_market/dynaContent/live_watch/get_quote/GetQuote.jsp?symbol=%s&illiquid=0&smeFlag=0&itpFlag=0"
	//MoneyControlURLFetcher = "http://www.moneycontrol.com/mccode/common/autosuggesion.php?query=%s%sp&type=1&section=mc_home"
	MoneyControlURLFetcher = "http://www.moneycontrol.com/mccode/common/autosuggesion.php?query=%s&type=1&callback=suggest1&section=mc_home"

	/* Moneycontrol Constants */
	Sector              = "SECTOR:"
	HighLow52Week       = "52 Wk Low/High"
	EPS                 = "EPS (TTM)"
	IndustryPE          = "INDUSTRY P/E"
	PE                  = "P/E"
	MarketCap           = "MARKET CAP (Rs Cr)"
	BookValue           = "BOOK VALUE (Rs)"
	Dividend            = "DIV (%)"
	PB                  = "PRICE/BOOK"
	PC                  = "P/C"
	FV                  = "FACE VALUE (Rs)"
	DivYield            = "DIV YIELD.(%)"
	ShareHoldingPattern = "Share Holding Pattern in (%)"
	FIIHolding          = "FII"
	DIIHolding          = "DII"
	PromoterHolding     = "Promoter"
	OtherHolding        = "Others"
	/* we can use regex matching to skyp 4 chars, match with [0-9.]*/
	MoneControlLiveQuoteSkipCharCount       = 4
	MoneControlSectorSkipCharCount          = 5
	MoneControlAlternateSectorSkipCharCount = 8
	MoneControlPromoterHoldingSkipCharCount = 1
	MoneControlFIIHoldingSkipCharCount      = 6
	CheckScriptChars                        = 3
	MoneControlOtherHoldingSkipCharCount    = 3
	FloatSizeBit32                          = 32

	InvalidStrategy         = 0
	StocksContra            = 1
	FiveConsecutiveDaysDown = 2
	MajorCorrection         = 3
	
	BloomBits = 71888
	BloomHashCount = 10
	BloomStockEntries = 5000
	BloomFPError = 0.001
	
	NSEBookAnalystsCount = 10
	
	NIFTY_50 = "NIFTY_50"
	NIFTY_NEXT_50 = "NIFTY_NEXT_50"
)

var NSESectoralIndexList = map[string]string{
	"AUTO":     "https://www.nseindia.com/content/indices/ind_niftyautolist.csv",
	"BANK":     "https://www.nseindia.com/content/indices/ind_niftybanklist.csv",
	"PRIVATE_BANKS":"https://www.nseindia.com/content/indices/ind_nifty_privatebanklist.csv",
	"FINANCE":  "https://www.nseindia.com/content/indices/ind_niftyfinancelist.csv",
	"FMCG":     "https://www.nseindia.com/content/indices/ind_niftyfmcglist.csv",
	"IT":       "https://www.nseindia.com/content/indices/ind_niftyitlist.csv",
	"MEDIA":    "https://www.nseindia.com/content/indices/ind_niftymedialist.csv",
	"METAL":    "https://www.nseindia.com/content/indices/ind_niftymetallist.csv",
	"PHARMA":   "https://www.nseindia.com/content/indices/ind_niftypharmalist.csv",
	"PSU_BANK": "https://www.nseindia.com/content/indices/ind_niftypsubanklist.csv",
	"REALTY":   "https://www.nseindia.com/content/indices/ind_niftyrealtylist.csv",
//	"INDUSTRY": "https://www.nseindia.com/content/indices/ind_niftyindustrylist.csv",
}

var NSEBroadMarketIndexList = map[string]string{
	"NIFTY_50":        "https://www.nseindia.com/content/indices/ind_nifty50list.csv",
	"NIFTY_NEXT_50":   "https://www.nseindia.com/content/indices/ind_niftynext50list.csv",
	"NIFTY_100":          "https://www.nseindia.com/content/indices/ind_nifty100list.csv",
	"NIFTY_200":          "https://www.nseindia.com/content/indices/ind_nifty200list.csv",
	"NIFTY_500":          "https://www.nseindia.com/content/indices/ind_nifty500list.csv",
	"NIFTY_MIDCAP_50":  "https://www.nseindia.com/content/indices/ind_niftymidcap50list.csv",
	"NIFTY_MIDCAP_FULL_100" : "https://www.nseindia.com/content/indices/ind_niftyfullmidcap100list.csv",
	"NIFTY_MIDCAP_100":  "https://www.nseindia.com/content/indices/ind_niftyfullmidcap100list.csv",
	"NIFTY_FREE_FLOAT_MIDCAP_100" : "https://www.nseindia.com/content/indices/ind_niftyfreefloatMidcap100list.csv",
	"NIFTY_SMALLCAP_250" : "https://www.nseindia.com/content/indices/ind_niftysmallcap250list.csv",
	"NIFTY_SMALLCAP_50" : "https://www.nseindia.com/content/indices/ind_niftysmallcap50list.csv",
	"NIFTY_SMALLCAP_FULL_100" : "https://www.nseindia.com/content/indices/ind_niftyfullsmallcap100list.csv",
	"NIFTY_FREE_FLOAT_SMALLCAP_100" : "https://www.nseindia.com/content/indices/ind_niftyfreefloatMidcap100list.csv",
	"NIFTY_MIDSMALLCAP_400" : "https://www.nseindia.com/content/indices/ind_niftymidsmallcap400list.csv",
	}
