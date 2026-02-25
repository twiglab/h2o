package hank

import (
	"net/http"
	"time"
)

type ElectricityItem struct {
	ElectricityMeter
	Time time.Time
}

type ElectricityPacket struct {
	Items map[string]*ElectricityItem
}

func NewElectricityPacket() *ElectricityPacket {
	return &ElectricityPacket{
		Items: make(map[string]*ElectricityItem),
	}
}

func (e *ElectricityPacket) Add(m ElectricityMeter) {
	e.Items[m.Code] = &ElectricityItem{ElectricityMeter: m, Time: time.Now()}
}

func (e *ElectricityPacket) List() (items []*ElectricityItem) {
	for _, item := range e.Items {
		items = append(items, item)
	}
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
