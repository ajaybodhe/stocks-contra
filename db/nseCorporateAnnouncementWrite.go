package db

import (
	//"database/sql"
	//"database/sql/driver"
	"fmt"
	"github.com/ajaybodhe/stocks-contra/coreStructures"
	"github.com/ajaybodhe/stocks-contra/util"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/golang/glog"
	//"time"
)

//const nseCorporateAnnouncementWriteQuery = "insert into NseCorporateAnnouncement(symbol, company, date, subject, announcement) values"

func deleteOldNseAnnouncements(db util.DB) error{
	deleteSql := "delete from NseCorporateAnnouncement where date < SUBDATE(NOW(), interval 3 day);"
	_, err := db.Exec(deleteSql)
	if err != nil {
		return fmt.Errorf("deleteOldNseAnnouncements: sql error:%s\n", err.Error())
	}
	return nil
}

func WriteNseCorporateAnnouncements(db util.DB, announcements []*coreStructures.NseCorporateAnnouncementData) error{
	nseCorporateAnnouncementWriteQuery := "insert into NseCorporateAnnouncement(symbol, company, date, subject, announcement) values"
	firstTime := true
	for _, announcement := range announcements {
		if firstTime {
			firstTime = false
			nseCorporateAnnouncementWriteQuery += fmt.Sprintf("(\"%v\", \"%v\", \"%v\", \"%v\", \"%v\")", announcement.Symbol, announcement.Company, announcement.Date, announcement.Subject, announcement.Announcement)
		} else {
			nseCorporateAnnouncementWriteQuery += fmt.Sprintf(",(\"%v\", \"%v\",\"%v\", \"%v\", \"%v\")", announcement.Symbol, announcement.Company, announcement.Date, announcement.Subject, announcement.Announcement)
		}
	}
	fmt.Printf("the query is: %s", nseCorporateAnnouncementWriteQuery)
	_, err := db.Exec(nseCorporateAnnouncementWriteQuery)
	if err != nil {
		return fmt.Errorf("nseCorporateAnnouncementWriteQuery: sql error:%s\n", err.Error())
	}
	
	err = deleteOldNseAnnouncements(db)
	if err != nil {
		return fmt.Errorf("deleteOldNseAnnouncements: sql error:%s\n", err.Error())
	}
	
	return nil
}