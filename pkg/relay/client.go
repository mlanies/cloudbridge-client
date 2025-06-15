package relay

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// Message types
const (
	MessageTypeRegister = "register"
	MessageTypeConnect  = "connect"
)

// Error types
const (
	ErrInvalidToken          = "invalid_token"
	ErrRateLimitExceeded     = "rate_limit_exceeded"
	ErrConnectionLimitReached = "connection_limit_reached"
	ErrServerUnavailable     = "server_unavailable"
)

// RegisterMessage represents a registration message
type RegisterMessage struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

// ConnectMessage represents a connection message
type ConnectMessage struct {
	Type      string `json:"type"`
	LocalPort int    `json:"local_port"`
}

// Response represents a server response
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// Client represents a CloudBridge Relay client
type Client struct {
	conn     net.Conn
	encoder  *json.Encoder
	decoder  *json.Decoder
	useTLS   bool
	config   *tls.Config
	reconnectDelay time.Duration
	maxRetries     int
}

// NewClient creates a new CloudBridge Relay client
func NewClient(useTLS bool, config *tls.Config) *Client {
	return &Client{
		useTLS:   useTLS,
		config:   config,
		reconnectDelay: time.Second * 5,
		maxRetries:     3,
	}
}

// Connect establishes a connection to the relay server
func (c *Client) Connect(host string, port int) error {
	var err error
	var conn net.Conn

	if c.useTLS {
		conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), c.config)
	} else {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	}

	if err != nil {
		return fmt.Errorf("failed to connect to relay: %w", err)
	}

	c.conn = conn
	c.encoder = json.NewEncoder(conn)
	c.decoder = json.NewDecoder(conn)
	return nil
}

// Register registers the client with the relay server using a JWT token
func (c *Client) Register(token string) error {
	msg := RegisterMessage{
		Type:  MessageTypeRegister,
		Token: token,
	}

	if err := c.encoder.Encode(msg); err != nil {
		return fmt.Errorf("failed to send register message: %w", err)
	}

	var response Response
	if err := c.decoder.Decode(&response); err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if response.Status != "ok" {
		return fmt.Errorf("registration failed: %s", response.Error)
	}

	return nil
}

// CreateTunnel creates a tunnel to the specified local port
func (c *Client) CreateTunnel(localPort int) error {
	msg := ConnectMessage{
		Type:      MessageTypeConnect,
		LocalPort: localPort,
	}

	if err := c.encoder.Encode(msg); err != nil {
		return fmt.Errorf("failed to send connect message: %w", err)
	}

	var response Response
	if err := c.decoder.Decode(&response); err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if response.Status != "ok" {
		return fmt.Errorf("tunnel creation failed: %s", response.Error)
	}

	return nil
}

// Close closes the connection to the relay server
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// HandleError handles relay-specific errors
func (c *Client) HandleError(err error) error {
	switch err.Error() {
	case ErrInvalidToken:
		return fmt.Errorf("invalid token, please request a new one")
	case ErrRateLimitExceeded:
		time.Sleep(time.Second)
		return fmt.Errorf("rate limit exceeded, retrying")
	case ErrConnectionLimitReached:
		return fmt.Errorf("connection limit reached, please close unused connections")
	case ErrServerUnavailable:
		return fmt.Errorf("server unavailable, please try another server")
	default:
		return err
	}
} 