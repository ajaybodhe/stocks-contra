package coreStructures

//import ()

type NseLiveQuoteDetails struct {
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
	AveragePrice             string `json:"averagePrice"`
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
