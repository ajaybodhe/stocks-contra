package api

import (
	"strconv"
	"errors"
	//"io/ioutil"
	"strings"
	"github.com/ajaybodhe/stocks-contra/coreStructures"
	"github.com/PuerkitoBio/goquery"
	"github.com/ajaybodhe/stocks-contra/util"
	"github.com/golang/glog"
	//"errors"
//	"crypto/tls"
	"net/http"
	//"encoding/json"
	//"io/ioutil"
	"io"
	//"strings"
	"os"
	//"os/exec"
	"fmt"
	"compress/gzip"
	"bytes"
	"time"
	"github.com/ajaybodhe/stocks-contra/workers"
)

const (
	nseCompany = "company"
	nseSymbol = "symbol"
	nseDescription = "desc"
	nseDate = "date"
	nseLink = "link" 
)

var LastAnnouncementDate string

func interestedSubjects() func(subject string) bool {
	interestedSubjectsMap := map[string]bool {
		"Updates" : true,
		"Press Release" : true,
		"Record Date" : false,
		"Financial Result Updates": true,
	}
	return func(subject string)bool {
		if val, ok := interestedSubjectsMap[subject]; ok {
			return val
		}
		return false
	}
}

func getISTTIme(date string ) time.Time{
	loc, _ := time.LoadLocation("Asia/Kolkata")
    test, err := time.ParseInLocation("02-Jan-2006 15:04", date, loc)
    if err != nil {
        panic(err)
    }
    //fmt.Println(test)
	return test
}

func downloadNseAnnouncementFile(client *http.Client, url string) (string, error){
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		glog.Fatalln(err)
		return "", err
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
		return "", err
	}
	defer resp.Body.Close()
	
	/* file path where we need to store delivery data */
	filePath := strings.Split(url, util.ForwardSlashChar)
	path := util.FileDownloadPath + filePath[len(filePath)-1]
	
	/* fetch n store file */
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return path, nil
}

func GetNseFullCorporateAnnoucement(client *http.Client, url string) (*coreStructures.NseCorporateAnnouncementData, error){
	req, err := http.NewRequest("GET", url, nil)
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
		return nil, err
	}
	defer resp.Body.Close()
	
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
	
	index := strings.Index(string(buf.Bytes()), "<table")
	if index < 0 {
		return nil, errors.New("no data found")
	}

	htmlData := buf.Bytes()[index:]

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlData))
	if nil != err {
		return nil, err
	}
	
	var date string
	var subject string
	var annoucement string
	var attachmentLink string
	
	doc.Find("tr").Each(func(rowIdx int, s *goquery.Selection) {
		s.Find("td").Each(func (columnIdx int, s *goquery.Selection) {
			if columnIdx == 1 {
				if rowIdx == 0 {
					date = s.Text()
				} else if rowIdx == 1 {
					subject = s.Text() 
				} else if rowIdx == 2 {
					annoucement = s.Text()
				} else if rowIdx == 3 {
					s.Find("a").Each(func(attrIndx int, s *goquery.Selection) {
						c, _ := s.Attr("href")
						attachmentLink = c	
					})	
				}
			}
		})
	})
	var path string
	if attachmentLink != "" {
		path, err = downloadNseAnnouncementFile(client, util.NseBaseURL + attachmentLink)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	
	dbDate := getISTTIme(date)
	return &coreStructures.NseCorporateAnnouncementData{Company : "",
										Symbol : "",
										Subject : subject,
										AttachementLink : util.NseBaseURL + attachmentLink,
										AttachmentFilePath : path,
										Date : dbDate,
										Announcement: annoucement}, 
										nil
}

func getNseCorporateAnnouncementDataValue(announcementDataStr string, substring string)(string, int) {
	var j int
	var c  rune
	i := strings.Index(announcementDataStr, substring)
	if i == -1 {
		return "", i
	}
	quote := false
	dataValue := ""
	for j,c = range announcementDataStr[i:] {
		if quote == true {
			dataValue = dataValue + string(c)
		}
		if c == '"' && quote == true {
			break
		} else if c == '"' {
			quote = true
		}
	}
	dataValue = dataValue[:len(dataValue) - 1]
	return dataValue, i+j
}

func GetNseCorporateAnnouncements(client *http.Client, jq chan workers.Job) error{
	
	url := util.NSECorporateAnnounceMentLink

	start := 0
	limit := 20
	
	check := interestedSubjects()
	
	var mostRecentDate string
	
	for  {
		url1 := url + "?start=" + strconv.Itoa(start) + "&limit=" + strconv.Itoa(limit)
		
		req, err := http.NewRequest("GET", url1, nil)
		if err != nil {
			fmt.Println(err)
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:44.0) Gecko/20100101 Firefox/4")
		req.Header.Set("Host", "www.nseindia.com")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(":Result:Fail:Error:", err.Error())
		}

		defer resp.Body.Close()
	
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
		announcementDataStr := buf.String()
		
		var count int
		for count = range announcementDataStr {
			if announcementDataStr[count] == '{' {
				break
			}
		}
		announcementDataStr = announcementDataStr[count:]
		
		lena :=  len(announcementDataStr)
		//var announcements coreStructures.NseCorporateAnnouncement
		charsCount := 0
		
		for  charsCount < lena {
			company, i := getNseCorporateAnnouncementDataValue(announcementDataStr[charsCount:], nseCompany)
			if i == -1 {
				break
			}
			charsCount += i
		
			symbol, i := getNseCorporateAnnouncementDataValue(announcementDataStr[charsCount:], nseSymbol)
			if i == -1 {
				break
			}
			charsCount += i
		
			desc, i := getNseCorporateAnnouncementDataValue(announcementDataStr[charsCount:], nseDescription)
			if i == -1 {
				break
			}
			charsCount += i
		
			link, i := getNseCorporateAnnouncementDataValue(announcementDataStr[charsCount:], nseLink)
			if i == -1 {
				break
			}
			charsCount += i
		
			sdate, i := getNseCorporateAnnouncementDataValue(announcementDataStr[charsCount:], nseDate)
			if i == -1 {
				break
			}
			charsCount += i
			
			if sdate >mostRecentDate {
				mostRecentDate = sdate
			}
			if sdate < LastAnnouncementDate {
				break
			}
			//TBD AJAY fetch the actual data attchements & mail
			fullDataURL := util.NseBaseURL + link
			if check(desc) == false {
				continue
			}
			
			jq <- workers.Job{Payload:
				coreStructures.NseShortCorporateAnnouncement{
					Symbol:symbol,
					Company:company,
					FullDataURL:fullDataURL,
				}}
			//TBD AJAY run this in goroutine after 15 seconds
			//TBD AJAY stop after the last read record is found
		}
		LastAnnouncementDate = mostRecentDate
		time.Sleep(1 * time.Minute)
	}
 }