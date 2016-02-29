package db

import (
	"database/sql"
	//"database/sql/driver"
	"fmt"
	"github.com/ajaybodhe/stocks-contra/coreStructures"
	"github.com/ajaybodhe/stocks-contra/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	//"time"
)

const nseCorporateAnnouncementReadQuery = "select company, symbol, date, subject, announcement from NseCorporateAnnouncement where date = (select max(date) from NseCorporateAnnouncement);"

func ReadMaxDateRecordNseCorporateAnnouncement(db util.DB) (*coreStructures.NseCorporateAnnouncementData, error){
	var rows *sql.Rows
	var err error
	rows, err = db.Query(nseCorporateAnnouncementReadQuery)
	if err != nil {
		return nil, fmt.Errorf("ReadMaxDateRecordNseCorporateAnnouncement err: sql error:%s\n", err.Error())
	}
	defer rows.Close()
	var nseCorporateAnnouncementData coreStructures.NseCorporateAnnouncementData
	for rows.Next() {
		err = rows.Scan(&nseCorporateAnnouncementData.Company,
						&nseCorporateAnnouncementData.Symbol,
						&nseCorporateAnnouncementData.Date,
						&nseCorporateAnnouncementData.Subject,
						&nseCorporateAnnouncementData.Announcement)
		if err != nil {
			glog.Error("error: while reading ReadMaxDateRecordNseCorporateAnnouncement :error:%s", err.Error())
			return nil, err
		}
	}
	return &nseCorporateAnnouncementData, nil
}