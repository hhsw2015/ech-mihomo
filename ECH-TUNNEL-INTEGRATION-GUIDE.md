# Mihomo ECH-Tunnel 集成 - 修改指南

## ✅ 已创建的文件

1. ✅ `adapter/outbound/echtunnel.go` - ECH-Tunnel Adapter
2. ✅ `transport/echtunnel/client.go` - 客户端实现
3. ✅ `transport/echtunnel/conn.go` - 连接包装

---

## 📝 需要手动修改的文件

### 1. `constant/adapters.go`

**位置 1**: 第 47 行后添加 (在 `Sudoku` 后面)

```go
const (
    Direct AdapterType = iota
    Reject
    RejectDrop
    Compatible
    Pass
    Dns

    Relay
    Selector
    Fallback
    URLTest
    LoadBalance

    Shadowsocks
    ShadowsocksR
    Snell
    Socks5
    Http
    Vmess
    Vless
    Trojan
    Hysteria
    Hysteria2
    WireGuard
    Tuic
    Ssh
    Mieru
    AnyTLS
    Sudoku
    ECHTunnel  // ← 添加这一行
)
```

**位置 2**: 第 214 行后添加 (在 `case Sudoku:` 后面)

```go
func (at AdapterType) String() string {
    switch at {
    case Direct:
        return "Direct"
    // ... 其他 case ...
    case Sudoku:
        return "Sudoku"
    case ECHTunnel:  // ← 添加这个 case
        return "ECHTunnel"
    case Relay:
        return "Relay"
    // ... 其他 case ...
    }
}
```

---

### 2. `adapter/parser.go`

**位置**: 第 161 行前添加 (在 `case "sudoku":` 后面,`default:` 前面)

```go
func ParseProxy(mapping map[string]any, options ...ProxyOption) (C.Proxy, error) {
    // ... 现有代码 ...
    
    switch proxyType {
    case "ss":
        // ...
    // ... 其他 case ...
    case "sudoku":
        sudokuOption := &outbound.SudokuOption{BasicOption: basicOption}
        err = decoder.Decode(mapping, sudokuOption)
        if err != nil {
            break
        }
        proxy, err = outbound.NewSudoku(*sudokuOption)
    case "echtunnel":  // ← 添加这个 case
        echTunnelOption := &outbound.ECHTunnelOption{BasicOption: basicOption}
        err = decoder.Decode(mapping, echTunnelOption)
        if err != nil {
            break
        }
        proxy, err = outbound.NewECHTunnel(*echTunnelOption)
    default:
        return nil, fmt.Errorf("unsupport proxy type: %s", proxyType)
    }
    
    // ... 其他代码 ...
}
```

---

### 3. `go.mod` (如果需要)

确保有 gorilla/websocket 依赖:

```go
require (
    // ... 其他依赖 ...
    github.com/gorilla/websocket v1.5.0
)
```

如果没有,运行:
```bash
go get github.com/gorilla/websocket@v1.5.0
go mod tidy
```

---

## 🔧 编译步骤

### 1. 检查依赖

```bash
cd e:\Download\新建文件夹\mihomo-Alpha
go mod tidy
```

### 2. 编译

```bash
# Windows
go build -o mihomo.exe

# 或使用 Makefile
make windows
```

### 3. 验证

```bash
# 检查版本
.\mihomo.exe -v

# 测试配置
.\mihomo.exe -t -f config.yaml
```

---

## 📝 配置示例

创建 `config.yaml`:

```yaml
# Mihomo 配置文件

port: 7890
socks-port: 7891
allow-lan: false
mode: rule
log-level: info

proxies:
  - name: "ECH-Tunnel-VPS"
    type: echtunnel
    server: your-vps.com
    port: 443
    ws-path: /tunnel
    token: your-secret-token
    ech-domain: cloudflare-ech.com
    dns: https://dns.alidns.com/dns-query
    udp: true

proxy-groups:
  - name: "🚀 代理"
    type: select
    proxies:
      - ECH-Tunnel-VPS
      - DIRECT

  - name: "🎯 直连"
    type: select
    proxies:
      - DIRECT
      - 🚀 代理

rules:
  # 国内直连
  - GEOIP,CN,🎯 直连
  - GEOSITE,CN,🎯 直连
  
  # 其他走代理
  - MATCH,🚀 代理
```

---

## ✅ 验证清单

### 编译前检查

- [ ] `constant/adapters.go` 已添加 `ECHTunnel` 类型
- [ ] `constant/adapters.go` 的 `String()` 方法已添加对应 case
- [ ] `adapter/parser.go` 已添加 `echtunnel` case
- [ ] 所有新文件都已创建
- [ ] `go.mod` 包含 gorilla/websocket 依赖

### 编译检查

- [ ] `go mod tidy` 运行成功
- [ ] `go build` 无错误
- [ ] 生成的 `mihomo.exe` 文件存在

### 功能检查

- [ ] 配置文件解析成功 (`mihomo.exe -t -f config.yaml`)
- [ ] 程序启动成功
- [ ] ECH-Tunnel 连接建立成功
- [ ] 流量正常转发
- [ ] 规则匹配正常

---

## 🐛 常见问题

### Q: 编译报错 "undefined: C.ECHTunnel"

**A**: 检查 `constant/adapters.go` 是否正确添加了 `ECHTunnel` 常量

### Q: 编译报错 "undefined: outbound.NewECHTunnel"

**A**: 检查 `adapter/outbound/echtunnel.go` 文件是否存在且正确

### Q: 编译报错 "package echtunnel is not in GOROOT"

**A**: 检查 `transport/echtunnel/` 目录是否存在且包含 `client.go` 和 `conn.go`

### Q: 配置文件解析失败

**A**: 检查 `adapter/parser.go` 是否正确添加了 `echtunnel` case

### Q: 连接失败

**A**: 
1. 检查服务器地址是否正确
2. 检查 Token 是否正确
3. 查看日志输出

---

## 📚 参考文档

- [实现计划](mihomo-ech-integration.md)
- [Mihomo 官方文档](https://wiki.metacubex.one/)
- [ECH-Tunnel 项目](https://github.com/...)

---

## 🎉 完成后

编译成功后,你将拥有一个集成了 ECH-Tunnel 支持的 Mihomo!

**优势**:
- ✅ 一个程序搞定所有
- ✅ 统一的 YAML 配置
- ✅ 强大的规则引擎
- ✅ 完整的生态支持

**使用方式**:
```bash
# 启动 Mihomo
.\mihomo.exe -f config.yaml

# 设置系统代理为 127.0.0.1:7890
# 享受智能分流!
```
