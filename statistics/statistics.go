package statistics

import (
	"fmt"
	"os"
	"sync"
	"text/tabwriter"
)

type Statistic struct {
	Duration float64
	Success  bool
}

type Statistics struct {
	sl              []Statistic
	totalDuration   float64
	ProcessDuration float64
	m               sync.Mutex
}

func (s *Statistics) AddStatistic(st Statistic) {
	s.m.Lock()
	s.sl = append(s.sl, st)
	s.m.Unlock()
}

func (s *Statistics) Report() {
	mean := 0.0
	unsuccess := 0
	success := 0

	for _, st := range s.sl {
		if st.Success {
			success += 1
		} else {
			unsuccess += 1
		}

		mean += st.Duration
	}

	s.totalDuration = mean
	mean = mean / float64(len(s.sl))

	fmt.Printf("===============================================================\n")
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(w, "Mean Duration:\t%fs\nTotal Duration:\t%fs\nProcess Duration:\t%fs\nSuccessful:\t%d\nUnsuccessful:\t%d\nTotal:\t%d\n", mean, s.totalDuration, s.ProcessDuration, success, unsuccess, len(s.sl))
	w.Flush()
	fmt.Printf("===============================================================\n")
}
