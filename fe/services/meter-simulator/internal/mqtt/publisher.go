package mqtt

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// PublisherConfig 发布者配置
type PublisherConfig struct {
	Broker   string
	ClientID string
	Username string
	Password string
}

// Publisher MQTT 发布者
type Publisher struct {
	config PublisherConfig
	client mqtt.Client
}

// NewPublisher 创建发布者
func NewPublisher(config PublisherConfig) *Publisher {
	return &Publisher{
		config: config,
	}
}

// Connect 连接到 MQTT 服务器
func (p *Publisher) Connect() error {
	opts := mqtt.NewClientOptions().
		AddBroker(p.config.Broker).
		SetClientID(p.config.ClientID).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetKeepAlive(60 * time.Second).
		SetCleanSession(true)

	if p.config.Username != "" {
		opts.SetUsername(p.config.Username)
	}
	if p.config.Password != "" {
		opts.SetPassword(p.config.Password)
	}

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Printf("[MQTT] Connected to broker: %s", p.config.Broker)
	})

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("[MQTT] Connection lost: %v", err)
	})

	opts.SetReconnectingHandler(func(client mqtt.Client, opts *mqtt.ClientOptions) {
		log.Printf("[MQTT] Reconnecting...")
	})

	p.client = mqtt.NewClient(opts)

	token := p.client.Connect()
	if !token.WaitTimeout(30 * time.Second) {
		return fmt.Errorf("connect timeout")
	}

	if err := token.Error(); err != nil {
		return fmt.Errorf("connect error: %w", err)
	}

	return nil
}

// Publish 发布消息
func (p *Publisher) Publish(topic string, qos byte, payload []byte) error {
	token := p.client.Publish(topic, qos, false, payload)
	if !token.WaitTimeout(10 * time.Second) {
		return fmt.Errorf("publish timeout")
	}

	if err := token.Error(); err != nil {
		return fmt.Errorf("publish error: %w", err)
	}

	return nil
}

// Disconnect 断开连接
func (p *Publisher) Disconnect() {
	if p.client != nil && p.client.IsConnected() {
		p.client.Disconnect(1000)
		log.Printf("[MQTT] Disconnected")
	}
}

// IsConnected 是否已连接
func (p *Publisher) IsConnected() bool {
	return p.client != nil && p.client.IsConnected()
}
