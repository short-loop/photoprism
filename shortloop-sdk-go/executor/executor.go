package executor

import (
	"time"
)

type ScheduledExecutor struct {
	quit chan int
}

func NewScheduledExecutor() ScheduledExecutor {
	se := ScheduledExecutor{
		quit: make(chan int, 1),
	}
	return se
}

func (se *ScheduledExecutor) Schedule(task func(), delay time.Duration) {
	go func() {
		select {
		case <-se.quit:
			//fmt.Println("exiting task")
			return
		case <-time.After(delay):
			//fmt.Println("executing task")
			task()
			break
		}
	}()
}

func (se *ScheduledExecutor) ScheduleAtFixedRate(task func(), delay time.Duration) {
	ticker := time.NewTicker(delay)
	go func() {
		defer func() {
			ticker.Stop()
		}()
		for {
			select {
			case <-se.quit:
				return
			case <-ticker.C:
				go task()
			}
		}
	}()
}

func (se *ScheduledExecutor) Shutdown() error {
	go func() {
		se.quit <- 1
	}()
	return nil
}
