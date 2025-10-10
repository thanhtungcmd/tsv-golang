package datetime

import (
	"time"
)

type datetime struct {
	time time.Time
}

func Datetime() *datetime {
	return &datetime{}
}

func (h *datetime) StartOfDay() int64 {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Unix()
}

func (h *datetime) EndOfDay() int64 {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 59, time.Local).Unix()
}

func (h *datetime) TimeNow() *datetime {
	h.time = time.Now()
	return h
}

func (h *datetime) CreateFromInt(input uint) *datetime {
	h.time = time.Unix(int64(input), 0)
	return h
}

func (h *datetime) ToString() string {
	formattedTime := h.time.Format("2006-01-02 15:04:05")
	return formattedTime
}

func (h *datetime) AddDays(i int) *datetime {
	h.time = h.time.AddDate(0, 0, i)
	return h
}

func (h *datetime) ToInt() int64 {
	return h.time.Unix()
}
