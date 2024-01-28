package pbar

import (
	"fmt"
	"time"
)

type SpinnerType []string

var (
	SPINNER_ARROW        = SpinnerType{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}
	SPINNER_POINTS       = SpinnerType{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	SPINNER_SIMPLEBAR    = SpinnerType{"|", "/", "-", "\\"}
	SPINNER_HISTOGRAM    = SpinnerType{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃"}
	SPINNER_PING         = SpinnerType{"▖", "▘", "▝", "▗"}
	SPINNER_MOON         = SpinnerType{"◐", "◓", "◑", "◒"}
	SPINNER_PENDULUM     = SpinnerType{"■ ■■ ", " ■■■ ", " ■■ ■", " ■■■ "}
	SPINNER_BIG_PENDULUM = SpinnerType{"■  ■■■■  ", " ■ ■■■■  ", "  ■■■■■  ", "  ■■■■ ■ ", "  ■■■■  ■", "  ■■■■ ■ ", "  ■■■■■  ", " ■ ■■■■  "}
	SPINNER_SQUARE       = SpinnerType{"■   ", "■■  ", "■■■ ", "■■■■", "■■■ ", "■■  "}
	SPINNER_SNAKE        = SpinnerType{"■   ", "■■  ", "■■■ ", "■■■■", " ■■■", "  ■■", "   ■"}
)

type ConfigSpinner struct {
	Spinner          SpinnerType
	StartMessage     string
	StopMessage      string
	ColorSpinner     color
	ColorTimer       color
	AnimationDelayMs uint16
}

type spinner struct {
	spType      SpinnerType
	startMsg    string
	stopMsg     string
	colorSp     color
	colorTimer  color
	animDelayMs uint16
	startedAt   time.Time
	stopCh      chan bool
	doneCh      chan bool
}

func NewDefaultSpinner() *spinner {
	return NewCustomSpinner(ConfigSpinner{})
}

func NewCustomSpinner(cfg ConfigSpinner) *spinner {
	sp := &spinner{
		spType:      SPINNER_PENDULUM,
		startMsg:    "Work in progress since:",
		stopMsg:     "Work finished in",
		colorSp:     BLUE_BRIGHT,
		colorTimer:  BLUE_BRIGHT,
		animDelayMs: 120,
		stopCh:      make(chan bool, 1),
		doneCh:      make(chan bool, 1),
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
	if cfg.ColorSpinner != "" {
		sp.colorSp = cfg.ColorSpinner
	}
	if cfg.ColorTimer != "" {
		sp.colorTimer = cfg.ColorTimer
	}
	if cfg.AnimationDelayMs != 0 {
		sp.animDelayMs = cfg.AnimationDelayMs
	}
}

func (sp *spinner) Start() {
	go func() {
		sp.startedAt = time.Now()

		for i := 0; i < len(sp.spType); i++ {
			fmt.Printf("\r%v", delete_line)

			select {
			case <-sp.stopCh:
				fmt.Printf("    %s %s%s%s\n", sp.stopMsg, sp.colorTimer, time.Since(sp.startedAt).Truncate(time.Millisecond), default_color)
				sp.doneCh <- true
				return
			default:
				fmt.Printf("  %s%v%s %s %s%s%s", sp.colorSp, sp.spType[i], default_color, sp.startMsg, sp.colorTimer, time.Since(sp.startedAt).Truncate(time.Second), default_color)
				time.Sleep(time.Duration(sp.animDelayMs) * time.Millisecond)
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
