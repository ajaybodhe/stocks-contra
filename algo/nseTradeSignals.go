package algo

import (
	"fmt"
	"github.com/ajaybodhe/stocks-contra/api"
	"github.com/ajaybodhe/stocks-contra/coreStructures"
	"github.com/ajaybodhe/stocks-contra/db"
	"github.com/ajaybodhe/stocks-contra/util"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NSESecuritiesBuySignal(dbHandle util.DB) error {
	// BUY SIGNAL ALGORITHM

	/* join full bhav copy daya, close price, % delivery data */
	nbdm, err := db.ReadNseBhavData(dbHandle, nil)
	if err != nil {
		fmt.Println("error executing ReadNseBhavData", err)
		return err
	}
	/* get eps, pe, pe-industry, 52 week high low */
	mcssCollection, err := db.ReadAllSecurityDetails(dbHandle, nil)
	if err != nil {
		fmt.Println("error executing ReadAllSecurityDetails", err)
		return err
	}
	bloomFilter,_ := db.GetInterestedSymbolsBloom(dbHandle)
	/* map of all securities corrected */
	nlsd := make(map[string]coreStructures.NseSecurityLongSignalData)
	/* check the stocks which have corrected for last few days n delivery % has gone up
	or either its above some comfort level */
	for symbol, nbda := range nbdm {
		
		if bloomFilter != nil {
			if bloomFilter.Test([]byte(symbol)) == false {
				continue
			}	
		}
		//fmt.Println("symbol:", symbol)
		//fmt.Println("nbda:", nbda)
		var max float32
		max = -1.0
		
		lenData := len(nbda)
		
		var rangeOver int
		if lenData >= 15 {
			rangeOver = 15
		} else if lenData >= 10 {
			rangeOver = 10
		} else {
			continue
		}
		for i:=1; i<=rangeOver; i++ {
			if max < nbda[lenData - i].ClosePrice {
				max = nbda[lenData - i].ClosePrice
				//maxInd = i
			}
		}
		
		correction := (max - nbda[lenData-1].ClosePrice) / max * 100
		//if symbol == "SIYSIL" || symbol == "DYNAMATECH" {
		//	fmt.Printf("\nsymbol=%v, max=%v, correction=%v, nbda[lenData-1].DelivPer=%v, nbda[lenData-1].ClosePrice=%v, nbda[lenData-2].ClosePrice=%v\n",
		//		symbol, max, correction, nbda[lenData-1].DelivPer, nbda[lenData-1].ClosePrice, nbda[lenData-2].ClosePrice)
		//}
		
		if max > nbda[lenData-1].ClosePrice &&
			nbda[lenData-1].DelivPer > 50 &&
			mcssCollection[symbol].EPS > 0.0 &&
			mcssCollection[symbol].BookValue > 0.0 &&
			correction > 5 &&
			mcssCollection[symbol].PE < 100 &&
			nbda[lenData-1].TtlTrdQnty > 1000 {
			
			/* we are assuming 5% diffrenerce in delivery percentages
			   this logic may change */
			if nbda[lenData-1].DelivPer >= (nbda[lenData-2].DelivPer - 5.0) && 
			   nbda[lenData-2].DelivPer >= (nbda[lenData-3].DelivPer - 5.0) &&
			   nbda[lenData-3].DelivPer >= (nbda[lenData-4].DelivPer - 5.0) &&
			   nbda[lenData-4].DelivPer >= (nbda[lenData-5].DelivPer - 5.0) {

				var highLowDiff float32
				
				if mcssCollection[symbol].High52 > mcssCollection[symbol].Low52 {
					highLowDiff = (mcssCollection[symbol].High52 - mcssCollection[symbol].Low52)
				} else {
					highLowDiff = 1.0
				}
				
				nlsd[symbol] = coreStructures.NseSecurityLongSignalData{
					Symbol:             symbol,
					PE:                 mcssCollection[symbol].PE,
					IndustryPE:         mcssCollection[symbol].IndustryPE,
					Correction:         correction,
					Closeness52WeekLow: (nbda[lenData-1].ClosePrice - mcssCollection[symbol].Low52) / highLowDiff * 100,
					DelivPer:           nbda[lenData-1].DelivPer,
					Sector:             mcssCollection[symbol].Sector,
					Strategy:           util.StocksContra,
				}
				continue
			}

			if  nbda[lenData-1].ClosePrice < nbda[lenData-2].ClosePrice &&
				nbda[lenData-2].ClosePrice < nbda[lenData-3].ClosePrice &&
				nbda[lenData-3].ClosePrice < nbda[lenData-4].ClosePrice &&
				nbda[lenData-4].ClosePrice < nbda[lenData-5].ClosePrice {
					
				var highLowDiff float32
				
				if mcssCollection[symbol].High52 > mcssCollection[symbol].Low52 {
					highLowDiff = (mcssCollection[symbol].High52 - mcssCollection[symbol].Low52)
				} else {
					highLowDiff = 1.0
				}
				
				nlsd[symbol] = coreStructures.NseSecurityLongSignalData{
					Symbol:             symbol,
					PE:                 mcssCollection[symbol].PE,
					IndustryPE:         mcssCollection[symbol].IndustryPE,
					Correction:         correction,
					Closeness52WeekLow: (nbda[lenData-1].ClosePrice - mcssCollection[symbol].Low52) / highLowDiff * 100,
					DelivPer:           nbda[lenData-1].DelivPer,
					Sector:             mcssCollection[symbol].Sector,
					Strategy:           util.FiveConsecutiveDaysDown,
				}
				continue
			}

			if correction > 5 {
				var highLowDiff float32
				if mcssCollection[symbol].High52 > mcssCollection[symbol].Low52 {
					highLowDiff = (mcssCollection[symbol].High52 - mcssCollection[symbol].Low52)
				} else {
					highLowDiff = 1.0
				}
				nlsd[symbol] = coreStructures.NseSecurityLongSignalData{
					Symbol:             symbol,
					PE:                 mcssCollection[symbol].PE,
					IndustryPE:         mcssCollection[symbol].IndustryPE,
					Correction:         correction,
					Closeness52WeekLow: (nbda[lenData-1].ClosePrice - mcssCollection[symbol].Low52) / highLowDiff * 100,
					DelivPer:           nbda[lenData-1].DelivPer,
					Sector:             mcssCollection[symbol].Sector,
					Strategy:           util.MajorCorrection,
				}
				continue
			}

		}
	}

	/* rank the stocks which have been corrected based on
	% of correction, PE vs industry PE, % of delivery,
	closeness % to 52 week low=(current price - 52 week low) /(52 week high - 52 week low)*100
	*/
	// TBD this algorithm

	// store the map in database
	err = db.WriteNSESecuritiesBuySignal(dbHandle, nlsd)
	if err != nil {
		fmt.Println("error executing NSESecuritiesBuySignal", err)
		return err
	}

	// SELL SIGNAL ALGORITHM

	return nil
}

/* poll current NSE order book */
func NseOrderBookAnalyser(client *http.Client, dbHandle util.DB) error {
	symbolStrategyMap, err := db.RetrieveAllSymbolsNStrategy(dbHandle)
	if err != nil {
		return err
	}
	// use sync.waitgroup to wait here
	done := make(chan bool)
	for k, v := range symbolStrategyMap {
		//if v != util.InvalidStrategy {
		go NseLiveQuoteAnalyser(done, client, k, v)
		///time.Sleep(10 * time.Second)
		//}
	}
	// sleep till 3.30 pm
	time.Sleep(2 * time.Hour)
	close(done)
	return nil
}

// all of these function should run between 9.15 am to 3.30
// yithen they should be destroyed by a signal on channel
func NseLiveQuoteAnalyser(done <-chan bool, client *http.Client, symbol string, strategy int) {
	///fmt.Println("inside NseLiveQuoteAnalyser:", symbol)
	for {
		select {
		case <-time.After(5 * time.Minute):
			fmt.Println("\n\nSYMBOL IS:", symbol)
			nseLQD := api.GetNSELiveQuote(client, symbol)
			fmt.Println("%d	%v	%v\n", len(nseLQD.Data), nseLQD.Data[0].TotalBuyQuantity, nseLQD.Data[0].TotalSellQuantity)
			if nseLQD != nil &&
				nseLQD.Data != nil &&
				len(nseLQD.Data) >= 1 {
				totalBuyQty := strings.Replace(nseLQD.Data[0].TotalBuyQuantity, util.CommaChar, util.EmptyString, -1)
				totalSellQty := strings.Replace(nseLQD.Data[0].TotalSellQuantity, util.CommaChar, util.EmptyString, -1)
				tBQ, _ := strconv.Atoi(totalBuyQty)
				tSQ, _ := strconv.Atoi(totalSellQty)
				if tBQ != 0 && tSQ != 0 {
					if ((tBQ - tSQ) / tSQ * 100) > 10 {
						//((nseLQD.Data[0].TotalBuyQuantity-nseLQD.Data[0].TotalSellQuantity)/nseLQD.Data[0].TotalSellQuantity*100) > 20 {
						fmt.Println("\nBuy symbol:", symbol, "	strategy:", strategy, "TotalBuyQuantity", nseLQD.Data[0].TotalBuyQuantity, "TotalSellQuantity", nseLQD.Data[0].TotalSellQuantity, "\n")
						//fmt.Printf("\n%v\n", nseLQD)
					} else if ((tSQ - tBQ) / tSQ * 100) > 10 {
						//((nseLQD.Data[0].TotalBuyQuantity-nseLQD.Data[0].TotalSellQuantity)/nseLQD.Data[0].TotalSellQuantity*100) > 20 {
						fmt.Println("\nSell symbol:", symbol, "	strategy:", strategy, "TotalBuyQuantity", nseLQD.Data[0].TotalBuyQuantity, "TotalSellQuantity", nseLQD.Data[0].TotalSellQuantity, "\n")
						//fmt.Printf("\n%v\n", nseLQD)
					}
				}
			}
		case <-done:
			//break
			return
		}
	}
}
