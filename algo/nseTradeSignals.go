package algo

import (
	"fmt"
	"github.com/ajaybodhe/stocks-contra/coreStructures"
	"github.com/ajaybodhe/stocks-contra/db"
	"github.com/ajaybodhe/stocks-contra/util"
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

	/* map of all securities corrected */
	nlsd := make(map[string]coreStructures.NseSecurityLongSignalData)
	/* check the stocks which have corrected for last few days n delivery % has gone up
	or either its above some comfort level */
	for symbol, nbda := range nbdm {
		//fmt.Println("symbol:", symbol)
		//fmt.Println("nbda:", nbda)
		var max float32
		//var maxInd int
		max = -1.0
		//maxInd = -1

		for _, nbr := range nbda {
			if max < nbr.ClosePrice {
				max = nbr.ClosePrice
				//maxInd = i
			}
		}
		lenData := len(nbda)
		//fmt.Println("lendata", lenData)
		correction := (max - nbda[lenData-1].ClosePrice) / max * 100
		if max > nbda[lenData-1].ClosePrice &&
			nbda[lenData-1].DelivPer > 50 &&
			mcssCollection[symbol].EPS > 0.0 &&
			mcssCollection[symbol].BookValue > 0.0 &&
			correction > 2.5 {
			/* create a map of stocks which are trading lower with delivery % getting higher */
			var eightDayAvg, fiveDayAvg, threeDayAvg float32
			if lenData >= 3 {
				threeDayAvg = (nbda[lenData-1].DelivPer + nbda[lenData-2].DelivPer + nbda[lenData-3].DelivPer) / 3
			}
			if lenData >= 5 {
				fiveDayAvg = ((threeDayAvg * 3) + nbda[lenData-4].DelivPer + nbda[lenData-5].DelivPer) / 5
			}
			if lenData >= 8 {
				eightDayAvg = ((fiveDayAvg * 5) + nbda[lenData-6].DelivPer + nbda[lenData-7].DelivPer + nbda[lenData-8].DelivPer) / 5
			}
			/* we are assuming 5% diffrenerce in delivery percentages
			   this logic may change */
			if threeDayAvg > (fiveDayAvg-5) &&
				threeDayAvg > (eightDayAvg-5) &&
				nbda[0].DelivPer > (threeDayAvg-5) {

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
				}

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

	/* poll current NSE order book */

	// SELL SIGNAL ALGORITHM

	return nil
}
