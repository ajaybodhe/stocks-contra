package api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/ajaybodhe/stocks-contra/coreStructures"
	"github.com/ajaybodhe/stocks-contra/util"
	"github.com/golang/glog"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func getNSEIndexList(client *http.Client, proddbhandle util.DB, list map[string]string) {
	for key, value := range list {
		glog.Infoln(key, value)

		req, err := http.NewRequest("GET", value, nil)
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
			continue
		}
		//defer resp.Body.Close()

		filePath := strings.Split(value, util.ForwardSlashChar)
		path := util.FileDownloadPath + filePath[len(filePath)-1]
		glog.Infoln(resp.Status)

		file, err := os.Create(path)
		if err != nil {
			glog.Errorln(err)
		}
		//defer file.Close()

		size, err := io.Copy(file, resp.Body)
		if err != nil {
			glog.Errorln(err)
		}
		glog.Infoln("%s with %v bytes downloaded", path, size)

		sqlQueryTruncateTable := fmt.Sprintf(util.TruncateTableQuery, key)
		rows, err := proddbhandle.Query(sqlQueryTruncateTable)
		if err != nil {
			glog.Errorln(err)
		}
		if rows != nil {
			rows.Close()
		}

		sqlQueryLoadFIle := fmt.Sprintf(util.LoadFileQuery, path, key, 1)
		rows, err = proddbhandle.Query(sqlQueryLoadFIle)
		if err != nil {
			glog.Errorln(err)
		}
		if rows != nil {
			rows.Close()
		}

		cmd := exec.Command("rm", "-f", path)
		err = cmd.Run()
		if err != nil {

			glog.Error("Error removing the file: ", path)
		}

		//log.Println(resp)
		resp.Body.Close()
		file.Close()
		//req.Close
	}

}

func GetNSEBroadMarketIndexLists(client *http.Client, proddbhandle util.DB) {
	glog.Infoln("============Getting NSE Broad Market Indices along with Listed Comapnies==============")
	/* TBD AJAY req/resp/client which objects should be created outside loop?*/
	getNSEIndexList(client, proddbhandle, util.NSEBroadMarketIndexList)
}

func GetNSESectoralIndexLists(client *http.Client, proddbhandle util.DB) {
	glog.Infoln("============Getting NSE Sectoral Indices along with Listed Comapnies==============")
	/* TBD AJAY req/resp/client which objects should be created outside loop?*/
	getNSEIndexList(client, proddbhandle, util.NSESectoralIndexList)
}

func GetNSESecuritiesFullBhavData(client *http.Client, proddbhandle util.DB, deleteFromTable bool) { //noOfDays int) {
	/* TBD ajay fetch data for today, this one is for yesterday */
	today := time.Now().Add(time.Duration(-86400*1) * time.Second)
	if today.Weekday() == time.Saturday || today.Weekday() == time.Sunday {
		return
	}

	/* preapre the http get req */
	req, err := http.NewRequest("GET", util.NSESecuritiesFullBhavDataLink, nil)
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

	/* get the response */
	resp, err := client.Do(req)
	if err != nil {
		glog.Errorln(":Result:Fail:Error:", err.Error())
		return
	}

	/* file path where we need to store delivery data */
	filePath := strings.Split(util.NSESecuritiesFullBhavDataLink, util.ForwardSlashChar)
	path := util.FileDownloadPath + filePath[len(filePath)-1]
	glog.Infoln(resp.Status)

	/* fetch n store file */
	file, err := os.Create(path)
	if err != nil {
		glog.Errorln(err)
	}
	//defer file.Close()
	size, err := io.Copy(file, resp.Body)
	if err != nil {
		glog.Errorln(err)
	}
	glog.Infoln("%s with %v bytes downloaded", path, size)

	/* Load csv into mysql */
	sqlQueryLoadFIle := fmt.Sprintf(util.LoadFileQueryNSEFBD, path, util.NSEFBDDateFormat)

	rows, err := proddbhandle.Query(sqlQueryLoadFIle)
	if err != nil {
		glog.Errorln(err)
		fmt.Println(err)
	}
	if rows != nil {
		rows.Close()
	}

	/* delete from table the data for oldest day */
	if deleteFromTable == true {
		rows, err = proddbhandle.Query(util.DeleteTableQueryNSEFBD)
		if err != nil {
			glog.Errorln(err)
			fmt.Println(err)
		}
		if rows != nil {
			rows.Close()
		}
	}

	/* delete the downloaded file */
	cmd := exec.Command("rm", "-f", path)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		glog.Error("Error removing the file: ", path)
	}

	/* free the stuff */
	resp.Body.Close()
	file.Close()
}

/* TBD Ajay we have the live quote here */
func GetNSELiveQuote(client *http.Client) {

	/* get quote for each script, update the 52 week high low
	read quote for each actively traded script n read 52 week high low
	insert or update into TradedCompanyInfo table */

	/* TBD AJAY we may have to convert iso-8859-1 to utf-8 */

	/* preapre the http get req */
	symbol := "ABB"
	reqURL := fmt.Sprintf(util.NSEGetLiveQuoteURL, symbol)
	fmt.Println("reqURL", reqURL)
	req, err := http.NewRequest("GET", reqURL, nil)
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

	/* get the response */
	resp, err := client.Do(req)
	if err != nil {
		glog.Errorln(":Result:Fail:Error:", err.Error())
		return
	}

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
	quoteDataStr := buf.String()

	/* TBD AJAY remove extra comma, from values */
	strs := strings.Split(quoteDataStr, "\n")
	for i := range strs {
		if strings.Contains(strs[i], "futLink") {
			var nseLQD coreStructures.NseLiveQuoteDetails
			if err = json.Unmarshal([]byte(strs[i]), &nseLQD); err != nil {
				panic(err)

			}
			//fmt.Printf("%+v", nseLQD)
			break
		}
	}

	resp.Body.Close()
}

/*
func GetNSEDeliveryPercentageData(noOfDays int) {
	count := 0
	for i := 0; count < noOfDays; i++ {
		// get date for day ith
		today := time.Now().Add(time.Duration(-86400*i) * time.Second)
		if today.Weekday() == time.Saturday || today.Weekday() == time.Sunday {
			continue
		}
		day := today.Day()
		month := int(today.Month())
		year := today.Year()
		NSEDeliveryPercentageDataUrl := fmt.Sprintf(NSEDeliveryPercentageDataLink, day, month, year)
		fmt.Println("Delivery daya for ", day, month, year, NSEDeliveryPercentageDataUrl)

		// preapre the http get req
		req, err := http.NewRequest("GET", NSEDeliveryPercentageDataUrl, nil)
		if err != nil {
			glog.Fatalln(err)
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:28.0) Gecko/20100101 Firefox/28.0")
		req.Header.Set("Host", "www.nseindia.com")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Accept-Encoding", "gzip, deflate")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,**;q=0.8")

		// get the response
		resp, err := client.Do(req)
		if err != nil {
			glog.Errorln(":Result:Fail:Error:", err.Error())
			continue
		}
		//defer resp.Body.Close()

		// file path where we need to store delivery data
		filePath := strings.Split(NSEDeliveryPercentageDataUrl, forwardSlashChar)
		path := fileDownloadPath + filePath[len(filePath)-1]
		glog.Infoln(resp.Status)

		// fetch n store file
		file, err := os.Create(path)
		if err != nil {
			glog.Errorln(err)
		}
		//defer file.Close()
		size, err := io.Copy(file, resp.Body)
		if err != nil {
			glog.Errorln(err)
		}
		glog.Infoln("%s with %v bytes downloaded", path, size)

		// create table query
		NSEDeliveryPercentageTableName := fmt.Sprintf(NSEDeliveryPercentageTable, day, month, year)
		sqlQueryCreateTable := fmt.Sprintf(createTableQuery, NSEDeliveryPercentageTableName)
		rows, err := proddbhandle.Query(sqlQueryCreateTable)
		if err != nil {
			glog.Errorln(err)
		}
		if rows != nil {
			rows.Close()
		}

		// truncate tables if already exist
		sqlQueryTruncateTable := fmt.Sprintf(truncateTableQuery, NSEDeliveryPercentageTableName)
		rows, err = proddbhandle.Query(sqlQueryTruncateTable)
		if err != nil {
			glog.Errorln(err)
		}
		if rows != nil {
			rows.Close()
		}

		// Load csv into mysql
		sqlQueryLoadFIle := fmt.Sprintf(loadFileQuery, path, NSEDeliveryPercentageTableName, 4)
		rows, err = proddbhandle.Query(sqlQueryLoadFIle)
		if err != nil {
			glog.Errorln(err)
		}
		if rows != nil {
			rows.Close()
		}

		// delete the downloaded file
		cmd := exec.Command("rm", "-f", path)
		err = cmd.Run()
		if err != nil {

			glog.Error("Error removing the file: ", path)
		}

		// free the stuff
		//log.Println(resp)
		resp.Body.Close()
		file.Close()
		//req.Close
		count++
	}
}
*/
