package output

import (
	"fmt"
	"time"
)

type spinner struct {
	Prefix   string
	doneChan chan struct{}
}

func NewSpinner() spinner {
	return spinner{
		doneChan: make(chan struct{}, 1),
	}
}

func (s *spinner) Start() {
	go func() {
		for {
			for _, r := range `-\|/` {
				select {
				case <-s.doneChan:
					return
				default:
					fmt.Printf("\r%s %c", s.Prefix, r)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()
}

func (s *spinner) Stop() {
	s.erase()
	s.doneChan <- struct{}{}
}

func (s *spinner) erase() {
	// loop through output, backspacing
	l := len(s.Prefix) + 1
	for i := 0; i < l; i++ {
		for _, c := range []string{"\b", " ", "\b"} {
			fmt.Printf(c)
		}
	}
}
