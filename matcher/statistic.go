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
	CountSlowInput uint64    `json:"count_slow_input"`
	CountOutHttp   uint64    `json:"count_output_http"`
	CountOutDns    uint64    `json:"count_output_dns"`
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

func (p *MyStatistic) InputSlowCount(n uint64) {
	atomic.AddUint64(&p.CountSlowInput, n)
}

// mutex is not used
func (p *MyStatistic) OutputHttpCount(n uint64) {
	atomic.AddUint64(&p.CountOutHttp, n)
	p.Lastoutputtime = time.Now()
}

func (p *MyStatistic) OutputDnsCount(n uint64) {
	atomic.AddUint64(&p.CountOutDns, n)
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
