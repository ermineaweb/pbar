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

type ConfigSpinner struct {
	Spinner      SpinnerType
	StartMessage string
	StopMessage  string
}

type spinner struct {
	spType    SpinnerType
	startMsg  string
	stopMsg   string
	startedAt time.Time
	stopCh    chan bool
	doneCh    chan bool
}

func NewSpinner(cfg ConfigSpinner) *spinner {
	sp := &spinner{
		spType:   SPINNER_POINTS,
		startMsg: "Work in progress since:",
		stopMsg:  "Work finished in",
		stopCh:   make(chan bool, 1),
		doneCh:   make(chan bool, 1),
	}

	sp.customSpinnerConfig(cfg)
	return sp
}

func (sp *spinner) customSpinnerConfig(cfg ConfigSpinner) {
	if len(cfg.Spinner) > 0 {
		sp.spType = cfg.Spinner
	}
	if cfg.StartMessage != "" {
		sp.startMsg = cfg.StartMessage
	}
	if cfg.StopMessage != "" {
		sp.stopMsg = cfg.StopMessage
	}
}

func (sp *spinner) Start() {
	sp.startedAt = time.Now()

	go func() {
		for i := 0; i < len(sp.spType); i++ {
			fmt.Printf("\r%v", DELETE_LINE)

			select {
			case <-sp.stopCh:
				fmt.Printf("    %s %s%s%s\n", sp.stopMsg, YELLOW_BOLD, time.Since(sp.startedAt).Truncate(time.Millisecond), RESET_COLOR)
				sp.doneCh <- true
				return
			default:
				fmt.Printf("  %s%v%s %s %s%s%s", YELLOW_BOLD, sp.spType[i], RESET_COLOR, sp.startMsg, YELLOW_BOLD, time.Since(sp.startedAt).Truncate(time.Second), RESET_COLOR)
				time.Sleep(100 * time.Millisecond)
				if i+1 == len(sp.spType) {
					i = -1
				}
			}
		}
	}()
}

func (sp *spinner) Stop() {
	sp.stopCh <- true
	<-sp.doneCh
}
