package hank

import (
	"cmp"
	"net/http"
	"slices"
	"time"

	"github.com/twiglab/h2o/pkg/common"
)

type window struct {
	candidate []common.Electricity
	pos       int
	l         int
}

func newWin(l int) *window {
	return &window{
		candidate: make([]common.Electricity, l, l),
		l:         l,
	}
}

func (w *window) Add(ele common.Electricity) {
	np := w.pos % w.l
	w.candidate[np] = ele
	w.pos = np + 1
}

type ElectricityItem struct {
	ElectricityMeter
	AcceptTime time.Time     `json:"accept_time"`
	Delay      time.Duration `json:"delay,format:sec"`
	TimeOut    time.Duration `json:"timeout,format:sec"`

	window *window `json:"-"`
}

func (ei *ElectricityItem) fill() {
	for _, ele := range ei.window.candidate {
		ei.Data.DataValue = max(ei.Data.DataValue, ele.DataValue)

		ei.Data.VoltageA = max(ei.Data.VoltageA, ele.VoltageA)
		ei.Data.VoltageB = max(ei.Data.VoltageB, ele.VoltageB)
		ei.Data.VoltageC = max(ei.Data.VoltageC, ele.VoltageC)

		ei.Data.CurrentA = max(ei.Data.CurrentA, ele.CurrentA)
		ei.Data.CurrentB = max(ei.Data.CurrentB, ele.CurrentB)
		ei.Data.CurrentC = max(ei.Data.CurrentC, ele.CurrentC)

		ei.Data.TotalActivePower = max(ei.Data.TotalActivePower, ele.TotalActivePower)
	}
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
		i = &ElectricityItem{ElectricityMeter: m, AcceptTime: time.Now(), window: newWin(10)}
		i.window.Add(m.Data)
		e.Items[i.Code] = i
		return
	}

	i.window.Add(m.Data)

	i.TimeOut = 0

	i.Device = m.Device
	i.Pos = m.Pos

	i.AcceptTime = time.Now()
	i.Delay = time.Since(m.DataTime)
}

func (e *ElectricityPacket) List() (items []*ElectricityItem) {
	for _, item := range e.Items {
		item.fill()
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
