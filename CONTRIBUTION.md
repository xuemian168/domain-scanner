# Contributing to Domain Scanner / 为域名扫描器贡献代码

[English](#english) | [中文](#中文)

---

## English

Thank you for your interest in contributing to Domain Scanner! This document provides guidelines for contributing to the project.

### Getting Started

#### 1. Fork the Repository
- Visit the [Domain Scanner repository](https://github.com/yourusername/domain_scanner)
- Click the "Fork" button in the top-right corner
- This creates a copy of the repository in your GitHub account

#### 2. Clone Your Fork
```bash
git clone https://github.com/YOUR_USERNAME/domain_scanner.git
cd domain_scanner
```

#### 3. Add Upstream Remote
```bash
git remote add upstream https://github.com/ORIGINAL_OWNER/domain_scanner.git
```

#### 4. Set Up Development Environment
```bash
# Download dependencies
go mod download

# Test the build
go build -o domain-scanner main.go

# Run tests (if you're adding functionality)
go run main.go -l 2 -s .test -p D
```

### Making Changes

#### 1. Create a New Branch
```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-description
```

#### 2. Make Your Changes
- Follow the existing code style and patterns
- Keep changes focused and atomic
- Test your changes thoroughly

#### 3. Commit Your Changes
```bash
git add .
git commit -m "feat: add new feature description"
# or
git commit -m "fix: resolve issue with domain checking"
```

#### Commit Message Guidelines
- Use conventional commit format: `type: description`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- Keep the first line under 72 characters
- Use present tense ("add feature" not "added feature")

#### 4. Push to Your Fork
```bash
git push origin feature/your-feature-name
```

#### 5. Create a Pull Request
- Go to your fork on GitHub
- Click "New Pull Request"
- Select your branch and provide a clear description
- Include:
  - What changes you made
  - Why you made them
  - How to test the changes
  - Any related issues

### Code Guidelines

#### Go Code Style
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions focused and small
- Handle errors appropriately

#### Testing
- Test your changes with various domain patterns
- Ensure the tool works with different suffixes and lengths
- Verify concurrency handling works correctly

### Pull Request Process

1. Ensure your PR has a clear title and description
2. Link any related issues
3. Make sure all checks pass
4. Be responsive to feedback and suggestions
5. Once approved, your PR will be merged

### Reporting Issues

- Use the GitHub issue tracker
- Provide clear reproduction steps
- Include your Go version and operating system
- Attach relevant logs or error messages

---

## 中文

感谢您对域名扫描器项目的关注！本文档为您提供了参与项目贡献的指南。

### 开始贡献

#### 1. Fork 仓库
- 访问[域名扫描器仓库](https://github.com/yourusername/domain_scanner)
- 点击右上角的 "Fork" 按钮
- 这将在您的 GitHub 账户中创建仓库的副本

#### 2. 克隆您的 Fork
```bash
git clone https://github.com/YOUR_USERNAME/domain_scanner.git
cd domain_scanner
```

#### 3. 添加上游远程仓库
```bash
git remote add upstream https://github.com/ORIGINAL_OWNER/domain_scanner.git
```

#### 4. 设置开发环境
```bash
# 下载依赖
go mod download

# 测试构建
go build -o domain-scanner main.go

# 运行测试（如果您添加了新功能）
go run main.go -l 2 -s .test -p D
```

### 进行更改

#### 1. 创建新分支
```bash
git checkout -b feature/your-feature-name
# 或
git checkout -b fix/issue-description
```

#### 2. 进行修改
- 遵循现有的代码风格和模式
- 保持更改的专注性和原子性
- 彻底测试您的更改

#### 3. 提交更改
```bash
git add .
git commit -m "feat: 添加新功能描述"
# 或
git commit -m "fix: 解决域名检查问题"
```

#### 提交消息指南
- 使用约定式提交格式：`type: description`
- 类型：`feat`、`fix`、`docs`、`style`、`refactor`、`test`、`chore`
- 第一行保持在 72 个字符以内
- 使用现在时态（"add feature" 而不是 "added feature"）

#### 4. 推送到您的 Fork
```bash
git push origin feature/your-feature-name
```

#### 5. 创建 Pull Request
- 转到 GitHub 上您的 fork
- 点击 "New Pull Request"
- 选择您的分支并提供清晰的描述
- 包括：
  - 您做了什么更改
  - 为什么进行这些更改
  - 如何测试更改
  - 任何相关问题

### 代码指南

#### Go 代码风格
- 遵循标准的 Go 格式化（`gofmt`）
- 使用有意义的变量和函数名
- 为复杂逻辑添加注释
- 保持函数专注和简小
- 适当处理错误

#### 测试
- 使用各种域名模式测试您的更改
- 确保工具适用于不同的后缀和长度
- 验证并发处理正常工作

### Pull Request 流程

1. 确保您的 PR 有清晰的标题和描述
2. 链接任何相关问题
3. 确保所有检查通过
4. 对反馈和建议做出响应
5. 一旦获得批准，您的 PR 将被合并

### 报告问题

- 使用 GitHub 问题跟踪器
- 提供清晰的复现步骤
- 包括您的 Go 版本和操作系统
- 附上相关日志或错误消息

---

## Development Tips / 开发提示

### Local Testing / 本地测试
```bash
# Quick functionality test
go run main.go -l 2 -s .li -p D -workers 5 -delay 100

# Test with regex filtering
go run main.go -l 3 -s .li -p D -r "^[a-z]{2}[0-9]$" -regex-mode full

# Performance test with minimal delay
go run main.go -l 2 -s .test -p D -workers 20 -delay 50
```

### Building for Release / 构建发布版本
```bash
# Build with optimization flags
go build -ldflags="-s -w" -o domain-scanner main.go

# Test release configuration
goreleaser check
goreleaser build --snapshot --clean
```

### Common Issues / 常见问题

- **Build Errors**: Ensure Go 1.19+ is installed
- **Network Issues**: Test with longer delays if rate-limited
- **Performance**: Adjust worker count based on your system capabilities

构建错误：确保安装了 Go 1.19+ 版本
网络问题：如果遇到速率限制，请测试更长的延迟
性能：根据系统能力调整工作线程数量

---

Thank you for contributing! / 感谢您的贡献！