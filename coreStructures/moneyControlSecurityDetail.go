package coreStructures

//import ()

type MoneyControlSecurityStructure struct {
	MarketCap       float32
	PE              float32
	BookValue       float32
	Dividend        float32
	IndustryPE      float32
	EPS             float32
	PC              float32
	PB              float32
	DivYield        float32
	FaceValue       float32
	PromoterHolding float32
	FIIHolding      float32
	DIIHolding      float32
	OtherHolding    float32
	High52          float32
	Low52           float32
	Sector          string
}

type MoneyControlAutoSuggestion struct {
	LinkSrc string `json:"link_src"`
	LinkTrack string `json:"link_track"`
	PdtDsNm string `json:"pdt_dis_nm"`
	ScId string `json:"sc_id"`
	StockName string `json:"stock_name"`
	ScSectorId string `json:"sc_sector_id"`
	ScSector string `json:"sc_sector"`
}

type MoneyControlAutoSuggestionArray struct {
	MoneyControlAutoSuggestions []*MoneyControlAutoSuggestion `json:"money_control_auto_suggestion"`
}