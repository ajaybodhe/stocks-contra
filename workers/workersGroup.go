package workers

import (
	"net/http"
	"crypto/tls"
)

const (
	MaxWorker = 10
	MaxQueue  = 10
)

// Job represents the job to be run
type Job struct {
	Payload interface{}
}

// function to be executed
type JobFunction func(job Job, client *http.Client)

// A buffered channel that we can send work requests on.
//var JobQueue chan Job

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool  chan chan Job
	JobChannel  chan Job
	quit    	chan bool
	JF JobFunction
	Client *http.Client
}

func NewWorker(workerPool chan chan Job, j JobFunction) *Worker {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return &Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	JF:j,
	Client: client}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w *Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel
			
			select {
			case job := <-w.JobChannel:
			// we have received a work request.
			// TODO execute the job
				w.JF(job, w.Client)
			//	if err := job.Payload.UploadToS3(); err != nil {
			//		log.Errorf("Error uploading to S3: %s", err.Error())
			//	}
			
			case <-w.quit:
			// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}


type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
	MaxWorkers int
	JF JobFunction
}

func NewDispatcher(maxWorkers int, j JobFunction) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool, JF:j, MaxWorkers:maxWorkers}
}

func (d *Dispatcher) Run(jobQueue chan Job) {
	// starting n number of workers
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool, d.JF)
		worker.Start()
	}
	
	go d.dispatch(jobQueue)
}

func (d *Dispatcher) dispatch(jobQueue chan Job) {
	for {
		select {
		case job := <-jobQueue:
		// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool
				
				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}

func CreateJobQueue(maxQueue int) (chan Job){
	jobQueue :=make(chan Job, maxQueue)
	return jobQueue
}

func init() {
}