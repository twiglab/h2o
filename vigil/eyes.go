package vigil

import (
	"cmp"
	"encoding/json/v2"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/twiglab/h2o/pkg/common"
)

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

func toView(data common.Electricity) EleMeterDataView {
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
}

type ElectricityEgg struct {
	Items map[string]*ElectricityItem
}

func NewElectricityPacket() *ElectricityEgg {
	return &ElectricityEgg{
		Items: make(map[string]*ElectricityItem),
	}
}

func (e *ElectricityEgg) Merge(m ElectricityMeter) {

	md := toView(m.Data)
	i := &ElectricityItem{
		EleMeterView: EleMeterView{
			Meter: m.Meter,
			Data:  md,
		},
	}

	e.Items[i.Code] = i
}

func (e *ElectricityEgg) List() (items []*ElectricityItem) {
	for _, item := range e.Items {
		items = append(items, item)
	}
	slices.SortFunc(items, func(a, b *ElectricityItem) int { return cmp.Compare(a.Code, b.Code) })
	return
}

func (e *ElectricityEgg) Get(code string) *ElectricityItem {
	if ei, ok := e.Items[code]; ok {
		return ei
	}
	return nil
}

func (e *ElectricityEgg) SetStatus(code string, status int) {
	if v, ok := e.Items[code]; ok {
		v.Status = status
	}
}

func EyesAll(ep *ElectricityEgg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.MarshalWrite(w, ep.List()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func EyesMux(ep *ElectricityEgg) http.Handler {
	mux := chi.NewMux()
	mux.HandleFunc("/all", EyesAll(ep))
	return mux
}
