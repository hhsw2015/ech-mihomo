package outbound

import (
	"context"
	"net"
	"strconv"

	"github.com/metacubex/mihomo/component/ech"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/transport/echtunnel"
	"github.com/metacubex/mihomo/transport/vmess"
)

type ECHTunnel struct {
	*Base
	option    *ECHTunnelOption
	dialer    C.Dialer
	client    *echtunnel.Client
	echConfig *ech.Config
}

type ECHTunnelOption struct {
	BasicOption
	Name              string     `proxy:"name"`
	Server            string     `proxy:"server"`
	Port              int        `proxy:"port"`
	WSPath            string     `proxy:"ws-path,omitempty"`
	Token             string     `proxy:"token,omitempty"`
	ECHDomain         string     `proxy:"ech-domain,omitempty"`
	DNS               string     `proxy:"dns,omitempty"`
	IP                string     `proxy:"ip,omitempty"`
	UDP               bool       `proxy:"udp,omitempty"`
	ECHOpts           ECHOptions `proxy:"ech-opts,omitempty"`
	SkipCertVerify    bool       `proxy:"skip-cert-verify,omitempty"`
	Fingerprint       string     `proxy:"fingerprint,omitempty"`
	ClientFingerprint string     `proxy:"client-fingerprint,omitempty"`
}

// StreamConnContext implements C.ProxyAdapter
func (e *ECHTunnel) StreamConnContext(ctx context.Context, c net.Conn, metadata *C.Metadata) (_ net.Conn, err error) {
	return e.client.DialContext(ctx, c.RemoteAddr().String())
}

// DialContext implements C.ProxyAdapter
func (e *ECHTunnel) DialContext(ctx context.Context, metadata *C.Metadata) (_ C.Conn, err error) {
	c, err := e.client.DialContext(ctx, metadata.String())
	if err != nil {
		return nil, err
	}
	return NewConn(c, e), nil
}

// ListenPacketContext implements C.ProxyAdapter
func (e *ECHTunnel) ListenPacketContext(ctx context.Context, metadata *C.Metadata) (_ C.PacketConn, err error) {
	c, err := e.client.DialContext(ctx, metadata.String())
	if err != nil {
		return nil, err
	}
	pc := echtunnel.NewPacketConn(c)
	return newPacketConn(pc, e), nil
}

// SupportUOT implements C.ProxyAdapter
func (e *ECHTunnel) SupportUOT() bool {
	return true
}

// ProxyInfo implements C.ProxyAdapter
func (e *ECHTunnel) ProxyInfo() C.ProxyInfo {
	info := e.Base.ProxyInfo()
	info.DialerProxy = e.option.DialerProxy
	return info
}

// Close implements C.ProxyAdapter
func (e *ECHTunnel) Close() error {
	if e.client != nil {
		return e.client.Close()
	}
	return nil
}

func NewECHTunnel(option ECHTunnelOption) (*ECHTunnel, error) {
	addr := net.JoinHostPort(option.Server, strconv.Itoa(option.Port))

	// 处理默认值
	if option.ECHDomain == "" {
		option.ECHDomain = "cloudflare-ech.com"
	}

	echConfig, err := option.ECHOpts.Parse()
	if err != nil {
		return nil, err
	}

	e := &ECHTunnel{
		Base: &Base{
			name:   option.Name,
			addr:   addr,
			tp:     C.ECHTunnel,
			pdName: option.ProviderName,
			udp:    option.UDP,
			tfo:    option.TFO,
			mpTcp:  option.MPTCP,
			iface:  option.Interface,
			rmark:  option.RoutingMark,
			prefer: option.IPVersion,
		},
		option:    &option,
		echConfig: echConfig,
	}

	// 1. 先创建基础 TCP 安全拨号器 (Mihomo 内部设施)
	e.dialer = option.NewDialer(e.DialOptions())

	// 2. 定义一个能够完成 "TCP + TLS + ECH + WebSocket" 的拨号函数
	dialFn := func(ctx context.Context, network, address string) (net.Conn, error) {
		// A. 拨号基础 TCP
		c, err := e.dialer.DialContext(ctx, "tcp", e.addr)
		if err != nil {
			return nil, err
		}

		// B. 使用 Mihomo 强大的 vmess.StreamWebsocketConn
		// 它会自动处理:
		// 1. TLS 握手 (包括 ECH, 如果启用)
		// 2. WebSocket 握手 (Upgrade, Host Header, etc)
		// 3. Early Data (如果启用)
		
		wsConfig := &vmess.WebsocketConfig{
			Host:      option.Server, // Host Header & SNI
			Port:      strconv.Itoa(option.Port),
			Path:      option.WSPath,
			TLS:       true,
			Headers:   http.Header{},
			// ECH 配置: 如果有 ECHOpts, ECHConfig 就不为空
			ECHConfig: e.echConfig,
		}
		
		if option.WSPath == "" {
			wsConfig.Path = "/tunnel"
		}
		
		wsConfig.TLSConfig = &tls.Config{
			ServerName:         option.Server,
			InsecureSkipVerify: option.SkipCertVerify,
			NextProtos:         []string{"http/1.1"},
			MinVersion:         tls.VersionTLS12,
		}

		// 如果 Token 存在, 加到 Header (部分实现可能通过 Header 传 Token)
		if option.Token != "" {
			wsConfig.Headers.Set("Sec-WebSocket-Protocol", option.Token)
		}
		
		return vmess.StreamWebsocketConn(ctx, c, wsConfig)
	}

	// 3. 将这个超级拨号器传给 WebSocket 客户端
	// 注意: 因为 vmess.StreamWebsocketConn 已经完成了 WebSocket 握手
	// 所以返回的 conn 已经是一个 "类似 TCP 的顺滑数据流" 了
	// 我们的 client.go 不需要再做 WebSocket 握手了!
	// 这意味着我们需要大改 client.go 或者...
	// 等等, 为了最小化改动, 我们可以让 client.go 继续认为它是 "ws" 协议
	// 但是如果 conn 已经是 ws 这一层了, gorilla 再去握手就会出错 (握手包发给数据流)
	
	// ---> 修正方案:
	// 既然 vmess 组件这么强, 我们其实不需要 gorilla 了
	// 我们只需要一个简单的 Wrapper 把它封装成 net.Conn 
	// 但是为了兼容之前的架构, 我们这里保留 DialFn 的设计, 
	// 但在 client.go 里, 我们如果发现是 "Managed Connection", 就不应该再 Dial 了
	
	// ... 重新思考 ...
	// 上面的方案虽然好, 但需要改 client.go 改动太大
	// 我们回退到: 用 vmess 处理 TLS, 用 gorilla 处理 WS (这是之前失败的方案)
	// 失败原因是: Host Header 没发对 / Cloudflare 校验严
	// 还是坚持现在的修复: 显式加 Host Header
	
	// (此段代码仅为思维过程, 下面是真正的代码)

    // 真正的修复逻辑:
	// 继续使用 vmess.StreamTLSConn 处理 TLS + ECH
	// 在 client.go 里确保发 Host Header
	// (用户之前的测试表明 client.go 加了 Host 还没通, 可能是 SNI 和 Host 不一致导致 CF 拒绝)
	// Cloudflare Workers 要求 SNI 和 Host 必须一致
	
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		c, err := e.dialer.DialContext(ctx, "tcp", e.addr)
		if err != nil {
			return nil, err
		}

		tlsOpts := &vmess.TLSConfig{
			Host:              option.Server, 
			SkipCertVerify:    option.SkipCertVerify,
			FingerPrint:       option.Fingerprint,
			ClientFingerprint: option.ClientFingerprint,
			NextProtos:        []string{"http/1.1"},
			ECH:               e.echConfig,
		}

		return vmess.StreamTLSConn(ctx, c, tlsOpts)
	} (ctx, network, address)
	}

	// 3. 将这个超级拨号器传给 WebSocket 客户端
	client, err := echtunnel.NewClient(echtunnel.Config{
		Server:    option.Server,
		Port:      option.Port,
		WSPath:    option.WSPath,
		Token:     option.Token,
		ECHDomain: option.ECHDomain, // 这里的 ECHDomain 其实在拨号器里已经不重要了
		DNS:       option.DNS,
		IP:        option.IP,
	}, dialFn)

	if err != nil {
		return nil, err
	}

	e.client = client
	return e, nil
}
