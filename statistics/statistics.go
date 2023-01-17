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
	Success         int
	Unsuccess       int
	RemoteIp        string
	ProcessDuration float64
	totalStatistic  map[string]float64
	avgStatistics   map[string]float64
	minStatistics   map[string]float64
	maxStatistics   map[string]float64
	slMutex         sync.RWMutex
	totalMutex      sync.RWMutex
	avgMutex        sync.RWMutex
	minMutex        sync.RWMutex
	maxMutex        sync.RWMutex
}

func (s *Statistics) AddStatistic(st map[string]float64) {
	s.slMutex.Lock()
	s.sl = append(s.sl, st)
	s.slMutex.Unlock()
}

func (s *Statistics) Report() {
	fmt.Println("==================================================================")
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(w, "SUCCESS:\t%d\nFAILED:\t%d\n", s.Success, s.Unsuccess)
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
	s.totalStatistic = NewStatistic()
	for _, st := range s.sl {
		for key, value := range st {
			s.totalMutex.Lock()
			s.totalStatistic[key] += value
			s.totalMutex.Unlock()
		}
	}
}

func (s *Statistics) calcavgStatistics() {
	s.avgStatistics = NewStatistic()
	for key, value := range s.totalStatistic {
		s.avgMutex.Lock()
		s.avgStatistics[key] = value / float64(len(s.sl))
		s.avgMutex.Unlock()
	}
}

func (s *Statistics) calcMinmaxStatistics() {
	min := math.MaxFloat64
	max := 0.0
	s.maxStatistics = NewStatistic()
	s.minStatistics = NewStatistic()

	for _, st := range s.sl {
		for key, value := range st {
			if value > max {
				s.maxMutex.Lock()
				s.maxStatistics[key] = value
				s.maxMutex.Unlock()
			}

			if value < min {
				s.minMutex.Lock()
				s.minStatistics[key] = value
				s.minMutex.Unlock()
			}
		}
	}
}

func NewStatistic() map[string]float64 {
	return map[string]float64{
		"DIAL":  0,
		"TOUCH": 0,
		"HELO":  0,
		"MAIL":  0,
		"RCPT":  0,
		"DATA":  0,
		"QUIT":  0,
	}
}

func NewStatistics(initial float64) *Statistics {
	return &Statistics{
		sl:             make([]map[string]float64, 0),
		totalStatistic: NewStatistic(),
		avgStatistics:  NewStatistic(),
		minStatistics:  NewStatistic(),
		maxStatistics:  NewStatistic(),
	}
}
