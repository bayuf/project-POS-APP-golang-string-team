package utils

import (
	"fmt"
	"sync"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
)

type EmailJob struct {
	Payload *dto.Email
}

type EmailSender interface {
	SendEmail(payload dto.Email) error
}

// Worker pool
func StartEmailWorkers(workerCount int, jobs <-chan EmailJob, stop <-chan struct{}, wg *sync.WaitGroup, emailUc EmailSender) {
	wg.Add(workerCount)

	for i := 1; i <= workerCount; i++ {
		go func(id int) {
			defer wg.Done()

			for {
				select {
				case job, ok := <-jobs:
					if !ok {
						fmt.Println("worker", id, "jobs channel closed")
						return
					}

					// send email
					if err := emailUc.SendEmail(*job.Payload); err != nil {
						fmt.Println("worker", id, "failed to send email:", err)
					}
				case <-stop:
					fmt.Println("worker", id, "received stop signal")
					return
				}
			}
		}(i)
	}
}
