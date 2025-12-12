# Mihomo-ECH 更新日志

## [1.0.0] - 2025-12-12

### 新增功能

- ✨ 添加 ECH-Tunnel 协议支持
- ✨ 添加 ECH-Workers 协议支持
- ✨ 新增 `echtunnel` 代理类型
- ✨ 支持 WebSocket + ECH 传输
- ✨ 支持 Token 认证
- ✨ 支持自定义 ECH 域名
- ✨ 支持 DoH DNS 配置
- ✨ 支持 UDP 转发

### 新增文件

- `adapter/outbound/echtunnel.go` - ECH-Tunnel Adapter 实现
- `transport/echtunnel/client.go` - ECH-Tunnel 客户端
- `transport/echtunnel/conn.go` - WebSocket 连接包装
- `ECH-TUNNEL-INTEGRATION-GUIDE.md` - 集成指南
- `CHANGELOG.md` - 更新日志

### 修改文件

- `constant/adapters.go` - 添加 ECHTunnel 类型定义
- `adapter/parser.go` - 添加 echtunnel 协议解析
- `README.md` - 更新项目说明

### 技术细节

- 基于 Mihomo (Clash Meta) 开发
- 使用 gorilla/websocket 实现 WebSocket 连接
- 完全兼容 ECH-Tunnel 和 ECH-Workers 服务端
- 支持与其他协议混用

### 配置示例

```yaml
proxies:
  - name: "ECH-Tunnel"
    type: echtunnel
    server: example.com
    port: 443
    token: your-token
```

---

## 计划中的功能

### [1.1.0] - 待定

- [ ] ECH 配置自动获取
- [ ] 连接池优化
- [ ] 多路复用支持
- [ ] 更详细的日志输出
- [ ] 性能优化

### [1.2.0] - 待定

- [ ] GUI 配置工具
- [ ] 健康检查增强
- [ ] 更多传输选项

---

## 贡献指南

欢迎提交 Issue 和 Pull Request!

如果你发现 bug 或有新功能建议,请:
1. 提交 Issue 描述问题
2. Fork 项目并创建新分支
3. 提交 Pull Request

---

**感谢所有贡献者!** ❤️
