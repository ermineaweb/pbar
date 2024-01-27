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

type Pbar struct {
	total     int
	actual    int
	width     int
	startedAt time.Time
	charDone  rune
	charTodo  rune
	signalCh  chan os.Signal
	lock      sync.Mutex
}

type windowSize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func NewPbar(total int) *Pbar {
	pb := &Pbar{
		total:    total,
		charDone: '#',
		charTodo: '-',
		signalCh: make(chan os.Signal, 1),
		lock:     sync.Mutex{},
	}

	pb.handleResize()
	pb.updateSize()
	pb.Add(0)
	return pb
}

func (pb *Pbar) Add(increment int) {
	if pb.width == 0 {
		return
	}

	pb.lock.Lock()
	defer pb.lock.Unlock()

	if pb.actual == 0 {
		pb.startedAt = time.Now()
	}

	pb.actual += increment
	percent := int(float64(pb.actual) / float64(pb.total) * 100.0)

	fmt.Printf("\r%v", DELETE_LINE)

	var pbar string

	if pb.width > 28 {
		widthTotal := pb.width - 28
		widthDone := int(float64(widthTotal) * float64(pb.actual) / float64(pb.total))
		done := strings.Repeat(string(pb.charDone), widthDone)
		todo := strings.Repeat(string(pb.charTodo), widthTotal-widthDone)
		pbar = fmt.Sprintf(" [%s%s]", done, todo)
	}

	if pb.actual >= pb.total {
		fmt.Printf("[%s%3d%%%s] [%4d/%4d] [%4s]%s", GREEN_BOLD, percent, RESET_COLOR, pb.actual, pb.total, time.Since(pb.startedAt).Truncate(time.Second), pbar)
	} else {
		fmt.Printf("[%s%3d%%%s] [%4d/%4d] [%4s]%s", YELLOW_BOLD, percent, RESET_COLOR, pb.actual, pb.total, time.Since(pb.startedAt).Truncate(time.Second), pbar)
	}
}

func (pb *Pbar) updateSize() {
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

	pb.width = int(winSize.Col)
}

func (pb *Pbar) handleResize() {
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
