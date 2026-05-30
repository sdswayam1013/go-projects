package worker

import "sync"

type Job struct{
	ID string
	Duration int
}

type JobStatus string
	 
const (
	Pending JobStatus = "pending"
	Completed JobStatus = "completed"
)

type Result struct{
	JobID string
	WorkerID string
	Status JobStatus
}

func main(){


	jobs := make(chan Job, 5)
	results := make(chan Result, 5)

	var wg sync.WaitGroup

	jobBackLog := make([]Job,10,20)
	for i :=0; i<10;i++{
		jobBackLogObj := Job( ID: "Job"+string(i), Duration: 8)
		jobBackLog = append(jobBackLog, jobBackLogObj)
		//obBackLog[i] = jobBackLog[]



	}
}

