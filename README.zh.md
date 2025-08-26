English | [中文](./README.zh.md)

[![Go Version](https://img.shields.io/badge/go-1.22-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-AGPL--3.0-green.svg)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/xuemian168/domain-scanner.svg?style=social)](https://github.com/xuemian168/domain-scanner/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/xuemian168/domain-scanner.svg?style=social)](https://github.com/xuemian168/domain-scanner/network/members)
[![GitHub Issues](https://img.shields.io/github/issues/xuemian168/domain-scanner.svg)](https://github.com/xuemian168/domain-scanner/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/xuemian168/domain-scanner.svg)](https://github.com/xuemian168/domain-scanner/pulls)

# 域名扫描器

一个强大的域名可用性检查工具，使用 Go 语言编写。该工具通过多种注册指标检查来帮助您找到可用的域名，并提供详细的验证结果。

### 网页版: [zli.li](https://zli.li)

![Star History Chart](https://api.star-history.com/svg?repos=xuemian168/domain-scanner&type=Date)

![screenshot](./imgs/image.png)

## 功能特点

- **多方法验证**：使用多种方法检查域名可用性：
  - DNS 记录（NS、A、MX）
  - WHOIS 信息
  - SSL 证书验证
- **高级过滤**：使用正则表达式过滤域名
- **并发处理**：可配置工作线程数的多线程域名检查
- **智能错误处理**：自动重试机制处理失败的查询
- **详细结果**：显示已注册域名的验证签名
- **进度跟踪**：实时显示当前/总数进度
- **文件输出**：将结果分别保存到可用和已注册域名的文件中
- **可配置延迟**：可调整查询间隔以防止速率限制

## 安装

```bash
git clone https://github.com/xuemian168/domain-scanner.git
cd domain-scanner
go mod download
```

## 使用方法

```bash
go run main.go [选项]
```

### 选项

- `-l int`: 域名长度（默认：3）
- `-s string`: 域名后缀（默认：.li）
- `-p string`: 域名模式：
  - `d`: 纯数字（例如：123.li）
  - `D`: 纯字母（例如：abc.li）
  - `a`: 字母数字组合（例如：a1b.li）
- `-delay int`: 查询间隔（毫秒）（默认：1000）
- `-workers int`: 并发工作线程数（默认：10）
- `-show-registered`: 在输出中显示已注册的域名（默认：false）
- `-h`: 显示帮助信息
- `-r string`: 域名正则表达式过滤器
- `-regex-mode string`: 正则表达式匹配模式（默认：full）
  - `full`：匹配完整域名
  - `prefix`：仅匹配域名前缀

### 示例

1. 使用 20 个工作线程检查 3 字母 .li 域名：
```bash
go run main.go -l 3 -s .li -p D -workers 20
```

2. 使用自定义延迟和工作线程数检查域名：
```bash
go run main.go -l 3 -s .li -p D -delay 500 -workers 15
```

3. 显示可用和已注册的域名：
```bash
go run main.go -l 3 -s .li -p D -show-registered
```

4. 使用完整域名正则表达式过滤：
```bash
go run main.go -l 3 -s .li -p D -r "^[a-z]{2}[0-9]$" -regex-mode full
```

5. 使用域名前缀正则表达式过滤：
```bash
go run main.go -l 3 -s .li -p D -r "^[a-z]{2}" -regex-mode prefix
```

## 输出格式

### 进度显示
```
[1/100] Domain abc.com AVAILABLE!
[2/100] Domain xyz.com REGISTERED [DNS_NS, WHOIS]
```

### 验证签名说明
- `DNS_NS`：域名有名称服务器记录
- `DNS_A`：域名有 IP 地址记录
- `DNS_MX`：域名有邮件服务器记录
- `WHOIS`：根据 WHOIS 信息域名已注册
- `SSL`：域名有有效的 SSL 证书

### 输出文件
- 可用域名：`available_domains_[模式]_[长度]_[后缀].txt`
- 已注册域名：`registered_domains_[模式]_[长度]_[后缀].txt`

## 错误处理

工具包含强大的错误处理机制：
- WHOIS 查询自动重试机制（3次尝试）
- SSL 证书检查超时设置
- 优雅处理网络问题
- 详细的错误报告

## 贡献

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](./CONTRIBUTING.md)

我们欢迎社区贡献！无论您是修复错误、添加新功能、改进文档还是报告问题，我们都非常感谢您的帮助。

### 如何贡献

1. **Fork 仓库**：创建项目的副本
2. **创建功能分支**：在专用分支中进行更改
3. **进行修改**：遵循编码规范并彻底测试
4. **提交 Pull Request**：描述您的更改并链接相关问题

有关详细的贡献指南、开发环境设置和编码标准，请阅读我们的 [CONTRIBUTING.md](./CONTRIBUTING.md) 文件。

### 贡献者快速开始

```bash
# Fork 并克隆仓库
git clone https://github.com/YOUR_USERNAME/domain-scanner.git
cd domain-scanner

# 设置开发环境
go mod download
go build -o domain-scanner main.go

# 创建功能分支
git checkout -b feature/your-feature-name

# 进行更改并测试
go run main.go -l 2 -s .test -p D

# 提交并推送
git commit -m "feat: 您的功能描述"
git push origin feature/your-feature-name
```

## 许可证

本项目采用 AGPL-3.0 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。

## 最近更新

### v1.3.2 - 2025-08-26
- **安全**：添加 ReDoS 攻击防护，正则表达式超时机制（100毫秒）
- **安全**：实现正则表达式复杂度验证和危险模式检测
- **性能**：恢复内存高效的流式架构
- **增强**：升级到 regexp2，支持高级正则功能（回溯引用、环视）
- **增强**：添加全面的正则表达式安全指南和示例
- **稳定性**：改进正则匹配操作的错误处理

### v1.3.1 - 2025-08-24
- **新增**：多 WHOIS 服务器支持，提升可靠性
- **新增**：指数退避重试机制，优化 WHOIS 查询  
- **新增**：全面的保留域名标识符（139 个模式）
- **性能**：误报率降低 67%（15% → 5%）
- **性能**：WHOIS 查询成功率提升 23%（~75% → ~92%）

### v1.3.0
- **性能优化**：显著提升域名检查速度
- **错误修复**：修复 .de 域名和其他 TLD 的 WHOIS 解析问题
- **代码质量**：重构内部架构，提高可维护性

📋 **[查看完整更新日志](docs/CHANGELOG.md)** - 查看详细版本历史、技术改进和所有变更。