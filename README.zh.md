[English](./README.md) | 中文

# 域名扫描器

一个用Go语言编写的强大且灵活的域名可用性检查工具。该工具可以帮助你根据各种模式和过滤器查找可用的域名。

## 功能特点

- 使用DNS和WHOIS查询检查域名可用性
- 基于不同模式生成域名：
  - 纯数字（例如：123.li）
  - 纯字母（例如：abc.li）
  - 字母数字组合（例如：a1b.li）
- 使用正则表达式进行高级筛选
- 将可用和已注册的域名分别保存到不同文件
- 可配置查询间隔以避免请求限制

## 安装

```bash
git clone https://github.com/yourusername/domain-scanner.git
cd domain-scanner
go mod tidy
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
- `-r string`: 域名正则表达式过滤器
- `-delay int`: 查询间隔（毫秒）（默认：1000）
- `-h`: 显示帮助信息

### 示例

1. 检查3字母的.li域名：
```bash
go run main.go -l 3 -s .li -p D
```

2. 检查3位数字的.li域名：
```bash
go run main.go -l 3 -s .li -p d
```

3. 检查包含'abc'的域名：
```bash
go run main.go -l 5 -s .li -p D -r '.*abc.*'
```

4. 检查以'a'开头，以'z'结尾的域名：
```bash
go run main.go -l 4 -s .li -p D -r '^a.*z$'
```

5. 检查只包含元音的域名：
```bash
go run main.go -l 3 -s .li -p D -r '^[aeiou]+$'
```

## 输出结果

程序会生成两个输出文件：
- `available_domains_[pattern]_[length]_[suffix].txt`（可用域名）
- `registered_domains_[pattern]_[length]_[suffix].txt`（已注册域名）

## 正则表达式示例

1. `^a.*` - 以'a'开头的域名
2. `.*z$` - 以'z'结尾的域名
3. `^[0-9]+$` - 纯数字域名
4. `^[a-z]+$` - 纯字母域名
5. `^[a-z][0-9][a-z]$` - 字母-数字-字母格式的域名
6. `.*(com|net|org)$` - 包含特定后缀的域名
7. `^[a-z]{2}\d{2}$` - 两个字母后跟两个数字的域名

## 贡献

欢迎贡献代码！请随时提交Pull Request。

## 许可证

本项目采用MIT许可证 - 详见LICENSE文件。 