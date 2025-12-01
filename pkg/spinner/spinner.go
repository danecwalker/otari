package spinner

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type SpinnerID int

const (
	LineSpinnerID SpinnerID = iota
	DotsSpinnerID
)

type SpinnerOptionFunc func(s *Spinner)

func WithSuccessSymbol(symbol string) SpinnerOptionFunc {
	return func(s *Spinner) {
		s.successSymbol = symbol
	}
}

func WithErrorSymbol(symbol string) SpinnerOptionFunc {
	return func(s *Spinner) {
		s.errorSymbol = symbol
	}
}

func WithInfoSymbol(symbol string) SpinnerOptionFunc {
	return func(s *Spinner) {
		s.infoSymbol = symbol
	}
}

type Spinner struct {
	message string
	frames  []string

	successSymbol string
	errorSymbol   string
	infoSymbol    string

	ticker       *time.Ticker
	done         chan struct{}
	out          io.Writer
	mu           sync.Mutex
	printlnCount int
	frameIndex   int
}

// ANSI sequences
const (
	hideCursor = "\x1b[?25l"
	showCursor = "\x1b[?25h"
	clearLine  = "\x1b[2K"
	cursorUp   = "\x1b[%dA"
)

func NewCustom(frames []string, opts ...SpinnerOptionFunc) *Spinner {
	s := &Spinner{
		frames: frames,
		out:    os.Stdout,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func New(spinner SpinnerID, opts ...SpinnerOptionFunc) *Spinner {
	var s *Spinner
	switch spinner {
	case LineSpinnerID:
		s = &Spinner{
			frames: []string{"-", "\\", "|", "/"},
			out:    os.Stdout,
		}
	case DotsSpinnerID:
		s = &Spinner{
			frames: []string{".    ", "..   ", "...  ", ".... ", ".....", " ....", "  ...", "   ..", "    .", "     "},
			out:    os.Stdout,
		}
	default:
		s = &Spinner{
			frames: []string{"-", "\\", "|", "/"},
			out:    os.Stdout,
		}
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Spinner) SetMessage(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = msg
}

func (s *Spinner) Println(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Clear spinner line before printing
	fmt.Fprintf(s.out, "\r%s\r", clearLine)

	s.printlnCount++

	// Print real line
	fmt.Fprintln(s.out, msg)

	// Redraw spinner after printed line
	if s.ticker != nil {
		frame := s.frames[s.frameIndex%len(s.frames)]
		fmt.Fprintf(s.out, "%s %s", frame, s.message)
	}
}

func (s *Spinner) Enable(interval time.Duration) {
	s.mu.Lock()

	if s.out == nil {
		s.out = os.Stdout
	}
	s.ticker = time.NewTicker(interval)
	s.done = make(chan struct{})

	// hide cursor
	fmt.Fprint(s.out, hideCursor)

	s.mu.Unlock()

	go func() {
		for {
			if s.ticker == nil {
				return
			}
			select {
			case <-s.ticker.C:
				s.mu.Lock()
				s.frameIndex = (s.frameIndex + 1) % len(s.frames)
				frame := s.frames[s.frameIndex]
				fmt.Fprintf(s.out, "\r%s %s", frame, s.message)
				s.mu.Unlock()
			case <-s.done:
				return
			}
		}
	}()
}

func (s *Spinner) clearOutputBlock() {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Move up printed lines + spinner line
	total := s.printlnCount + 1 // spinner line + println lines
	if total <= 0 {
		fmt.Fprint(s.out, showCursor)
		return
	}

	// Assume we're currently on the spinner line
	fmt.Fprint(s.out, "\r")

	for i := range total {
		// Clear current line
		fmt.Fprint(s.out, clearLine)
		if i < total-1 {
			// Move cursor up one line and to column 0,
			// but don't go above the top of the spinner block
			fmt.Fprintf(s.out, cursorUp, 1)
		}
	}

	// Now we're at the topmost cleared line, column 0
}

func (s *Spinner) stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = nil
	}
	if s.done != nil {
		close(s.done)
		s.done = nil
	}
}

func (s *Spinner) FinishWithMessage(msg string) {
	s.stop()
	s.clearOutputBlock()
	fmt.Fprint(s.out, showCursor)
	fmt.Fprintf(s.out, "%s\n", msg)
}

func (s *Spinner) FinishWithSuccess(msg string) {
	if s.successSymbol != "" {
		msg = fmt.Sprintf("%s %s", s.successSymbol, msg)
	}
	s.FinishWithMessage(msg)
}
func (s *Spinner) FinishWithError(msg string) {
	if s.errorSymbol != "" {
		msg = fmt.Sprintf("%s %s", s.errorSymbol, msg)
	}
	s.FinishWithMessage(msg)
}

func (s *Spinner) FinishWithInfo(msg string) {
	if s.infoSymbol != "" {
		msg = fmt.Sprintf("%s %s", s.infoSymbol, msg)
	}
	s.FinishWithMessage(msg)
}

func (s *Spinner) Finish() { s.FinishWithMessage("") }
