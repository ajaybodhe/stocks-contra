package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	TOTAL_PAGES = 6
)

type CompanyDetail struct {
	CompanyCode string
	CompanyName string
}

type EdelweissRatio struct {
	CompanyCode      string
	CompanyName      string
	EdelweissMetrics []*EdelweissMetric
}

type EdelweissMetric struct {
	Title         string
	MetricDetails []*MetricDetail
}

type MetricDetail struct {
	Id         int
	ShortName  string
	FullName   string
	Percentage string
}

func parseMetrics(companyCode string, pageNumber int,
	metricTitleMap map[string]int, edelweissMetrics *[]*EdelweissMetric) []*EdelweissMetric {
	http_url := "https://www.edelweiss.in/Market/CompAPI.aspx"

	resp, err := http.PostForm(http_url,
		url.Values{"CoCode": {companyCode},
			"Type":        {"RATIOS"},
			"currpage":    {strconv.Itoa(pageNumber)},
			"isfirstload": {"false"},
		})

	if nil != err {
		log.Panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	index := strings.Index(string(body), "<table")

	if index < 0 {
		return nil
	}

	htmlData := body[index:]

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlData))

	if nil != err {
		return nil
	}

	tempArray := make([]int, 0)

	doc.Find("tr").Each(func(rowIdx int, s *goquery.Selection) {
		if rowIdx == 0 {
			s.Find("td").Each(func(columnIdx int, s *goquery.Selection) {
				metricTitle := s.Text()

				mIdx := metricTitleMap[metricTitle]
				if 0 == mIdx {
					metricTitleMapSize := len(metricTitleMap)
					tempArray = append(tempArray, metricTitleMapSize)
					metricTitleMap[metricTitle] = metricTitleMapSize + 1

					edelweissMetric := &EdelweissMetric{}
					edelweissMetric.Title = metricTitle
					*edelweissMetrics = append(*edelweissMetrics, edelweissMetric)
				} else {
					tempArray = append(tempArray, -1)
				}
			})
		}

		if rowIdx == 2 {
			idx := 0
			s.Find("table").Each(func(tableIdx int, tableSelection *goquery.Selection) {
				metrixIndex := tempArray[idx]

				if -1 != metrixIndex {

					edelweissMetric := (*edelweissMetrics)[metrixIndex]

					metricDetails := make([]*MetricDetail, 0)

					tableSelection.Find("tr").Each(func(mRow int, metric *goquery.Selection) {
						metricDetail := &MetricDetail{}
						metricDetail.Id = mRow
						metric.Find("td").Each(func(cellId int, cell *goquery.Selection) {
							if 0 == cellId {
								title, _ := cell.Attr("title")
								//TODO: need to check ret val
								metricDetail.ShortName = cell.Text()
								metricDetail.FullName = title
							}

							if 2 == cellId {
								metricDetail.Percentage = cell.Text()
							}
						})

						metricDetails = append(metricDetails, metricDetail)
					})

					edelweissMetric.MetricDetails = metricDetails
				}
				idx++
			})
		}
	})

	return *edelweissMetrics
}

func fetchEdelweissMetricsByCompany(companyDetails []*CompanyDetail) {

	for _, companyDetail := range companyDetails {

		edelweissRatio := &EdelweissRatio{}
		edelweissMetrics := make([]*EdelweissMetric, 0)

		metricTitleMap := map[string]int{}

		for i := 1; i <= TOTAL_PAGES; i++ {
			parseMetrics(companyDetail.CompanyCode, i, metricTitleMap, &edelweissMetrics)
		}

		edelweissRatio.CompanyCode = companyDetail.CompanyCode
		edelweissRatio.CompanyName = companyDetail.CompanyName

		edelweissRatio.EdelweissMetrics = edelweissMetrics

		edelweissRatiosJSON, err := json.Marshal(edelweissRatio)

		if nil != err {
			fmt.Println("Error occured while marshaling EdelweissRatio List")
			return
		}

		fmt.Println(string(edelweissRatiosJSON))
	}
}

func getCompanyCode(http_url string) string {

	resp, err := http.Get(http_url)

	if nil != err {
		log.Panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	pageSource := string(body)

	rexpr := regexp.MustCompile("var CoCode=[a-z0-9]+;")
	matchedStr := rexpr.FindString(pageSource)

	tokens := strings.Split(matchedStr, "=")

	code := tokens[1]

	return code[:len(code)-1]
}

func listCompanyWithCompanyCodes() []*CompanyDetail {
	companyDetails := make([]*CompanyDetail, 0)

	pageList := []string{"a", "b", "c", "d", "e", "f", "g",
		"h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r",
		"s", "t", "u", "v", "w", "x", "y", "z", "0"}

	for _, page := range pageList {

		http_url := "https://www.edelweiss.in/market/liststocks/" + page + ".html"

		resp, err := http.PostForm(http_url, nil)

		if nil != err {
			log.Panic(err)
		}

		doc, err := goquery.NewDocumentFromResponse(resp)

		if nil != err {
			log.Panic(err)
		}

		doc.Find("#divIndexstklist").Find("table").Each(func(tableIdx int,
			tableSelection *goquery.Selection) {
			tableSelection.Find("tr").Each(func(trIdx int, trSelection *goquery.Selection) {
				trSelection.Find("td").Each(func(tdIdx int, tdSelection *goquery.Selection) {
					tdSelection.Find("a").Each(func(aIdx int, aSelection *goquery.Selection) {
						url, _ := aSelection.Attr("href")
						cname, _ := aSelection.Find("span").Attr("title")
						cocode := getCompanyCode(url)
						//fmt.Printf("%s|%s\n", cocode, cname)
						companyDetail := &CompanyDetail{CompanyCode: cocode, CompanyName: cname}
						companyDetails = append(companyDetails, companyDetail)
					})
				})
			})
		})
	}

	return companyDetails
}

func crawl() {
	//1. Enable this to fetch all company details
	// This operation take very long time, so its better we store this
	// company details in one table / file and fetch it to initialize companyDetails
	// array
	//companyDetails := listCompanyWithCompanyCodes()

	// Initialized for testing perpose only
	companyDetails := make([]*CompanyDetail, 0)
	companyDetail := &CompanyDetail{CompanyCode: "242", CompanyName: "A B B"}
	companyDetails = append(companyDetails, companyDetail)
	companyDetail = &CompanyDetail{CompanyCode: "2323", CompanyName: "HMT Ltd"}
	companyDetails = append(companyDetails, companyDetail)

	/*
		companyDetailsJSON, err := json.Marshal(companyDetails)

		if nil != err {
			fmt.Println("Error occured while marshaling EdelweissRatio List")
			return
		}

		fmt.Println(string(companyDetailsJSON))
	*/

	fetchEdelweissMetricsByCompany(companyDetails)
}
