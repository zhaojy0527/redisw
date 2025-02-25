# Redisw

Redisw 是一个简单易用的 Redis 服务器连接管理工具，支持多服务器配置和快速切换。它能帮助开发者快速在多个 Redis 服务器之间进行切换，提高开发和运维效率。

## 功能特点

- 支持多个 Redis 服务器配置
- 交互式服务器选择界面
- 支持服务器名称搜索
- 支持密码保护的 Redis 连接
- 支持 Redis 集群连接
- 支持自定义配置文件路径
- 命令行界面简洁直观

## 安装要求

- Go 1.16 或更高版本
- Redis CLI 工具
- MacOS/Linux 系统

### 安装依赖

MacOS:
```bash
brew install redis
```

Linux (Ubuntu/Debian):
```bash
apt-get install redis-tools
```

## 安装

### 通过 Homebrew 安装

```bash
brew install redisw
```

### 从源码安装

```bash
git clone https://github.com/zhaojy0527/redisw.git
cd redisw
make build
```

## 配置

### 配置文件位置
默认配置文件位置：`~/.config/redisw/config.yml`

### 配置文件示例

```yaml
- name: "本地Redis"
  host: "localhost"
  port: 6379
  password: ""

- name: "开发环境"
  host: "dev.redis.example.com"
  port: 6379
  password: "your-password"

- name: "测试集群"
  host: "test.redis.cluster"
  port: 6379
  password: "cluster-password"
```

## 使用方法

### 基本命令

```bash
# 使用默认配置文件启动
redisw

# 指定配置文件启动
redisw -config /path/to/config.yml

# 显示版本信息
redisw -version
```

### 交互式界面操作

1. 启动后会显示已配置的 Redis 服务器列表
2. 使用上下箭头键选择要连接的服务器
3. 输入关键字可以快速搜索服务器
4. 按 Enter 键连接选中的服务器
5. 按 Ctrl+C 可以返回服务器选择界面
6. 再次按 Ctrl+C 可以退出程序

## 开发

### 获取源码

```bash
git clone https://github.com/zhaojy0527/redisw.git
cd redisw
```

### 构建

```bash
make build
```

### 运行测试

```bash
make test
```

### 清理构建文件

```bash
make clean
```

## 问题反馈

如果您在使用过程中遇到任何问题，或有任何建议，请：

1. 提交 [Issue](https://github.com/zhaojy0527/redisw/issues)
2. 发送邮件至：zhaojianyong0527@gmail.com

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件