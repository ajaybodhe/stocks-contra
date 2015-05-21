package api

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/ajaybodhe/stocks-contra/util"
	"github.com/golang/glog"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

func GetMoneycontrolLiveQuote(client *http.Client, symbol string) {
	quoteURL := getMoneycontrolLiveQuoteURL(client, symbol)
	//args := fmt.Sprintf("-dump >/tmp/quote", quoteURL)
	cmd := exec.Command("/usr/bin/lynx", "-dump >/tmp/quote", quoteURL)
	out, err := cmd.Output()
	//err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		glog.Error("Error downloading the quote: ", quoteURL)
	}
	quoteStr := string(out)
	strings.Index(quoteStr, "")
	/*
		err = ioutil.WriteFile("/tmp/quote", out, os.ModePerm)
		if err != nil {
			glog.Errorln(":Result:Fail:Error:", err.Error())
			return
		}
	*/
}

func getMoneycontrolLiveQuoteURL(client *http.Client, symbol string) string {
	//symbol := "ABB"
	symbolComma := symbol + ","
	symbolForwardSlah := symbol + "/"
	reqURL := fmt.Sprintf(util.MoneyControlURLFetcher, symbol, "%20")
	fmt.Println("reqURL", reqURL)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		glog.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:28.0) Gecko/20100101 Firefox/28.0")
	req.Header.Set("Host", "www.moneycontrol.com")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	/* get the response */
	resp, err := client.Do(req)
	if err != nil {
		glog.Errorln(":Result:Fail:Error:", err.Error())
		return ""
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
	quoteDataStr := buf.String()

	subStr := strings.Split(quoteDataStr, "a href=\"")
	for _, each := range subStr {
		//fmt.Println(index, each)
		if strings.Contains(each, symbolComma) || strings.Contains(each, symbolForwardSlah) {
			quoteURL := strings.Split(each, "\">")
			fmt.Println(quoteURL[0])
			return quoteURL[0]
		}

	}

	return ""
}
