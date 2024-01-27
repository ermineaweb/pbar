package pbar

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

type ConfigPbar struct {
	TotalTasks           uint16
	CharDone             rune
	CharTodo             rune
	ColorPercentWorking  color
	ColorPercentFinished color
	ColorCharDone        color
	ColorCharTodo        color
}

type pbar struct {
	total       uint16
	actual      uint16
	width       uint16
	charDone    rune
	charTodo    rune
	colorPct    color
	colorPctEnd color
	colorDone   color
	colorTodo   color
	startedAt   time.Time
	signalCh    chan os.Signal
	lock        sync.Mutex
}

type windowSize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func NewPbar(cfg ConfigPbar) *pbar {
	pb := &pbar{
		total:       1,
		charDone:    '#',
		charTodo:    '-',
		colorPct:    YELLOW,
		colorPctEnd: GREEN,
		colorDone:   default_color,
		colorTodo:   default_color,
		signalCh:    make(chan os.Signal, 1),
		lock:        sync.Mutex{},
	}

	pb.customPbarConfig(cfg)
	pb.handleSignals()
	pb.updateSize()
	pb.Add(0)

	return pb
}

func (pb *pbar) customPbarConfig(cfg ConfigPbar) {
	if cfg.TotalTasks > 0 {
		pb.total = cfg.TotalTasks
	}
	if cfg.CharDone != 0 {
		pb.charDone = cfg.CharDone
	}
	if cfg.CharTodo != 0 {
		pb.charTodo = cfg.CharTodo
	}
	if cfg.ColorPercentWorking != "" {
		pb.colorPct = cfg.ColorPercentWorking
	}
	if cfg.ColorPercentFinished != "" {
		pb.colorPctEnd = cfg.ColorPercentFinished
	}
	if cfg.ColorCharDone != "" {
		pb.colorDone = cfg.ColorCharDone
	}
	if cfg.ColorCharTodo != "" {
		pb.colorTodo = cfg.ColorCharTodo
	}
}

func (pb *pbar) Add(increment int) {
	if pb.width == 0 {
		return
	}

	pb.lock.Lock()
	defer pb.lock.Unlock()

	if pb.actual == 0 {
		pb.startedAt = time.Now()
	}

	if pb.actual += uint16(increment); pb.actual > pb.total {
		return
	}

	percent := int(float64(pb.actual) / float64(pb.total) * 100.0)
	pb.render(percent)
}

func (pb *pbar) render(percent int) {
	fmt.Printf("\r%v", delete_line)

	var pbar string
	if pb.actual >= pb.total {
		pbar = fmt.Sprintf("[%s%3d%%%s] [%4d/%-4d] [%4s]", pb.colorPctEnd, percent, default_color, pb.actual, pb.total, time.Since(pb.startedAt).Truncate(time.Second))
	} else {
		pbar = fmt.Sprintf("[%s%3d%%%s] [%4d/%-4d] [%4s]", pb.colorPct, percent, default_color, pb.actual, pb.total, time.Since(pb.startedAt).Truncate(time.Second))
	}

	widthTotal := int(pb.width) - len(pbar) + 10
	widthDone := int(float64(widthTotal) * float64(pb.actual) / float64(pb.total))
	done := strings.Repeat(string(pb.charDone), widthDone)
	todo := strings.Repeat(string(pb.charTodo), widthTotal-widthDone)

	if pb.width > 30 {
		pbar = fmt.Sprintf("%s [%s%s%s%s%s%s]", pbar, pb.colorDone, done, default_color, pb.colorTodo, todo, default_color)
	}

	fmt.Print(pbar)

	if pb.actual >= pb.total {
		fmt.Println()
	}
}

func (pb *pbar) updateSize() {
	winSize := &windowSize{}

	if _, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(winSize))); err != 0 {
		if err == syscall.ENOTTY || err == syscall.ENODEV {
			return
		} else {
			return
		}
	}

	pb.width = winSize.Col
}

func (pb *pbar) handleSignals() {
	signal.Notify(pb.signalCh, syscall.SIGWINCH, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			select {
			case sig := <-pb.signalCh:
				switch sig {
				case syscall.SIGWINCH:
					// terminal is resized
					pb.updateSize()

				case syscall.SIGTERM, syscall.SIGINT:
					os.Exit(1)
				}
			}
		}
	}()
}
