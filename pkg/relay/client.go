package relay

import (
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

// Client represents a CloudBridge Relay client
type Client struct {
	conn    net.Conn
	reader  *bufio.Reader
	writer  *bufio.Writer
	useTLS  bool
	config  *tls.Config

	missedHeartbeats int32
	stopHeartbeat   chan struct{}
}

// NewClient creates a new CloudBridge Relay client
func NewClient(useTLS bool, config *tls.Config) *Client {
	return &Client{
		useTLS: useTLS,
		config: config,
		stopHeartbeat: make(chan struct{}),
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