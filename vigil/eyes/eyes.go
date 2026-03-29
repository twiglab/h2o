package eyes

import (
	"encoding/json/v2"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/twiglab/h2o/pkg/common"
	"github.com/twiglab/h2o/vigil"
)

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

		ActivePowerTotal: float64(data.ActivePowerTotal) * 10,
	}
}

type Extra struct {
	StdCurrent float64 `json:"std_current"`
}

type ElectricityItem struct {
	vigil.Meter
	Data  EleMeterDataView `json:"data"`
	Extra Extra            `json:"extra"`
}

type ElectricityEgg struct {
	Items map[string]*ElectricityItem
}

func NewElectricityEgg() *ElectricityEgg {
	return &ElectricityEgg{
		Items: make(map[string]*ElectricityItem),
	}
}

func (e *ElectricityEgg) Merge(m vigil.ElectricityMeter) {
	data := toView(m.Data)
	i := &ElectricityItem{
		Meter: m.Meter,
		Data:  data,
		Extra: Extra{
			StdCurrent: std(data.CurrentA, data.CurrentB, data.CurrentC),
		},
	}
	e.Items[i.Code] = i
}

func (e *ElectricityEgg) List() (items []*ElectricityItem) {
	for _, item := range e.Items {
		items = append(items, item)
	}
	return
}

func (e *ElectricityEgg) Get(code string) *ElectricityItem {
	if ei, ok := e.Items[code]; ok {
		return ei
	}
	return nil
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
