package utils

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
)

// 默认错误日志数量为200
const ERROR_SIZE = 200

// 定义错误码缓存
type ErrorLog struct {
	Timestamp time.Time `json:"timestamp"`
	Code      int       `json:"code"`
	Message   string    `json:"message"`
}

type MyErrors struct {
	r *Ring
}

func NewMyErrors(capacity int) *MyErrors {
	r := &Ring{}
	r.SetCapacity(capacity)
	if r.Capacity() != capacity {
		return nil
	}
	return &MyErrors{r: r}
}

func (p *MyErrors) Push(errlog *ErrorLog) {
	p.r.Enqueue(errlog)
}

func (p *MyErrors) Dump() []byte {
	errlogs := p.r.Values()
	for i, j := 0, len(errlogs)-1; i < j; i, j = i+1, j-1 {
		errlogs[i], errlogs[j] = errlogs[j], errlogs[i] // reverse the slice
	}
	b, _ := json.MarshalIndent(errlogs, " ", "")
	return b

}

// for logrus hook
type MyErrorHook struct {
	Myerrors *MyErrors
}

func NewMyErrorHook() *MyErrorHook {
	return &MyErrorHook{Myerrors: NewMyErrors(ERROR_SIZE)}
}

func (p *MyErrorHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (p *MyErrorHook) Fire(entry *logrus.Entry) error {
	if entry.Level == logrus.ErrorLevel {
		errlog := &ErrorLog{
			Timestamp: time.Now(),
			Code:      500,
			Message:   entry.Message,
		}
		p.Myerrors.Push(errlog)
	}
	return nil
}

func AddLogHook() *MyErrors {
	var hook = NewMyErrorHook()
	logrus.AddHook(hook)
	return hook.Myerrors
}
