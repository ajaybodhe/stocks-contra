package api

import (
	"strconv"
	//"errors"
	//"io/ioutil"
	//"strings"
	//"github.com/ajaybodhe/stocks-contra/coreStructures"
	"github.com/PuerkitoBio/goquery"
	//"github.com/ajaybodhe/stocks-contra/db"
	//"github.com/ajaybodhe/stocks-contra/conf"
	"github.com/ajaybodhe/stocks-contra/util"
	"github.com/golang/glog"
	"errors"
	"net/http"
	//"encoding/json"
	//"io/ioutil"
	"io"
	"strings"
	//"os"
	//"os/exec"
	"fmt"
	"compress/gzip"
	"bytes"
	//"os/exec"
	//"time"
)

func interestedBseSubjects() func(subject string) bool {
	interestedBseSubjectsMap := map[string]bool {
		"Company Update" : true,
		"Press Release" : true,
		"Record Date" : false,
		"Financial Result Updates": true,
	}
	return func(subject string)bool {
		if val, ok := interestedBseSubjectsMap[subject]; ok {
			return val
		}
		return false
	}
}

func GetBseCorporateAnnouncements(client *http.Client, proddbhandle util.DB, url string) error{
	start := 1
	limit := 20
	stopFlag := false
//	announcements := make([]*coreStructures.BseCorporateAnnouncementData, 10)
	check := interestedBseSubjects()
	
	for start < 21 && stopFlag == false {
		url1 := url + "?curpg=" + strconv.Itoa(start) + "&annflag=1&dt=&dur=D&dtto=&cat=&scrip=&anntype=C"
		fmt.Printf("\nurl=%v\n", url1)
		req, err := http.NewRequest("GET", url1, nil)
		if err != nil {
			glog.Fatalln(err)
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:28.0) Gecko/20100101 Firefox/28.0")
		req.Header.Set("Host", "www.bseindia.com")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Accept-Encoding", "gzip, deflate")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		resp, err := client.Do(req)
		if err != nil {
			glog.Errorln(":Result:Fail:Error:", err.Error())
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
		
		index := strings.Index(string(buf.Bytes()), "<span id=\"ctl00_ContentPlaceHolder1_lblann\">")
		if index < 0 {
			return errors.New("no data found")
		}

		htmlData := buf.Bytes()[index:]
		
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlData))
		if nil != err {
			return err
		}
	
		var date string
		var subject string
		var announcement string
		var attachmentLink string
		var description string
	
		doc.Find("table").Each(func(rowIdx int, s *goquery.Selection) {
			annoucementCount := 0
			s.Find("tr").Each(func (columnIdx int, s *goquery.Selection) {
				s.Find("td").Each(func (columnIdx1 int, s *goquery.Selection) {
					if columnIdx1 == 0 && columnIdx != 0{
						if annoucementCount == 0 {
							description = s.Text()
							fmt.Printf("description=%s\n", s.Text())
							annoucementCount = 1
						} else if annoucementCount == 1{
							announcement = s.Text()
							fmt.Printf("announcement=%s\n", s.Text())
							annoucementCount = 2
						} else if annoucementCount == 2{
							fmt.Printf("nothing=%s\n\n\n\n", s.Text())
							annoucementCount= 0
							// TBD Ajay we should complete corporate announcement here
							
						}
					} else if columnIdx1 == 1 {
						subject = s.Text()
						if check(subject) == false {
							return
						}
						fmt.Printf("subject=%s\n", s.Text())
					} else if columnIdx1 == 2 {
						s.Find("a").Each(func(attrIndx int, s *goquery.Selection) {
							c, _ := s.Attr("href")
							attachmentLink = c
							fmt.Printf("link=%s\n", c)	
						})	
					} else if columnIdx1 == 3 {
						date = s.Text()
						fmt.Printf("date=%s\n", s.Text())
					}
				})
			})
		})
		//fmt.Printf("\nResponse Status = %v\n\nResponse Body = %+v\n", resp.Status, string(htmlData))
		fmt.Printf("descriptionOutside=%s\n", description)
		fmt.Printf("announcementOutside=%s\n", announcement)
		
		start += limit
	}
	return nil
 }