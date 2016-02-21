package coreStructures

//import(
//	"encoding/json"
//)

type NseCorporateAnnouncementData struct {
	Company string `json:"company"`
	Symbol string `json:"symbol"`
	Subject string `json:"desc"`
	AttachementLink string `json:"link"`
	AttachmentFilePath string `json:"path"`
	Date string `json:"date"`
	Announcement string `json:"announcement"`
}

type NseCorporateAnnouncement struct{
	Data []NseCorporateAnnouncementData `json:"rows"`
	Success string `json:"success"`
	ResultCount int `json:"results"`
}