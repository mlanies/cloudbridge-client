package relay

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/2gc-dev/cloudbridge-client/pkg/auth"
	"github.com/2gc-dev/cloudbridge-client/pkg/config"
	"github.com/2gc-dev/cloudbridge-client/pkg/errors"
	"github.com/2gc-dev/cloudbridge-client/pkg/heartbeat"
	"github.com/2gc-dev/cloudbridge-client/pkg/tunnel"
	"github.com/2gc-dev/cloudbridge-client/pkg/types"
)

// Client represents a CloudBridge Relay client
type Client struct {
	config         *types.Config
	conn           net.Conn
	encoder        *json.Encoder
	decoder        *json.Decoder
	authManager    *auth.AuthManager
	tunnelManager  *tunnel.Manager
	heartbeatMgr   *heartbeat.Manager
	retryStrategy  *errors.RetryStrategy
	mu             sync.RWMutex
	connected      bool
	clientID       string
	ctx            context.Context
	cancel         context.CancelFunc
}

// Message types as defined in the requirements
const (
	MessageTypeHello         = "hello"
	MessageTypeHelloResponse = "hello_response"
	MessageTypeAuth          = "auth"
	MessageTypeAuthResponse  = "auth_response"
	MessageTypeTunnelInfo    = "tunnel_info"
	MessageTypeTunnelResponse = "tunnel_response"
	MessageTypeHeartbeat     = "heartbeat"
	MessageTypeHeartbeatResponse = "heartbeat_response"
	MessageTypeError         = "error"
)

// NewClient creates a new CloudBridge Relay client
func NewClient(cfg *types.Config) (*Client, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create authentication manager
	authManager, err := auth.NewAuthManager(&auth.AuthConfig{
		Type:     cfg.Auth.Type,
		Secret:   cfg.Auth.Secret,
		Keycloak: cfg.Auth.Keycloak,
	})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create auth manager: %w", err)
	}

	// Create retry strategy
	retryStrategy := errors.NewRetryStrategy(
		cfg.RateLimiting.MaxRetries,
		cfg.RateLimiting.BackoffMultiplier,
		cfg.RateLimiting.MaxBackoff,
	)

	client := &Client{
		config:        cfg,
		authManager:   authManager,
		retryStrategy: retryStrategy,
		ctx:           ctx,
		cancel:        cancel,
	}

	// Create tunnel manager
	client.tunnelManager = tunnel.NewManager(client)

	// Create heartbeat manager
	client.heartbeatMgr = heartbeat.NewManager(client)

	return client, nil
}

// Connect establishes a connection to the relay server
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return fmt.Errorf("already connected")
	}

	// Create TLS config
	tlsConfig, err := config.CreateTLSConfig(c.config)
	if err != nil {
		return fmt.Errorf("failed to create TLS config: %w", err)
	}

	// Establish connection
	var conn net.Conn
	if tlsConfig != nil {
		conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", c.config.Relay.Host, c.config.Relay.Port), tlsConfig)
	} else {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", c.config.Relay.Host, c.config.Relay.Port))
	}

	if err != nil {
		return errors.NewRelayError(errors.ErrTLSHandshakeFailed, fmt.Sprintf("failed to connect: %v", err))
	}

	c.conn = conn
	c.encoder = json.NewEncoder(conn)
	c.decoder = json.NewDecoder(conn)

	// Send hello message
	if err := c.sendHello(); err != nil {
		conn.Close()
		return fmt.Errorf("failed to send hello: %w", err)
	}

	// Receive hello response
	if err := c.receiveHelloResponse(); err != nil {
		conn.Close()
		return fmt.Errorf("failed to receive hello response: %w", err)
	}

	c.connected = true
	return nil
}

// Authenticate authenticates with the relay server
func (c *Client) Authenticate(token string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return fmt.Errorf("not connected")
	}

	// Create auth message
	authMsg, err := c.authManager.CreateAuthMessage(token)
	if err != nil {
		return fmt.Errorf("failed to create auth message: %w", err)
	}

	// Send auth message
	if err := c.sendMessage(authMsg); err != nil {
		return fmt.Errorf("failed to send auth message: %w", err)
	}

	// Receive auth response
	response, err := c.receiveMessage()
	if err != nil {
		return fmt.Errorf("failed to receive auth response: %w", err)
	}

	// Check response type
	if response["type"] != MessageTypeAuthResponse {
		return fmt.Errorf("unexpected response type: %s", response["type"])
	}

	// Check status
	if status, ok := response["status"].(string); !ok || status != "ok" {
		errorMsg := "authentication failed"
		if msg, ok := response["error"].(string); ok {
			errorMsg = msg
		}
		return errors.NewRelayError(errors.ErrAuthenticationFailed, errorMsg)
	}

	// Store client ID
	if clientID, ok := response["client_id"].(string); ok {
		c.clientID = clientID
	}

	return nil
}

// CreateTunnel creates a tunnel with the specified parameters
func (c *Client) CreateTunnel(tunnelID string, localPort int, remoteHost string, remotePort int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return fmt.Errorf("not connected")
	}

	// Create tunnel info message
	tunnelMsg := map[string]interface{}{
		"type":        MessageTypeTunnelInfo,
		"tunnel_id":   tunnelID,
		"local_port":  localPort,
		"remote_host": remoteHost,
		"remote_port": remotePort,
	}

	// Send tunnel info message
	if err := c.sendMessage(tunnelMsg); err != nil {
		return fmt.Errorf("failed to send tunnel info: %w", err)
	}

	// Receive tunnel response
	response, err := c.receiveMessage()
	if err != nil {
		return fmt.Errorf("failed to receive tunnel response: %w", err)
	}

	// Check response type
	if response["type"] != MessageTypeTunnelResponse {
		return fmt.Errorf("unexpected response type: %s", response["type"])
	}

	// Check status
	if status, ok := response["status"].(string); !ok || status != "ok" {
		errorMsg := "tunnel creation failed"
		if msg, ok := response["error"].(string); ok {
			errorMsg = msg
		}
		return errors.NewRelayError(errors.ErrTunnelCreationFailed, errorMsg)
	}

	// Register tunnel with tunnel manager
	return c.tunnelManager.RegisterTunnel(tunnelID, localPort, remoteHost, remotePort)
}

// StartHeartbeat starts the heartbeat mechanism
func (c *Client) StartHeartbeat() error {
	if !c.connected {
		return fmt.Errorf("not connected")
	}

	return c.heartbeatMgr.Start()
}

// StopHeartbeat stops the heartbeat mechanism
func (c *Client) StopHeartbeat() {
	c.heartbeatMgr.Stop()
}

// Close closes the connection to the relay server
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cancel()

	if c.heartbeatMgr != nil {
		c.heartbeatMgr.Stop()
	}

	if c.conn != nil {
		c.connected = false
		return c.conn.Close()
	}

	return nil
}

// IsConnected returns true if the client is connected
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// GetClientID returns the client ID assigned by the relay server
func (c *Client) GetClientID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.clientID
}

// sendHello sends a hello message to the relay server
func (c *Client) sendHello() error {
	helloMsg := map[string]interface{}{
		"type":     MessageTypeHello,
		"version":  "1.0",
		"features": []string{"tls", "heartbeat", "tunnel_info"},
	}

	return c.sendMessage(helloMsg)
}

// receiveHelloResponse receives and validates hello response
func (c *Client) receiveHelloResponse() error {
	response, err := c.receiveMessage()
	if err != nil {
		return err
	}

	if response["type"] != MessageTypeHelloResponse {
		return fmt.Errorf("unexpected response type: %s", response["type"])
	}

	// Validate version
	if version, ok := response["version"].(string); !ok || version != "1.0" {
		return fmt.Errorf("unsupported protocol version: %v", response["version"])
	}

	return nil
}

// sendMessage sends a JSON message to the relay server
func (c *Client) sendMessage(msg map[string]interface{}) error {
	if c.encoder == nil {
		return fmt.Errorf("encoder not initialized")
	}

	return c.encoder.Encode(msg)
}

// receiveMessage receives a JSON message from the relay server
func (c *Client) receiveMessage() (map[string]interface{}, error) {
	if c.decoder == nil {
		return nil, fmt.Errorf("decoder not initialized")
	}

	var msg map[string]interface{}
	if err := c.decoder.Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	// Check for error messages
	if msg["type"] == MessageTypeError {
		code, _ := msg["code"].(string)
		message, _ := msg["message"].(string)
		return nil, errors.NewRelayError(code, message)
	}

	return msg, nil
}

// SendHeartbeat sends a heartbeat message
func (c *Client) SendHeartbeat() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return fmt.Errorf("not connected")
	}

	heartbeatMsg := map[string]interface{}{
		"type": MessageTypeHeartbeat,
	}

	if err := c.sendMessage(heartbeatMsg); err != nil {
		return errors.NewRelayError(errors.ErrHeartbeatFailed, fmt.Sprintf("failed to send heartbeat: %v", err))
	}

	// Receive heartbeat response
	response, err := c.receiveMessage()
	if err != nil {
		return errors.NewRelayError(errors.ErrHeartbeatFailed, fmt.Sprintf("failed to receive heartbeat response: %v", err))
	}

	if response["type"] != MessageTypeHeartbeatResponse {
		return errors.NewRelayError(errors.ErrHeartbeatFailed, "unexpected heartbeat response type")
	}

	return nil
}

// GetConfig returns the client configuration
func (c *Client) GetConfig() *types.Config {
	return c.config
}

// GetRetryStrategy returns the retry strategy
func (c *Client) GetRetryStrategy() *errors.RetryStrategy {
	return c.retryStrategy
} 