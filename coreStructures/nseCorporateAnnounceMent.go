package coreStructures

import (
	"time"
)

//import(
//	"encoding/json"
//)

type NseCorporateAnnouncementData struct {
	Company string `json:"company"`
	Symbol string `json:"symbol"`
	Subject string `json:"desc"`
	AttachementLink string `json:"link"`
	AttachmentFilePath string `json:"path"`
	Date time.Time `json:"date"`
	Announcement string `json:"announcement"`
}

type NseCorporateAnnouncement struct{
	Data []NseCorporateAnnouncementData `json:"rows"`
	Success string `json:"success"`
	ResultCount int `json:"results"`
}

type NseShortCorporateAnnouncement struct {
	Company string `json:"company"`
	Symbol string `json:"symbol"`
	//LatAnnouncementDate time.Time `json:"date"`
	FullDataURL string `json:"full_data_url"`
}