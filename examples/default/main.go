package main

import (
	"math/rand"
	"time"

	"github.com/ermineaweb/pbar"
)

func main() {
	use := "progress"
	tasks := 150

	switch use {

	case "spinner":
		sp := pbar.NewDefaultSpinner()

		sp.Start()
		defer sp.Stop()

		for i := 1; i <= tasks; i++ {
			work(i)
		}

	case "progress":
		pbar := pbar.NewDefaultPbar(tasks)

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
