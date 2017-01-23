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
	"github.com/ajaybodhe/stocks-contra/workers"
)

func NSESecuritiesBuySignal() error {
	// BUY SIGNAL ALGORITHM

	/* join full bhav copy daya, close price, % delivery data */
	nbdm, err := db.ReadNseBhavData(nil)
	if err != nil {
		fmt.Println("error executing ReadNseBhavData", err)
		return err
	}
	/* get eps, pe, pe-industry, 52 week high low */
	mcssCollection, err := db.ReadAllSecurityDetails(nil)
	if err != nil {
		fmt.Println("error executing ReadAllSecurityDetails", err)
		return err
	}
	bloomFilter,_ := db.GetInterestedSymbolsBloom()
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
	err = db.WriteNSESecuritiesBuySignal(nlsd)
	if err != nil {
		fmt.Println("error executing NSESecuritiesBuySignal", err)
		return err
	}

	// SELL SIGNAL ALGORITHM

	return nil
}

func NseOrderBookAnalyser() error {
	
	dispatcher := workers.NewDispatcher(workers.MaxWorker, NseLiveQuoteAnalyser)
	
	//TODO
	// algo/nseTradeSignals.go:187: cannot use [2]string literal (type [2]string) as type []string in argument to db.GetInterestedSymbols
	tables :=make([]string, 2)
	tables[0]=util.NIFTY_50
	tables[1]=util.NIFTY_NEXT_50
	
	// read all symbols from NSE_50 and NSE_Next_50
	symbols, err := db.GetInterestedSymbols(tables)
	if err != nil {
		return err
	}
	
	// create jobqueue
	jq := workers.CreateJobQueue(workers.MaxQueue)
	
	// run the dispatcher
	dispatcher.Run(jq)

	for {
		for i := range symbols {
			jq<-workers.Job{Payload:symbols[i]}
		}
		
		time.Sleep(2 * time.Minute)
	}
	return nil
}

// all of these function should run between 9.15 am to 3.30
// yithen they should be destroyed by a signal on channel
func NseLiveQuoteAnalyser(job workers.Job, client *http.Client) {
	symbol := job.Payload.(string)
	nseLQD := api.GetNSELiveQuote(client, symbol)
	
	if nseLQD != nil &&
		nseLQD.Data != nil &&
		len(nseLQD.Data) >= 1 {

		totalBuyQty := strings.Replace(nseLQD.Data[0].TotalBuyQuantity, util.CommaChar, util.EmptyString, -1)
		totalSellQty := strings.Replace(nseLQD.Data[0].TotalSellQuantity, util.CommaChar, util.EmptyString, -1)

		tBQ, _ := strconv.Atoi(totalBuyQty)
		tSQ, _ := strconv.Atoi(totalSellQty)
		
		if tBQ != 0 && tSQ != 0 {
			fmt.Printf("step 4, symbol=%v, tbq=%v, tsq=%v\n", symbol, tBQ, tSQ)
			if ((tBQ - tSQ) / tSQ * 100) > 20 {
				
				//((nseLQD.Data[0].TotalBuyQuantity-nseLQD.Data[0].TotalSellQuantity)/nseLQD.Data[0].TotalSellQuantity*100) > 20 {
				fmt.Println("\nBuy symbol:", symbol, "TotalBuyQuantity", nseLQD.Data[0].TotalBuyQuantity, "TotalSellQuantity", nseLQD.Data[0].TotalSellQuantity, "\n")
				//fmt.Printf("\n%v\n", nseLQD)
			} else if ((tSQ - tBQ) / tSQ * 100) > 50 {
				//((nseLQD.Data[0].TotalBuyQuantity-nseLQD.Data[0].TotalSellQuantity)/nseLQD.Data[0].TotalSellQuantity*100) > 20 {
				fmt.Println("\nSell symbol:", symbol, "TotalBuyQuantity", nseLQD.Data[0].TotalBuyQuantity, "TotalSellQuantity", nseLQD.Data[0].TotalSellQuantity, "\n")
				//fmt.Printf("\n%v\n", nseLQD)
			}
		}
	}

}
