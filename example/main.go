package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/ermineaweb/pbar"
)

func main() {
	var wg sync.WaitGroup

	use := "progress"
	tasks := 25

	switch use {

	case "spinner":
		spinner := pbar.NewSpinner(
			pbar.ConfigSpinner{
				Spinner:      pbar.SPINNER_ARROW,
				StartMessage: "Let's work!",
				StopMessage:  "Job's done!",
			},
		)

		spinner.Start()
		defer spinner.Stop()

		for i := 1; i <= tasks; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				work(i)
			}(i)
		}

	case "progress":
		pbar := pbar.NewPbar(
			pbar.ConfigPbar{
				TotalTasks: uint16(tasks),
				CharDone:   '#',
				CharTodo:   '-',
			},
		)

		for i := 1; i <= tasks; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				work(i)
				pbar.Add(1)
			}(i)
		}
	}

	wg.Wait()
}

func work(i int) {
	rnd := rand.Intn(8000) + 1000
	time.Sleep(time.Duration(rnd) * time.Millisecond)
}
