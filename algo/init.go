package algo

import (
	"github.com/ajaybodhe/stocks-contra/workers"
)

var NseOrderBookQueue chan workers.Job
//var LastAnnouncementDate string
//var LastAnnouncementDateMutex *sync.Mutex

func init() {
	NseOrderBookQueue = workers.CreateJobQueue(workers.MaxQueue)
	//LastAnnouncementDateMutex = &sync.Mutex{}
}