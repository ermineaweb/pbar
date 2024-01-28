package main

import (
	"math/rand"
	"time"

	"github.com/ermineaweb/pbar"
)

func main() {
	use := "spinner"
	tasks := 150

	switch use {

	case "spinner":
		spinner := pbar.NewCustomSpinner(
			pbar.ConfigSpinner{
				Spinner:          pbar.SPINNER_BIG_PENDULUM,
				StartMessage:     "Let's work!",
				StopMessage:      "Job's done!",
				ColorSpinner:     pbar.RED_BRIGHT,
				ColorTimer:       pbar.BLUE_BRIGHT,
				AnimationDelayMs: 140,
			},
		)

		spinner.Start()
		defer spinner.Stop()

		for i := 1; i <= tasks; i++ {
			work(i)
		}

	case "progress":
		pbar := pbar.NewCustomPbar(
			pbar.ConfigPbar{
				TotalTasks:           uint64(tasks),
				Infos:                true,
				CharDone:             'o',
				CharTodo:             '-',
				ColorPercentWorking:  pbar.RED_BRIGHT,
				ColorPercentFinished: pbar.GREEN,
				ColorCharDone:        pbar.RED_BRIGHT,
				ColorCharTodo:        pbar.BLACK_BRIGHT,
			},
		)

		for i := 1; i <= tasks; i++ {
			work(i)
			pbar.Add(1)

		}
	}
}

func work(i int) {
	rnd := rand.Intn(100) + 20
	time.Sleep(time.Duration(rnd) * time.Millisecond)
}
