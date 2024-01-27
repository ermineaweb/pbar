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
	TotalTasks uint16
	CharDone   rune
	CharTodo   rune
}

type pbar struct {
	total     uint16
	actual    uint16
	width     uint16
	charDone  rune
	charTodo  rune
	startedAt time.Time
	signalCh  chan os.Signal
	lock      sync.Mutex
}

type windowSize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func NewPbar(cfg ConfigPbar) *pbar {
	pb := &pbar{
		total:    1,
		charDone: '#',
		charTodo: '-',
		signalCh: make(chan os.Signal, 1),
		lock:     sync.Mutex{},
	}

	pb.customPbarConfig(cfg)
	pb.handleResize()
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

	var pbar string

	if pb.width > 28 {
		widthTotal := int(pb.width - 28)
		widthDone := int(float64(widthTotal) * float64(pb.actual) / float64(pb.total))
		done := strings.Repeat(string(pb.charDone), widthDone)
		todo := strings.Repeat(string(pb.charTodo), widthTotal-widthDone)
		pbar = fmt.Sprintf(" [%s%s]", done, todo)
	}

	fmt.Printf("\r%v", DELETE_LINE)

	if pb.actual >= pb.total {
		fmt.Printf("[%s%3d%%%s] [%4d/%4d] [%4s]%s\n", GREEN_BOLD, percent, RESET_COLOR, pb.actual, pb.total, time.Since(pb.startedAt).Truncate(time.Second), pbar)
	} else {
		fmt.Printf("[%s%3d%%%s] [%4d/%4d] [%4s]%s", YELLOW_BOLD, percent, RESET_COLOR, pb.actual, pb.total, time.Since(pb.startedAt).Truncate(time.Second), pbar)
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

func (pb *pbar) handleResize() {
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
