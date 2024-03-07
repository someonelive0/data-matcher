package matcher

import (
	"encoding/json"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type MyStatistic struct {
	Starttime      time.Time `json:"start_time"`
	Lastinputtime  time.Time `json:"last_input_time"`
	Lastoutputtime time.Time `json:"last_output_time"`
	CountInput     uint64    `json:"count_input"`
	CountOutput    uint64    `json:"count_output"`
	// Matchedpolicy  MatchedPolicy     `json:"matchd_policy"`
	Dailystatistic []*DailyStatistic `json:"daily_statistic"`
	dailymap       map[string]*DailyStatistic
	mutex          sync.Mutex
}

func NewMyStatistic(t time.Time) *MyStatistic {
	return &MyStatistic{
		Starttime:      t,
		Dailystatistic: make([]*DailyStatistic, 0),
		dailymap:       make(map[string]*DailyStatistic),
	}
}

type DailyStatistic struct {
	Day   string `json:"day"`
	Count uint64 `json:"count"`
}

func (p *MyStatistic) InputCount(n uint64) {
	atomic.AddUint64(&p.CountInput, n)
	p.Lastinputtime = time.Now()
	day := p.Lastinputtime.Format(time.DateOnly)

	p.mutex.Lock()
	if i, ok := p.dailymap[day]; !ok {
		dailysta := &DailyStatistic{Day: day, Count: n}
		p.dailymap[day] = dailysta
		p.Dailystatistic = append(p.Dailystatistic, dailysta)
	} else {
		i.Count++
	}
	p.mutex.Unlock()
}

// mutex is not used
func (p *MyStatistic) OutputCount(n uint64) {
	atomic.AddUint64(&p.CountOutput, n)
	p.Lastoutputtime = time.Now()
}

func (p *MyStatistic) Dump() []byte {
	p.mutex.Lock()
	sort.Slice(p.Dailystatistic, func(i, j int) bool {
		return p.Dailystatistic[i].Day > p.Dailystatistic[j].Day
	})
	jsonb, _ := json.Marshal(p)
	p.mutex.Unlock()
	return jsonb
}
