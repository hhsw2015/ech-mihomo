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

	// 2. 定义一个能够完成 "TCP + TLS + ECH" 的拨号函数
	dialFn := func(ctx context.Context, network, address string) (net.Conn, error) {
		// 强制使用配置的服务器地址进行 TCP 连接
		c, err := e.dialer.DialContext(ctx, "tcp", e.addr)
		if err != nil {
			return nil, err
		}

		// 使用 vmess.StreamTLSConn 进行 TLS 握手 (支持 ECH)
		// 注意: StreamTLSConn 的 host 参数应该是内层 SNI (这里是 option.Server)
		// 如果开启了 ECH, StreamTLSConn 会自动使用 echConfig 中的 PublicName 作为外层 SNI
		tlsOpts := &vmess.TLSConfig{
			Host:              option.Server, // Host Header (用于 websocket) & Target SNI
			SkipCertVerify:    option.SkipCertVerify,
			FingerPrint:       option.Fingerprint,
			ClientFingerprint: option.ClientFingerprint,
			NextProtos:        []string{"http/1.1"},
			ECH:               e.echConfig,
		}

		return vmess.StreamTLSConn(ctx, c, tlsOpts)
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
