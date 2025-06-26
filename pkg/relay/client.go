package relay

import (
<<<<<<< HEAD
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"sync/atomic"
	"sync"
	"time"
)

// Message types
const (
	MessageTypeHello      = "hello"
	MessageTypeAuth       = "auth"
	MessageTypeAuthResp   = "auth_response"
	MessageTypeRegister   = "register"
	MessageTypeRegisterResp = "register_response"
	MessageTypeHeartbeat  = "heartbeat"
	MessageTypeHeartbeatResp = "heartbeat_response"
	MessageTypeError      = "error"

	MaxMessageSize           = 1024 * 1024 // 1MB
	ConnectTimeout           = 10 * time.Second
	ReadWriteTimeout         = 30 * time.Second
	HeartbeatInterval        = 30 * time.Second
	HeartbeatTimeout         = 5 * time.Second
	MaxMissedHeartbeats      = 3
	MaxErrorWindowSeconds = 60
	MaxErrorCount         = 3
)

// Error types
const (
	ErrInvalidToken          = "invalid_token"
	ErrRateLimitExceeded     = "rate_limit_exceeded"
	ErrConnectionLimitReached = "connection_limit_reached"
	ErrServerUnavailable     = "server_unavailable"
)

// AuthMessage represents an authentication message
type AuthMessage struct {
	Type       string                 `json:"type"`
	Token      string                 `json:"token"`
	Version    string                 `json:"version"`
	ClientInfo map[string]interface{} `json:"client_info"`
}

type AuthResponse struct {
	Type       string                 `json:"type"`
	Status     string                 `json:"status"`
	ServerInfo map[string]interface{} `json:"server_info"`
}

type RegisterMessage struct {
	Type       string                 `json:"type"`
	TunnelInfo map[string]interface{} `json:"tunnel_info"`
}

type RegisterResponse struct {
	Type    string                 `json:"type"`
	Status  string                 `json:"status"`
	TunnelID string                `json:"tunnel_id"`
	Config  map[string]interface{} `json:"config"`
}

type HeartbeatMessage struct {
	Type     string                 `json:"type"`
	TunnelID string                 `json:"tunnel_id"`
	Stats    map[string]interface{} `json:"stats"`
}

type HeartbeatResponse struct {
	Type       string `json:"type"`
	Status     string `json:"status"`
	ServerTime string `json:"server_time"`
}

type ErrorMessage struct {
	Type    string                 `json:"type"`
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// Tunnel represents a managed tunnel connection
type Tunnel struct {
	ID          string
	LocalPort   int
	RemoteHost  string
	RemotePort  int
	Protocol    string
	Options     map[string]interface{}
	stopChan    chan struct{}
	proxyCmd    *exec.Cmd
}

// Client represents a CloudBridge Relay client
type Client struct {
	conn    net.Conn
	reader  *bufio.Reader
	writer  *bufio.Writer
	useTLS  bool
	config  *tls.Config

	missedHeartbeats int32
	stopHeartbeat   chan struct{}
	tunnels         map[string]*Tunnel
	tunnelMutex     sync.RWMutex
}

// NewClient creates a new CloudBridge Relay client
func NewClient(useTLS bool, config *tls.Config) *Client {
	return &Client{
		useTLS: useTLS,
		config: config,
		stopHeartbeat: make(chan struct{}),
		tunnels: make(map[string]*Tunnel),
	}
}

// Connect establishes a connection to the relay server
func (c *Client) Connect(host string, port int) error {
	var err error
	var conn net.Conn
	dialer := &net.Dialer{Timeout: ConnectTimeout}
	address := fmt.Sprintf("%s:%d", host, port)

	if c.useTLS {
		conn, err = tls.DialWithDialer(dialer, "tcp", address, c.config)
	} else {
		conn, err = dialer.Dial("tcp", address)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to relay: %w", err)
	}

	c.conn = conn
	c.reader = bufio.NewReaderSize(conn, MaxMessageSize)
	c.writer = bufio.NewWriter(conn)
	return nil
}

// Close closes the connection to the relay server
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// sendMessage отправляет JSON-сообщение с \n
func (c *Client) sendMessage(msg interface{}) error {
	c.conn.SetWriteDeadline(time.Now().Add(ReadWriteTimeout))
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if len(data) > MaxMessageSize {
		return fmt.Errorf("message too large")
	}
	if _, err := c.writer.Write(append(data, '\n'));
		err != nil {
		return err
	}
	return c.writer.Flush()
}

// readMessage читает строку, парсит JSON, ограничивает размер
func (c *Client) readMessage() (map[string]interface{}, error) {
	c.conn.SetReadDeadline(time.Now().Add(ReadWriteTimeout))
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if len(line) > MaxMessageSize {
		return nil, fmt.Errorf("message too large")
	}
	line = strings.TrimSpace(line)
	var msg map[string]interface{}
	if err := json.Unmarshal([]byte(line), &msg); err != nil {
		return nil, err
	}
	return msg, nil
}

// Handshake: ждет hello, отправляет auth, ждет auth_response
func (c *Client) Handshake(token string, version string) error {
	// 1. Ждем hello
	hello, err := c.readMessage()
	if err != nil {
		return fmt.Errorf("failed to read hello: %w", err)
	}
	if t, _ := hello["type"].(string); t != MessageTypeHello {
		return fmt.Errorf("expected hello, got: %v", hello)
	}
	log.Printf("Received hello: %v", hello)

	// 2. Отправляем auth
	auth := AuthMessage{
		Type:    MessageTypeAuth,
		Token:   token,
		Version: version,
		ClientInfo: map[string]interface{}{
			"os":      runtime.GOOS,
			"version": version,
		},
	}
	if err := c.sendMessage(auth); err != nil {
		return fmt.Errorf("failed to send auth: %w", err)
	}

	// 3. Ждем auth_response
	resp, err := c.readMessage()
	if err != nil {
		return fmt.Errorf("failed to read auth_response: %w", err)
	}
	if t, _ := resp["type"].(string); t != MessageTypeAuthResp {
		return fmt.Errorf("expected auth_response, got: %v", resp)
	}
	if resp["status"] != "success" {
		return fmt.Errorf("auth failed: %v", resp)
	}
	log.Printf("Auth success: %v", resp)
	return nil
}

// startHeartbeat запускает автоматическую отправку heartbeat
func (c *Client) startHeartbeat(tunnelID string) {
	go func() {
		ticker := time.NewTicker(HeartbeatInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				msg := HeartbeatMessage{
					Type:     MessageTypeHeartbeat,
					TunnelID: tunnelID,
					Stats:    map[string]interface{}{},
				}
				if err := c.sendMessage(msg); err != nil {
					log.Printf("Failed to send heartbeat: %v", err)
				}
				atomic.AddInt32(&c.missedHeartbeats, 1)
				if atomic.LoadInt32(&c.missedHeartbeats) > MaxMissedHeartbeats {
					log.Printf("Missed too many heartbeats, closing connection")
					c.Close()
					return
				}
			case <-c.stopHeartbeat:
				return
			}
		}
	}()
}

// stopHeartbeatLoop останавливает heartbeat loop
func (c *Client) stopHeartbeatLoop() {
	close(c.stopHeartbeat)
}

// startLocalProxy запускает TCP-прокси с локального порта на remote_host:remote_port
func startLocalProxy(tunnel map[string]interface{}) {
	localPort := int(tunnel["local_port"].(float64))
	remoteHost := tunnel["remote_host"].(string)
	remotePort := int(tunnel["remote_port"].(float64))

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", localPort))
	if err != nil {
		log.Printf("Failed to start proxy: %v", err)
		return
	}
	log.Printf("Proxy started on 127.0.0.1:%d -> %s:%d", localPort, remoteHost, remotePort)
	go func() {
		for {
			clientConn, err := ln.Accept()
			if err != nil {
				log.Printf("Proxy accept error: %v", err)
				continue
			}
			go func() {
				defer clientConn.Close()
				serverConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", remoteHost, remotePort))
				if err != nil {
					log.Printf("Proxy dial error: %v", err)
					return
				}
				defer serverConn.Close()
				go io.Copy(serverConn, clientConn)
				io.Copy(clientConn, serverConn)
			}()
		}
	}()
}

type errorStats struct {
	mu      sync.Mutex
	history map[string][]time.Time
}

func newErrorStats() *errorStats {
	return &errorStats{history: make(map[string][]time.Time)}
}

func (e *errorStats) add(code string) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := time.Now()
	window := now.Add(-MaxErrorWindowSeconds * time.Second)
	lst := e.history[code]
	// Оставляем только ошибки за последний window
	filtered := make([]time.Time, 0, len(lst))
	for _, t := range lst {
		if t.After(window) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, now)
	e.history[code] = filtered
	return len(filtered)
}

// EventLoop - основной цикл обработки сообщений (после аутентификации)
func (c *Client) EventLoop(onTunnel func(tunnelInfo map[string]interface{})) error {
	var tunnelID string
	errStats := newErrorStats()
	for {
		msg, err := c.readMessage()
		if err != nil {
			if err == io.EOF {
				log.Println("Connection closed by server")
				return nil
			}
			return fmt.Errorf("failed to decode message: %w", err)
		}
		typeVal, _ := msg["type"].(string)
		switch typeVal {
		case "tunnel_info":
			log.Printf("Tunnel info: %v", msg)
			if tunnel, ok := msg["tunnel_id"].(string); ok {
				log.Printf("Tunnel ID: %s", tunnel)
			}
			go startLocalProxy(msg)
		case MessageTypeRegister:
			log.Printf("Register tunnel: %v", msg)
			if onTunnel != nil {
				onTunnel(msg["tunnel_info"].(map[string]interface{}))
			}
			// Отправить register_response (заглушка)
			resp := RegisterResponse{
				Type:    MessageTypeRegisterResp,
				Status:  "success",
				TunnelID: "tunnel-uuid",
				Config:  map[string]interface{}{"endpoint": "relay.2gc.ru:443"},
			}
			c.sendMessage(resp)
			// Запустить heartbeat loop для этого туннеля
			tunnelID = resp.TunnelID
			c.startHeartbeat(tunnelID)
		case MessageTypeHeartbeatResp:
			log.Printf("Heartbeat response: %v", msg)
			atomic.StoreInt32(&c.missedHeartbeats, 0)
		case MessageTypeHeartbeat:
			log.Printf("Heartbeat: %v", msg)
			resp := HeartbeatResponse{
				Type:       MessageTypeHeartbeatResp,
				Status:     "ok",
				ServerTime: time.Now().Format(time.RFC3339),
			}
			c.sendMessage(resp)
		case MessageTypeError:
			log.Printf("Error: %v", msg)
			code, _ := msg["code"].(string)
			count := errStats.add(code)
			if count >= MaxErrorCount {
				log.Printf("Too many errors of type %s in %d seconds, closing connection", code, MaxErrorWindowSeconds)
				return fmt.Errorf("error threshold exceeded for %s", code)
			}
			// Можно добавить обработку по коду ошибки (например, reconnect, backoff и т.д.)
		default:
			log.Printf("Unknown message type: %v", msg)
		}
	}
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

// CreateTunnel creates a new tunnel based on tunnel info
func (c *Client) CreateTunnel(tunnelInfo map[string]interface{}) (*Tunnel, error) {
	tunnel := &Tunnel{
		ID:         tunnelInfo["tunnel_id"].(string),
		LocalPort:  int(tunnelInfo["local_port"].(float64)),
		RemoteHost: tunnelInfo["remote_host"].(string),
		RemotePort: int(tunnelInfo["remote_port"].(float64)),
		Protocol:   tunnelInfo["protocol"].(string),
		Options:    tunnelInfo["options"].(map[string]interface{}),
		stopChan:   make(chan struct{}),
	}

	c.tunnelMutex.Lock()
	c.tunnels[tunnel.ID] = tunnel
	c.tunnelMutex.Unlock()

	if err := tunnel.start(); err != nil {
		c.tunnelMutex.Lock()
		delete(c.tunnels, tunnel.ID)
		c.tunnelMutex.Unlock()
		return nil, err
	}

	return tunnel, nil
}

// CloseTunnel closes and removes a tunnel
func (c *Client) CloseTunnel(tunnelID string) error {
	c.tunnelMutex.Lock()
	tunnel, exists := c.tunnels[tunnelID]
	if !exists {
		c.tunnelMutex.Unlock()
		return fmt.Errorf("tunnel %s not found", tunnelID)
	}
	delete(c.tunnels, tunnelID)
	c.tunnelMutex.Unlock()

	return tunnel.stop()
}

// GetTunnel returns a tunnel by ID
func (c *Client) GetTunnel(tunnelID string) (*Tunnel, bool) {
	c.tunnelMutex.RLock()
	defer c.tunnelMutex.RUnlock()
	tunnel, exists := c.tunnels[tunnelID]
	return tunnel, exists
}

// ListTunnels returns all active tunnels
func (c *Client) ListTunnels() []*Tunnel {
	c.tunnelMutex.RLock()
	defer c.tunnelMutex.RUnlock()
	tunnels := make([]*Tunnel, 0, len(c.tunnels))
	for _, tunnel := range c.tunnels {
		tunnels = append(tunnels, tunnel)
	}
	return tunnels
}

// Tunnel methods
func (t *Tunnel) start() error {
	var cmd *exec.Cmd
	switch t.Protocol {
	case "rdp":
		cmd = t.startRDPProxy()
	case "ssh":
		cmd = t.startSSHProxy()
	case "http", "https":
		cmd = t.startHTTPProxy()
	default:
		return fmt.Errorf("unsupported protocol: %s", t.Protocol)
	}

	if cmd == nil {
		return fmt.Errorf("failed to start proxy for protocol %s", t.Protocol)
	}

	t.proxyCmd = cmd
	go t.monitorProxy()
	return nil
}

func (t *Tunnel) stop() error {
	close(t.stopChan)
	if t.proxyCmd != nil && t.proxyCmd.Process != nil {
		return t.proxyCmd.Process.Kill()
	}
	return nil
}

func (t *Tunnel) monitorProxy() {
	if t.proxyCmd == nil {
		return
	}

	done := make(chan error, 1)
	go func() {
		done <- t.proxyCmd.Wait()
	}()

	select {
	case err := <-done:
		log.Printf("Tunnel %s proxy stopped: %v", t.ID, err)
	case <-t.stopChan:
		log.Printf("Tunnel %s stopping proxy", t.ID)
	}
}

func (t *Tunnel) startRDPProxy() *exec.Cmd {
	cmd := exec.Command("xfreerdp",
		fmt.Sprintf("/v:%s", t.RemoteHost),
		fmt.Sprintf("/port:%d", t.RemotePort),
		fmt.Sprintf("/u:%s", t.Options["username"]),
		fmt.Sprintf("/p:%s", t.Options["password"]),
		fmt.Sprintf("/d:%s", t.Options["domain"]))
	return cmd
}

func (t *Tunnel) startSSHProxy() *exec.Cmd {
	cmd := exec.Command("ssh",
		"-L", fmt.Sprintf("%d:%s:%d", t.LocalPort, t.RemoteHost, t.RemotePort),
		fmt.Sprintf("%s@%s", t.Options["username"], t.RemoteHost))
	return cmd
}

func (t *Tunnel) startHTTPProxy() *exec.Cmd {
	// Implement HTTP proxy based on your requirements
	return nil
=======
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
>>>>>>> ebb63d9 (feat: implement CloudBridge Relay Client with TLS 1.3, JWT auth, tunnels, heartbeat, rate limiting, comprehensive docs and tests)
} 