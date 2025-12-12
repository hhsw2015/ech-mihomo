# Mihomo-ECH 配置示例

## 基础配置

### 最小配置

```yaml
port: 7890
socks-port: 7891
mode: rule

proxies:
  - name: "ECH-Tunnel"
    type: echtunnel
    server: your-server.com
    port: 443

proxy-groups:
  - name: "Proxy"
    type: select
    proxies:
      - ECH-Tunnel
      - DIRECT

rules:
  - GEOIP,CN,DIRECT
  - MATCH,Proxy
```

---

## 完整配置

### ECH-Tunnel (VPS 服务端)

```yaml
proxies:
  - name: "ECH-VPS"
    type: echtunnel
    server: vps.example.com
    port: 443
    ws-path: /tunnel
    token: your-secret-token
    ech-domain: cloudflare-ech.com
    dns: https://dns.alidns.com/dns-query
    ip: 1.2.3.4
    udp: true
```

### ECH-Workers (Cloudflare Workers 服务端)

```yaml
proxies:
  - name: "ECH-Workers"
    type: echtunnel
    server: workers.example.com
    port: 443
    ws-path: /
    token: your-workers-token
    ech-domain: cloudflare-ech.com
```

---

## 高级配置

### 多服务端自动切换

```yaml
proxies:
  - name: "ECH-VPS-1"
    type: echtunnel
    server: vps1.example.com
    port: 443
    token: token1

  - name: "ECH-VPS-2"
    type: echtunnel
    server: vps2.example.com
    port: 443
    token: token2

  - name: "ECH-Workers-1"
    type: echtunnel
    server: workers1.example.com
    port: 443
    token: token3

proxy-groups:
  - name: "自动选择"
    type: url-test
    proxies:
      - ECH-VPS-1
      - ECH-VPS-2
      - ECH-Workers-1
    url: 'https://www.gstatic.com/generate_204'
    interval: 300
```

### 混合协议配置

```yaml
proxies:
  # ECH-Tunnel
  - name: "ECH-1"
    type: echtunnel
    server: ech.example.com
    port: 443
    token: xxx

  # Trojan
  - name: "Trojan-1"
    type: trojan
    server: trojan.example.com
    port: 443
    password: yyy

  # VMess
  - name: "VMess-1"
    type: vmess
    server: vmess.example.com
    port: 443
    uuid: zzz

proxy-groups:
  - name: "负载均衡"
    type: load-balance
    proxies:
      - ECH-1
      - Trojan-1
      - VMess-1
    strategy: consistent-hashing
```

### 智能分流配置

```yaml
proxy-groups:
  - name: "🚀 代理"
    type: select
    proxies:
      - 自动选择
      - ECH-VPS
      - ECH-Workers
      - DIRECT

  - name: "自动选择"
    type: url-test
    proxies:
      - ECH-VPS
      - ECH-Workers

  - name: "🎬 流媒体"
    type: select
    proxies:
      - ECH-VPS
      - ECH-Workers

  - name: "🎯 直连"
    type: select
    proxies:
      - DIRECT

rules:
  # 局域网直连
  - DOMAIN-SUFFIX,local,🎯 直连
  - IP-CIDR,192.168.0.0/16,🎯 直连
  - IP-CIDR,10.0.0.0/8,🎯 直连
  - IP-CIDR,172.16.0.0/12,🎯 直连
  - IP-CIDR,127.0.0.0/8,🎯 直连

  # 国内直连
  - GEOIP,CN,🎯 直连
  - GEOSITE,CN,🎯 直连

  # 流媒体
  - GEOSITE,netflix,🎬 流媒体
  - GEOSITE,youtube,🎬 流媒体
  - GEOSITE,disney,🎬 流媒体

  # 其他走代理
  - MATCH,🚀 代理
```

---

## 完整示例配置

```yaml
# Mihomo-ECH 完整配置示例

# 基础设置
port: 7890
socks-port: 7891
mixed-port: 7892
allow-lan: false
bind-address: '*'
mode: rule
log-level: info
ipv6: true

# DNS 设置
dns:
  enable: true
  listen: 0.0.0.0:53
  enhanced-mode: fake-ip
  fake-ip-range: 198.18.0.1/16
  nameserver:
    - https://dns.alidns.com/dns-query
    - https://doh.pub/dns-query
  fallback:
    - https://1.1.1.1/dns-query
    - https://dns.google/dns-query

# 代理配置
proxies:
  # ECH-Tunnel VPS
  - name: "ECH-VPS-HK"
    type: echtunnel
    server: hk.example.com
    port: 443
    ws-path: /tunnel
    token: your-token-1
    ech-domain: cloudflare-ech.com
    dns: https://dns.alidns.com/dns-query
    udp: true

  - name: "ECH-VPS-US"
    type: echtunnel
    server: us.example.com
    port: 443
    ws-path: /tunnel
    token: your-token-2
    ech-domain: cloudflare-ech.com
    udp: true

  # ECH-Workers
  - name: "ECH-Workers-CF"
    type: echtunnel
    server: workers.example.com
    port: 443
    ws-path: /
    token: your-workers-token
    ech-domain: cloudflare-ech.com

# 代理组
proxy-groups:
  - name: "🚀 代理"
    type: select
    proxies:
      - 自动选择
      - 香港节点
      - 美国节点
      - Workers
      - DIRECT

  - name: "自动选择"
    type: url-test
    proxies:
      - ECH-VPS-HK
      - ECH-VPS-US
      - ECH-Workers-CF
    url: 'https://www.gstatic.com/generate_204'
    interval: 300

  - name: "香港节点"
    type: select
    proxies:
      - ECH-VPS-HK

  - name: "美国节点"
    type: select
    proxies:
      - ECH-VPS-US

  - name: "Workers"
    type: select
    proxies:
      - ECH-Workers-CF

  - name: "🎬 流媒体"
    type: select
    proxies:
      - 香港节点
      - 美国节点

  - name: "🎯 直连"
    type: select
    proxies:
      - DIRECT

# 规则
rules:
  # 局域网
  - DOMAIN-SUFFIX,local,🎯 直连
  - IP-CIDR,192.168.0.0/16,🎯 直连
  - IP-CIDR,10.0.0.0/8,🎯 直连
  - IP-CIDR,172.16.0.0/12,🎯 直连
  - IP-CIDR,127.0.0.0/8,🎯 直连

  # 流媒体
  - GEOSITE,netflix,🎬 流媒体
  - GEOSITE,youtube,🎬 流媒体
  - GEOSITE,disney,🎬 流媒体
  - GEOSITE,hbo,🎬 流媒体

  # 国内直连
  - GEOIP,CN,🎯 直连
  - GEOSITE,CN,🎯 直连

  # 默认代理
  - MATCH,🚀 代理
```

---

## 参数说明

### ECH-Tunnel 配置参数

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `type` | string | ✅ | - | 必须为 `echtunnel` |
| `name` | string | ✅ | - | 节点名称 |
| `server` | string | ✅ | - | 服务器地址 |
| `port` | int | ✅ | - | 服务器端口 |
| `ws-path` | string | ❌ | `/tunnel` | WebSocket 路径 |
| `token` | string | ❌ | - | 认证令牌 |
| `ech-domain` | string | ❌ | `cloudflare-ech.com` | ECH 域名 |
| `dns` | string | ❌ | - | DoH 服务器 |
| `ip` | string | ❌ | - | 指定 IP 地址 |
| `udp` | bool | ❌ | `false` | 启用 UDP |

---

## 使用建议

1. **服务端选择**
   - VPS 服务端: 延迟低,稳定性好
   - Workers 服务端: 免费,但可能有限制

2. **代理组配置**
   - 使用 `url-test` 实现自动切换
   - 使用 `fallback` 实现故障转移
   - 使用 `load-balance` 实现负载均衡

3. **规则配置**
   - 国内网站直连,节省流量
   - 流媒体使用特定节点
   - 其他流量走自动选择

4. **性能优化**
   - 启用 UDP 支持
   - 合理设置测试间隔
   - 使用 DNS 缓存
