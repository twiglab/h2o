package mqtt

import (
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

// ClientConfig MQTT客户端配置
type ClientConfig struct {
	Broker   string
	ClientID string
	Username string
	Password string
}

// Client MQTT客户端
type Client struct {
	config  ClientConfig
	client  mqtt.Client
	logger  *zap.Logger
	handler MessageHandler
	mu      sync.RWMutex
}

// MessageHandler 消息处理器接口
type MessageHandler interface {
	HandleMessage(topic string, payload []byte)
}

// NewClient 创建MQTT客户端
func NewClient(config ClientConfig, logger *zap.Logger) *Client {
	return &Client{
		config: config,
		logger: logger,
	}
}

// SetHandler 设置消息处理器
func (c *Client) SetHandler(handler MessageHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handler = handler
}

// Connect 连接到MQTT服务器
func (c *Client) Connect() error {
	opts := mqtt.NewClientOptions().
		AddBroker(c.config.Broker).
		SetClientID(c.config.ClientID).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetKeepAlive(60 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetCleanSession(false).
		SetOrderMatters(false)

	if c.config.Username != "" {
		opts.SetUsername(c.config.Username)
	}
	if c.config.Password != "" {
		opts.SetPassword(c.config.Password)
	}

	// 连接回调
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		c.logger.Info("connected to MQTT broker",
			zap.String("broker", c.config.Broker))
	})

	// 断开回调
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		c.logger.Warn("MQTT connection lost",
			zap.Error(err))
	})

	// 重连回调
	opts.SetReconnectingHandler(func(client mqtt.Client, opts *mqtt.ClientOptions) {
		c.logger.Info("reconnecting to MQTT broker...")
	})

	c.client = mqtt.NewClient(opts)

	token := c.client.Connect()
	if !token.WaitTimeout(30 * time.Second) {
		return fmt.Errorf("MQTT connect timeout")
	}

	if err := token.Error(); err != nil {
		return fmt.Errorf("MQTT connect error: %w", err)
	}

	return nil
}

// Subscribe 订阅主题
func (c *Client) Subscribe(topic string, qos byte) error {
	token := c.client.Subscribe(topic, qos, c.messageCallback)
	if !token.WaitTimeout(10 * time.Second) {
		return fmt.Errorf("MQTT subscribe timeout")
	}

	if err := token.Error(); err != nil {
		return fmt.Errorf("MQTT subscribe error: %w", err)
	}

	c.logger.Info("subscribed to topic",
		zap.String("topic", topic),
		zap.Uint8("qos", qos))

	return nil
}

// messageCallback MQTT消息回调
func (c *Client) messageCallback(client mqtt.Client, msg mqtt.Message) {
	c.mu.RLock()
	handler := c.handler
	c.mu.RUnlock()

	if handler != nil {
		handler.HandleMessage(msg.Topic(), msg.Payload())
	}
}

// Publish 发布消息
func (c *Client) Publish(topic string, qos byte, retained bool, payload []byte) error {
	token := c.client.Publish(topic, qos, retained, payload)
	if !token.WaitTimeout(10 * time.Second) {
		return fmt.Errorf("MQTT publish timeout")
	}

	if err := token.Error(); err != nil {
		return fmt.Errorf("MQTT publish error: %w", err)
	}

	return nil
}

// Disconnect 断开连接
func (c *Client) Disconnect() {
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(1000)
		c.logger.Info("disconnected from MQTT broker")
	}
}

// IsConnected 是否已连接
func (c *Client) IsConnected() bool {
	return c.client != nil && c.client.IsConnected()
}
