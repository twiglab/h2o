package hank

import (
	"cmp"
	"net/http"
	"slices"
	"time"
)

type ElectricityItem struct {
	ElectricityMeter
	AcctpeTime time.Time
	Delay      time.Duration `json:"delay,format:sec"`
	TimeOut    time.Duration `json:"timeout,format:sec"`
}

type ElectricityPacket struct {
	Items map[string]*ElectricityItem
}

func NewElectricityPacket() *ElectricityPacket {
	return &ElectricityPacket{
		Items: make(map[string]*ElectricityItem),
	}
}

func (e *ElectricityPacket) Merge(m ElectricityMeter) {
	var (
		i  *ElectricityItem
		ok bool
	)

	if i, ok = e.Items[m.Code]; !ok {
		i = &ElectricityItem{ElectricityMeter: m, AcctpeTime: time.Now()}
		e.Items[i.Code] = i
		return
	}

	i.TimeOut = 0

	i.Device = m.Device
	i.Pos = m.Pos
	i.Data = m.Data

	i.AcctpeTime = time.Now()
	i.Delay = time.Since(m.DataTime)

	if m.Data.DataValue > 0 {
		i.Data.DataValue = m.Data.DataValue
	}

}

func (e *ElectricityPacket) List() (items []*ElectricityItem) {
	for _, item := range e.Items {
		item.TimeOut = time.Since(item.DataTime)
		items = append(items, item)
	}
	slices.SortFunc(items, func(a, b *ElectricityItem) int { return cmp.Compare(a.Code, b.Code) })
	return
}

func (e *ElectricityPacket) SetStatus(code string, status int) {
	if v, ok := e.Items[code]; ok {
		v.Status = status
	}
}

func EyesAll(ep *ElectricityPacket) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := marshalWrite(w, ep.List()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
