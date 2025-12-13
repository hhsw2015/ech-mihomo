package echtunnel

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Config struct {
	Server    string
	Port      int
	WSPath    string
	Token     string
	ECHDomain string
	DNS       string
	IP        string
}

// 定义拨号函数类型
type DialFn func(ctx context.Context, network, address string) (net.Conn, error)

type Client struct {
	config Config
	dialer *websocket.Dialer
}

func NewClient(config Config, dialFn DialFn) (*Client, error) {
	// 设置默认值
	if config.WSPath == "" {
		config.WSPath = "/tunnel"
	}
	if config.ECHDomain == "" {
		config.ECHDomain = "cloudflare-ech.com"
	}

	// TLS 配置
	tlsConfig := &tls.Config{
		ServerName: config.ECHDomain,
		MinVersion: tls.VersionTLS13,
		// TODO: 添加 ECH 配置
		// 可以使用 Mihomo 现有的 ECH 支持
	}

	dialer := &websocket.Dialer{
		TLSClientConfig:  tlsConfig,
		HandshakeTimeout: 30 * time.Second,
		NetDialContext:   dialFn, // 使用传入的安全拨号器
	}

	return &Client{
		config: config,
		dialer: dialer,
	}, nil
}

func (c *Client) DialContext(ctx context.Context, address string) (net.Conn, error) {
	// 构建 WebSocket URL
	// 注意: 这里使用 ws 协议而不是 wss, 因为 dialer.NetDialContext 已经返回了安全连接(TLS/ECH)
	// 如果这里用 wss, gorilla 会试图再包裹一层 TLS, 导致错误
	u := url.URL{
		Scheme: "ws",
		Host:   net.JoinHostPort(c.config.Server, strconv.Itoa(c.config.Port)),
		Path:   c.config.WSPath,
	}

	// 添加 Token
	if c.config.Token != "" {
		q := u.Query()
		q.Set("token", c.config.Token)
		u.RawQuery = q.Encode()
	}

	// 建立 WebSocket 连接
	// 显式设置 Host Header, 这对 CDN (Cloudflare) 非常重要
	headers := http.Header{}
	headers.Set("Host", c.config.Server)
	headers.Set("User-Agent", "Mihomo/1.0 ECHTunnel/0.1")

	// 许多实现(如 v2ray-core, xray) 使用 Sec-WebSocket-Protocol 传递 Token/UUID
	if c.config.Token != "" {
		headers.Set("Sec-WebSocket-Protocol", c.config.Token)
	}

	conn, resp, err := c.dialer.DialContext(ctx, u.String(), headers)
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("websocket dial failed:Status=%s, err=%w", resp.Status, err)
		}
		return nil, fmt.Errorf("websocket dial failed: %w", err)
	}

	// 发送目标地址
	err = conn.WriteMessage(websocket.TextMessage, []byte(address))
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("send target address failed: %w", err)
	}

	// 返回包装后的连接
	return &WSConn{Conn: conn}, nil
}

func (c *Client) Close() error {
	return nil
}
