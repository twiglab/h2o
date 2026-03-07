package hank

import (
	"cmp"
	"net/http"
	"slices"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/twiglab/h2o/pkg/common"
)

const windowSize = 4

type window struct {
	Candidate []EleMeterDataView
	c         int
	l         int
}

func newWin(l int) *window {
	return &window{
		Candidate: make([]EleMeterDataView, l),
		l:         l,
	}
}

func (w *window) Add(ele EleMeterDataView) {
	np := w.c % w.l
	w.Candidate[np] = ele
	w.c = np + 1
}

func (w *window) Clear() {
	clear(w.Candidate)
}

type EleMeterView struct {
	Meter
	Data EleMeterDataView `json:"data"`
}

type EleMeterDataView struct {
	DataValue float64 `json:"data_value"`

	VoltageA float64 `json:"voltage_a,omitempty"`
	VoltageB float64 `json:"voltage_b,omitempty"`
	VoltageC float64 `json:"voltage_c,omitempty"`

	CurrentA float64 `json:"current_a,omitempty"`
	CurrentB float64 `json:"current_b,omitempty"`
	CurrentC float64 `json:"current_c,omitempty"`

	ActivePowerTotal   float64 `json:"active_power_total,omitempty"`   // 总有功功率  P
	ReactivePowerTotal float64 `json:"reactive_power_total,omitempty"` // 总无功功率  Q
	ApparentPowerTotal float64 `json:"apparent_power_total,omitempty"` // 总视在功率  S
}

func toEleMeterData(data common.Electricity) EleMeterDataView {
	return EleMeterDataView{
		DataValue: float64(data.DataValue) / 100,

		VoltageA: float64(data.VoltageA) / 100,
		VoltageB: float64(data.VoltageB) / 100,
		VoltageC: float64(data.VoltageC) / 100,

		CurrentA: float64(data.CurrentA) / 100,
		CurrentB: float64(data.CurrentB) / 100,
		CurrentC: float64(data.CurrentC) / 100,

		ActivePowerTotal: float64(data.ActivePowerTotal) / 100,
	}
}

type ElectricityItem struct {
	EleMeterView
	AcceptTime time.Time     `json:"accept_time"`
	Delay      time.Duration `json:"delay,format:sec"`
	TimeOut    time.Duration `json:"timeout,format:sec"`

	window *window `json:"-"`

	L int
	C int
}

func (ei *ElectricityItem) fill() {
	ei.Data = EleMeterDataView{}
	for _, ele := range ei.window.Candidate {
		ei.Data.DataValue = max(ei.Data.DataValue, ele.DataValue)

		ei.Data.VoltageA = max(ei.Data.VoltageA, ele.VoltageA)
		ei.Data.VoltageB = max(ei.Data.VoltageB, ele.VoltageB)
		ei.Data.VoltageC = max(ei.Data.VoltageC, ele.VoltageC)

		ei.Data.CurrentA = max(ei.Data.CurrentA, ele.CurrentA)
		ei.Data.CurrentB = max(ei.Data.CurrentB, ele.CurrentB)
		ei.Data.CurrentC = max(ei.Data.CurrentC, ele.CurrentC)

		ei.Data.ActivePowerTotal = max(ei.Data.ActivePowerTotal, ele.ActivePowerTotal)
	}

	ei.L = ei.window.l
	ei.C = ei.window.c
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
		md := toEleMeterData(m.Data)
		i = &ElectricityItem{
			EleMeterView: EleMeterView{
				Meter: m.Meter, Data: md,
			},
			AcceptTime: time.Now(),
			window:     newWin(windowSize),
		}

		i.window.Add(md)
		e.Items[i.Code] = i
		return
	}

	i.window.Add(i.Data)

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

func (e *ElectricityPacket) Get(code string) *ElectricityItem {
	if ei, ok := e.Items[code]; ok {
		ei.fill()
		return ei
	}
	return nil

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

func EyesMux(ep *ElectricityPacket) http.Handler {
	mux := chi.NewMux()
	mux.HandleFunc("/all", EyesAll(ep))
	return mux
}
