package statistics

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
)

type Statistics struct {
	sl              []map[string]float64
	success         int
	unsuccess       int
	RemoteIp        string
	ProcessDuration float64
	totalStatistic  map[string]float64
	avgStatistics   map[string]float64
	minStatistics   map[string]float64
	maxStatistics   map[string]float64
	slMutex         sync.RWMutex
	successMutex    sync.RWMutex
}

func (s *Statistics) AddStatistic(st map[string]float64) {
	s.slMutex.Lock()
	s.sl = append(s.sl, st)
	s.slMutex.Unlock()
}

func (s *Statistics) AddSuccess(success bool) {
	s.successMutex.Lock()
	if success {
		s.success += 1
	} else {
		s.unsuccess += 1
	}
	s.successMutex.Unlock()
}

func (s *Statistics) Report() {
	fmt.Println("==================================================================")
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(w, "SUCCESS:\t%d\nFAILED:\t%d\n", s.success, s.unsuccess)
	w.Flush()
	fmt.Println("==================================================================")

	s.sumStatistics()
	s.calcavgStatistics()
	s.calcMinmaxStatistics()

	w = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "\tMIN\tAVG\tMAX\tTOTAL")
	for key, value := range s.totalStatistic {
		fmt.Fprintf(w, "%s:\t%fs\t%fs\t%fs\t%fs\n", key, s.minStatistics[key], s.avgStatistics[key], s.maxStatistics[key], value)
	}

	w.Flush()
	fmt.Print("\n\n")
	fmt.Printf("Total Process Duration: %fs\n", s.ProcessDuration)
	fmt.Println("==================================================================")
}

func (s *Statistics) sumStatistics() {
	for _, st := range s.sl {
		for key, value := range st {
			s.totalStatistic[key] += value
		}
	}
}

func (s *Statistics) calcavgStatistics() {
	for key, value := range s.totalStatistic {
		s.avgStatistics[key] = value / float64(len(s.sl))
	}
}

func (s *Statistics) calcMinmaxStatistics() {
	for _, st := range s.sl {
		for key, value := range st {
			if value > s.maxStatistics[key] {
				s.maxStatistics[key] = value
			}

			if value < s.minStatistics[key] {
				s.minStatistics[key] = value
			}
		}
	}
}

func NewStatistic(initial float64) map[string]float64 {
	return map[string]float64{
		"DIAL":  initial,
		"TOUCH": initial,
		"HELO":  initial,
		"MAIL":  initial,
		"RCPT":  initial,
		"DATA":  initial,
		"QUIT":  initial,
	}
}

func NewStatistics() *Statistics {
	return &Statistics{
		sl:             make([]map[string]float64, 0),
		totalStatistic: NewStatistic(0),
		avgStatistics:  NewStatistic(0),
		minStatistics:  NewStatistic(math.MaxFloat64),
		maxStatistics:  NewStatistic(0),
	}
}
