# Mihomo-ECH

> Mihomo with ECH-Tunnel Protocol Support

[English](#english) | [中文](#中文)

---

## 中文

### 📖 项目简介

**Mihomo-ECH** 是在 [Mihomo (Clash Meta)](https://github.com/MetaCubeX/mihomo) 基础上集成了 **ECH-Tunnel** 协议支持的增强版本。

通过本项目,你可以使用 Mihomo 作为统一客户端,连接到 ECH-Tunnel 或 ECH-Workers 服务端,享受强大的规则引擎和智能分流功能。

### ✨ 核心特性

- ✅ **ECH-Tunnel 协议支持** - 完整实现 ECH-Tunnel 客户端
- ✅ **ECH-Workers 支持** - 兼容 Cloudflare Workers 服务端
- ✅ **统一客户端** - 一个程序支持所有代理协议
- ✅ **智能分流** - 内置 GeoIP/GeoSite 规则引擎
- ✅ **自动切换** - 支持 url-test/fallback/load-balance
- ✅ **混合协议** - 可与 Trojan/VMess/Shadowsocks 等混用

### 🆚 与原版 Mihomo 的区别

| 特性 | 原版 Mihomo | Mihomo-ECH |
|------|------------|-----------|
| ECH-Tunnel 协议 | ❌ 不支持 | ✅ **完整支持** |
| ECH-Workers 协议 | ❌ 不支持 | ✅ **完整支持** |
| 其他协议 | ✅ 支持 | ✅ 支持 |
| 规则引擎 | ✅ 强大 | ✅ 强大 |

### 📦 新增内容

本项目在原版 Mihomo 基础上新增:

```
adapter/outbound/echtunnel.go    # ECH-Tunnel Adapter
transport/echtunnel/client.go    # ECH-Tunnel 客户端
transport/echtunnel/conn.go      # WebSocket 连接包装
constant/adapters.go             # 添加 ECHTunnel 类型
adapter/parser.go                # 添加 echtunnel 解析
```

### 🚀 快速开始

#### 1. 编译

```bash
# 克隆项目
git clone https://github.com/YOUR_USERNAME/mihomo-ech.git
cd mihomo-ech

# 安装依赖
go mod tidy

# 编译
go build -o mihomo.exe
```

#### 2. 配置

创建 `config.yaml`:

```yaml
port: 7890
socks-port: 7891
allow-lan: false
mode: rule
log-level: info

proxies:
  # ECH-Tunnel (VPS 服务端)
  - name: "ECH-VPS"
    type: echtunnel
    server: your-vps.com
    port: 443
    ws-path: /tunnel
    token: your-secret-token
    ech-domain: cloudflare-ech.com
    dns: https://dns.alidns.com/dns-query
    udp: true

  # ECH-Workers (Cloudflare Workers 服务端)
  - name: "ECH-Workers"
    type: echtunnel
    server: workers.example.com
    port: 443
    ws-path: /
    token: your-workers-token
    ech-domain: cloudflare-ech.com

proxy-groups:
  - name: "🚀 代理"
    type: select
    proxies:
      - ECH-VPS
      - ECH-Workers
      - DIRECT

rules:
  - GEOIP,CN,DIRECT
  - MATCH,🚀 代理
```

#### 3. 运行

```bash
.\mihomo.exe -f config.yaml
```

### 📝 配置说明

#### ECH-Tunnel 配置项

| 参数 | 必填 | 说明 | 示例 |
|------|------|------|------|
| `type` | ✅ | 协议类型 | `echtunnel` |
| `name` | ✅ | 节点名称 | `ECH-VPS` |
| `server` | ✅ | 服务器地址 | `vps.example.com` |
| `port` | ✅ | 服务器端口 | `443` |
| `ws-path` | ❌ | WebSocket 路径 | `/tunnel` (默认) |
| `token` | ❌ | 认证令牌 | `your-token` |
| `ech-domain` | ❌ | ECH 域名 | `cloudflare-ech.com` (默认) |
| `dns` | ❌ | DoH 服务器 | `https://dns.alidns.com/dns-query` |
| `ip` | ❌ | 指定 IP | `1.2.3.4` |
| `udp` | ❌ | 启用 UDP | `true` / `false` |

### 🎯 使用场景

#### 场景 1: 多服务端自动切换

```yaml
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

#### 场景 2: 混合协议负载均衡

```yaml
proxy-groups:
  - name: "负载均衡"
    type: load-balance
    proxies:
      - ECH-Tunnel-1
      - Trojan-1
      - VMess-1
```

#### 场景 3: 智能分流

```yaml
rules:
  - GEOIP,CN,DIRECT
  - GEOSITE,CN,DIRECT
  - GEOSITE,netflix,ECH-VPS-1
  - MATCH,自动选择
```

### 📚 文档

- [集成指南](ECH-TUNNEL-INTEGRATION-GUIDE.md) - 详细的集成步骤
- [实现计划](docs/mihomo-ech-integration.md) - 技术实现细节
- [Mihomo 官方文档](https://wiki.metacubex.one/) - Mihomo 使用文档

### 🤝 贡献

欢迎提交 Issue 和 Pull Request!

### 📄 许可证

本项目基于 [Mihomo](https://github.com/MetaCubeX/mihomo) 开发,遵循 GPL-3.0 许可证。

### 🙏 致谢

- [Mihomo (Clash Meta)](https://github.com/MetaCubeX/mihomo) - 强大的代理工具
- [ECH-Tunnel](https://github.com/...) - ECH-Tunnel 项目
- 所有贡献者

---

## English

### 📖 Introduction

**Mihomo-ECH** is an enhanced version of [Mihomo (Clash Meta)](https://github.com/MetaCubeX/mihomo) with **ECH-Tunnel** protocol support integrated.

With this project, you can use Mihomo as a unified client to connect to ECH-Tunnel or ECH-Workers servers, enjoying powerful rule engines and intelligent traffic routing.

### ✨ Features

- ✅ **ECH-Tunnel Protocol** - Full ECH-Tunnel client implementation
- ✅ **ECH-Workers Support** - Compatible with Cloudflare Workers backend
- ✅ **Unified Client** - One program for all proxy protocols
- ✅ **Smart Routing** - Built-in GeoIP/GeoSite rule engine
- ✅ **Auto Switch** - Support url-test/fallback/load-balance
- ✅ **Mixed Protocols** - Works with Trojan/VMess/Shadowsocks etc.

### 🚀 Quick Start

#### 1. Build

```bash
git clone https://github.com/YOUR_USERNAME/mihomo-ech.git
cd mihomo-ech
go mod tidy
go build -o mihomo.exe
```

#### 2. Configure

Create `config.yaml`:

```yaml
port: 7890
socks-port: 7891
mode: rule

proxies:
  - name: "ECH-Tunnel"
    type: echtunnel
    server: your-server.com
    port: 443
    token: your-token

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

#### 3. Run

```bash
./mihomo -f config.yaml
```

### 📝 Configuration

See [Integration Guide](ECH-TUNNEL-INTEGRATION-GUIDE.md) for detailed configuration options.

### 📚 Documentation

- [Integration Guide](ECH-TUNNEL-INTEGRATION-GUIDE.md)
- [Mihomo Wiki](https://wiki.metacubex.one/)

### 📄 License

GPL-3.0 License (same as Mihomo)

### 🙏 Credits

- [Mihomo (Clash Meta)](https://github.com/MetaCubeX/mihomo)
- [ECH-Tunnel Project](https://github.com/...)
- All contributors

---

## ⭐ Star History

If you find this project useful, please consider giving it a star!

---

**Made with ❤️ for the community**