package utils

import (
	"code.google.com/p/go.crypto/ssh/terminal"
	"fmt"
	"github.com/cheggaaa/pb"
	"syscall"
)

const (
	codePrint = iota
	codeProgress
	codeHideProgress
)

type printTask struct {
	code    int
	message string
}

// Progress is a progress displaying subroutine, it allows to show download and other operations progress
// mixed with progress bar
type Progress struct {
	stop     chan bool
	stopped  chan bool
	queue    chan printTask
	bar      *pb.ProgressBar
	barShown bool
}

// NewProgress creates new progress instance
func NewProgress() *Progress {
	return &Progress{
		stop:    make(chan bool),
		stopped: make(chan bool),
		queue:   make(chan printTask, 100),
	}
}

// Start makes progress start its work
func (p *Progress) Start() {
	go p.worker()
}

// Shutdown shuts down progress display
func (p *Progress) Shutdown() {
	p.ShutdownBar()
	p.stop <- true
	<-p.stopped
}

// InitBar creates progressbar for count bytes
func (p *Progress) InitBar(count int64, isBytes bool) {
	if p.bar != nil {
		panic("bar already initialized")
	}
	if terminal.IsTerminal(syscall.Stdout) {
		p.bar = pb.New(0)
		p.bar.Total = count
		p.bar.NotPrint = true
		p.bar.Callback = func(out string) {
			p.queue <- printTask{code: codeProgress, message: out}
		}

		if isBytes {
			p.bar.SetUnits(pb.U_BYTES)
			p.bar.ShowSpeed = true
		}
		p.barShown = false
		p.bar.Start()
	}
}

// ShutdownBar stops progress bar and hides it
func (p *Progress) ShutdownBar() {
	if p.bar == nil {
		return
	}
	p.bar.Finish()
	p.bar = nil
	p.queue <- printTask{code: codeHideProgress}
}

// Write is implementation of io.Writer to support updating of progress bar
func (p *Progress) Write(s []byte) (int, error) {
	if p.bar != nil {
		p.bar.Add(len(s))
	}
	return len(s), nil
}

// AddBar increments progress for progress bar
func (p *Progress) AddBar(count int) {
	if p.bar != nil {
		p.bar.Add(count)
	}
}

// Printf does printf but in safe manner: not overwriting progress bar
func (p *Progress) Printf(msg string, a ...interface{}) {
	p.queue <- printTask{code: codePrint, message: fmt.Sprintf(msg, a...)}
}

func (p *Progress) worker() {
	for {
		select {
		case task := <-p.queue:
			switch task.code {
			case codePrint:
				if p.barShown {
					fmt.Print("\r\033[2K")
					p.barShown = false
				}
				fmt.Print(task.message)
			case codeProgress:
				fmt.Print("\r" + task.message)
				p.barShown = true
			case codeHideProgress:
				if p.barShown {
					fmt.Print("\r\033[2K")
					p.barShown = false
				}
			}
		case <-p.stop:
			p.stopped <- true
			return
		}
	}
}