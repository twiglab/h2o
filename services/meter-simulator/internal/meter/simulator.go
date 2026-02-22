package meter

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// MeterState 电表状态
type MeterState struct {
	Address      byte    // Modbus 从站地址
	MeterNo      string  // 电表编号
	CurrentKWh   float64 // 当前电能读数 (kWh)
	MinIncrement float64 // 最小增量 (kWh)
	MaxIncrement float64 // 最大增量 (kWh)
}

// Simulator 电表模拟器
type Simulator struct {
	dtuID  string
	meters []*MeterState
	mu     sync.RWMutex
	rand   *rand.Rand
}

// NewSimulator 创建模拟器
func NewSimulator(dtuID string) *Simulator {
	return &Simulator{
		dtuID:  dtuID,
		meters: make([]*MeterState, 0),
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// AddMeter 添加电表
func (s *Simulator) AddMeter(address byte, meterNo string, initialKWh float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	meter := &MeterState{
		Address:      address,
		MeterNo:      meterNo,
		CurrentKWh:   initialKWh,
		MinIncrement: 0.1, // 最小增量 0.1 kWh
		MaxIncrement: 2.0, // 最大增量 2.0 kWh
	}
	s.meters = append(s.meters, meter)
}

// AddMeterWithRange 添加带自定义增量范围的电表
func (s *Simulator) AddMeterWithRange(address byte, meterNo string, initialKWh, minInc, maxInc float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	meter := &MeterState{
		Address:      address,
		MeterNo:      meterNo,
		CurrentKWh:   initialKWh,
		MinIncrement: minInc,
		MaxIncrement: maxInc,
	}
	s.meters = append(s.meters, meter)
}

// GetDTUID 获取 DTU ID
func (s *Simulator) GetDTUID() string {
	return s.dtuID
}

// GetMeters 获取所有电表
func (s *Simulator) GetMeters() []*MeterState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*MeterState, len(s.meters))
	copy(result, s.meters)
	return result
}

// UpdateAndGetReading 更新电表读数并返回
func (s *Simulator) UpdateAndGetReading(address byte) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, meter := range s.meters {
		if meter.Address == address {
			// 生成随机增量
			increment := meter.MinIncrement + s.rand.Float64()*(meter.MaxIncrement-meter.MinIncrement)
			meter.CurrentKWh += increment

			// 四舍五入到小数点后2位
			meter.CurrentKWh = float64(int(meter.CurrentKWh*100+0.5)) / 100

			return meter.CurrentKWh, nil
		}
	}

	return 0, fmt.Errorf("meter with address %d not found", address)
}

// GetReading 获取当前读数（不更新）
func (s *Simulator) GetReading(address byte) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, meter := range s.meters {
		if meter.Address == address {
			return meter.CurrentKWh, nil
		}
	}

	return 0, fmt.Errorf("meter with address %d not found", address)
}

// GetMeterByAddress 根据地址获取电表
func (s *Simulator) GetMeterByAddress(address byte) (*MeterState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, meter := range s.meters {
		if meter.Address == address {
			return meter, nil
		}
	}

	return nil, fmt.Errorf("meter with address %d not found", address)
}

// DefaultSimulator 创建默认的3电表模拟器
func DefaultSimulator(dtuID string) *Simulator {
	sim := NewSimulator(dtuID)

	// 添加3个电表，使用不同的初始读数和增量范围
	// 电表1: 商铺A，较低用电量
	sim.AddMeterWithRange(0x01, "METER001", 1000.00, 0.1, 1.0)

	// 电表2: 商铺B，中等用电量
	sim.AddMeterWithRange(0x02, "METER002", 2500.00, 0.5, 2.0)

	// 电表3: 商铺C，较高用电量
	sim.AddMeterWithRange(0x03, "METER003", 5000.00, 1.0, 3.0)

	return sim
}
