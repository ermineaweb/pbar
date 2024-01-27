package pbar

import (
	"fmt"
	"time"
)

type SpinnerType []string

var (
	SPINNER_ARROW  = SpinnerType{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}
	SPINNER_POINTS = SpinnerType{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
)

type Spinner struct {
	spType       SpinnerType
	startedAt    time.Time
	startMessage string
	stopMessage  string
	stopCh       chan bool
	doneCh       chan bool
}

func NewSpinner(spType SpinnerType) *Spinner {
	return &Spinner{
		spType:       spType,
		startMessage: "Work in progress since:",
		stopMessage:  "Work finished in",
		stopCh:       make(chan bool, 1),
		doneCh:       make(chan bool, 1),
	}
}

func (sp *Spinner) Start() {
	sp.startedAt = time.Now()

	go func() {
		for i := 0; i < len(sp.spType); i++ {
			fmt.Printf("\r%v", DELETE_LINE)

			select {
			case <-sp.stopCh:
				fmt.Printf("    %s %s%s%s\n", sp.stopMessage, YELLOW_BOLD, time.Since(sp.startedAt).Truncate(time.Millisecond), RESET_COLOR)
				sp.doneCh <- true
				return
			default:
				fmt.Printf("  %s%v%s %s %s%s%s", YELLOW_BOLD, sp.spType[i], RESET_COLOR, sp.startMessage, YELLOW_BOLD, time.Since(sp.startedAt).Truncate(time.Second), RESET_COLOR)
				time.Sleep(100 * time.Millisecond)
				if i+1 == len(sp.spType) {
					i = -1
				}
			}
		}
	}()
}

func (sp *Spinner) Stop() {
	sp.stopCh <- true
	<-sp.doneCh
}
