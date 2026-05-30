package worker

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID       string
	Duration int
}

type Result struct {
	JobID    string
	WorkerID int
	Status   string
}

// Worker function
func worker(workerID int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d picked %s\n", workerID, job.ID)

		time.Sleep(time.Duration(job.Duration) * time.Second)

		results <- Result{
			JobID:    job.ID,
			WorkerID: workerID,
			Status:   "Done",
		}
	}
}

// Producer: moves jobs from backlog → jobs channel
func produceJobs(jobBacklog []Job, jobs chan<- Job) {
	for _, job := range jobBacklog {
		fmt.Printf("Enqueuing %s\n", job.ID)
		jobs <- job
	}
	close(jobs)
}

func main() {
	numWorkers := 3

	// Step 1: Create backlog with string Job IDs
	jobBacklog := []Job{
		{ID: "Job1", Duration: 2},
		{ID: "Job2", Duration: 1},
		{ID: "Job3", Duration: 3},
		{ID: "Job4", Duration: 2},
		{ID: "Job5", Duration: 1},
	}

	jobs := make(chan Job, len(jobBacklog))
	results := make(chan Result, len(jobBacklog))

	var wg sync.WaitGroup

	// Step 2: Start workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	// Step 3: Start producer
	go produceJobs(jobBacklog, jobs)

	// Step 4: Close results after workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Step 5: Read results
	for res := range results {
		fmt.Printf("%s completed by Worker %d\n", res.JobID, res.WorkerID)
	}
}
