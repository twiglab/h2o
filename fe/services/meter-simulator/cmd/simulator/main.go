package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"meter-simulator/internal/meter"
	"meter-simulator/internal/modbus"
	"meter-simulator/internal/mqtt"
)

var (
	mqttBroker   = flag.String("broker", "tcp://localhost:1883", "MQTT broker address")
	mqttClientID = flag.String("client-id", "meter-simulator", "MQTT client ID")
	mqttUsername = flag.String("username", "", "MQTT username")
	mqttPassword = flag.String("password", "", "MQTT password")
	dtuID        = flag.String("dtu-id", "DTU001", "DTU device ID")
	interval     = flag.Duration("interval", 30*time.Second, "Report interval")
)

func main() {
	flag.Parse()

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Printf("Starting meter simulator...")
	log.Printf("  DTU ID: %s", *dtuID)
	log.Printf("  MQTT Broker: %s", *mqttBroker)
	log.Printf("  Report Interval: %s", *interval)

	// 创建 MQTT 发布者
	publisher := mqtt.NewPublisher(mqtt.PublisherConfig{
		Broker:   *mqttBroker,
		ClientID: *mqttClientID,
		Username: *mqttUsername,
		Password: *mqttPassword,
	})

	// 连接 MQTT
	if err := publisher.Connect(); err != nil {
		log.Fatalf("Failed to connect MQTT: %v", err)
	}
	defer publisher.Disconnect()

	// 创建电表模拟器（3个电表）
	sim := meter.DefaultSimulator(*dtuID)

	// 打印电表信息
	log.Printf("Simulating %d meters:", len(sim.GetMeters()))
	for _, m := range sim.GetMeters() {
		log.Printf("  - Address: 0x%02X, MeterNo: %s, Initial: %.2f kWh",
			m.Address, m.MeterNo, m.CurrentKWh)
	}

	// 创建退出信号通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 创建定时器
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	// 立即发送一次
	reportAllMeters(publisher, sim)

	// 主循环
	log.Printf("Simulator running, press Ctrl+C to stop...")
	for {
		select {
		case <-ticker.C:
			reportAllMeters(publisher, sim)
		case <-quit:
			log.Printf("Shutting down...")
			return
		}
	}
}

// reportAllMeters 上报所有电表数据
func reportAllMeters(publisher *mqtt.Publisher, sim *meter.Simulator) {
	meters := sim.GetMeters()
	topic := fmt.Sprintf("dtu/%s/data", sim.GetDTUID())

	log.Printf("--- Reporting %d meters ---", len(meters))

	for _, m := range meters {
		// 更新读数（模拟用电增长）
		reading, err := sim.UpdateAndGetReading(m.Address)
		if err != nil {
			log.Printf("[ERROR] Failed to update meter %s: %v", m.MeterNo, err)
			continue
		}

		// 构建 Modbus RTU 响应帧
		frame := modbus.BuildEnergyReadingResponse(m.Address, reading)

		// 发布到 MQTT
		if err := publisher.Publish(topic, 1, frame); err != nil {
			log.Printf("[ERROR] Failed to publish meter %s: %v", m.MeterNo, err)
			continue
		}

		log.Printf("[OK] Meter %s (0x%02X): %.2f kWh -> %s",
			m.MeterNo, m.Address, reading, hex.EncodeToString(frame))

		// 短暂延迟，模拟真实轮询
		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("--- Report complete ---")
}
