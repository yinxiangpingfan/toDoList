# 微服务日志聚合方案

## 📋 问题分析

### 当前问题
在微服务架构中，每个服务都在不同的服务器/容器中运行，日志分散在各处：
```
User服务 → /logs/user.log (服务器A)
Task服务 → /logs/task.log (服务器B)
Gateway → /logs/gateway.log (服务器C)
```

**问题：**
- ❌ 日志分散，难以查询
- ❌ 无法追踪跨服务的请求链路
- ❌ 排查问题需要登录多台服务器
- ❌ 日志容易丢失（服务器故障）

### 错误方案
❌ **把所有日志写到一个文件** - 这在微服务中不可行，因为：
- 服务在不同机器上，无法共享文件
- 即使用NFS共享，会有性能和并发问题
- 违背微服务独立部署的原则

## ✅ 正确方案：集中式日志管理

### 方案对比

| 方案 | 复杂度 | 成本 | 适用场景 |
|------|--------|------|---------|
| **ELK Stack** | ⭐⭐⭐⭐ | 高 | 大型企业 |
| **Loki + Grafana** | ⭐⭐⭐ | 中 | 中大型项目（推荐） |
| **Fluentd + ES** | ⭐⭐⭐⭐ | 高 | 大型项目 |
| **云服务日志** | ⭐⭐ | 按量付费 | 云上部署 |
| **简单方案：文件+rsync** | ⭐ | 低 | 小型项目 |

## 🎯 推荐方案：Loki + Grafana

### 架构图

```
┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│ User服务     │ ──────► │   Promtail  │ ──────► │    Loki     │
│ (输出日志)   │  文件   │ (日志采集器) │  推送   │ (日志存储)   │
└─────────────┘         └─────────────┘         └─────────────┘
                                                        │
┌─────────────┐         ┌─────────────┐                │
│ Task服务     │ ──────► │   Promtail  │ ───────────────┤
│ (输出日志)   │  文件   │ (日志采集器) │  推送          │
└─────────────┘         └─────────────┘                │
                                                        │
┌─────────────┐         ┌─────────────┐                │
│ Gateway服务  │ ──────► │   Promtail  │ ───────────────┤
│ (输出日志)   │  文件   │ (日志采集器) │  推送          │
└─────────────┘         └─────────────┘                │
                                                        ▼
                                                ┌─────────────┐
                                                │   Grafana   │
                                                │ (日志查询UI) │
                                                └─────────────┘
```

### 组件说明

1. **各微服务** - 继续输出日志到本地文件（不需要改代码）
2. **Promtail** - 轻量级日志采集器，监控日志文件并推送到Loki
3. **Loki** - 日志存储和索引系统（类似Prometheus，但用于日志）
4. **Grafana** - 统一的日志查询和可视化界面

### 优势

- ✅ **轻量级** - Loki比Elasticsearch轻量很多
- ✅ **易部署** - Docker Compose一键部署
- ✅ **低成本** - 资源占用少
- ✅ **强大查询** - 支持LogQL查询语言
- ✅ **无需改代码** - 服务继续输出到文件即可

## 🚀 实施步骤

### 第1步：修改日志格式（JSON格式）

**为什么要JSON格式？**
- 便于解析和查询
- 支持结构化字段
- 便于添加元数据（服务名、traceID等）

**修改：`global/logger/log.go`**

```go
package logger

import (
    "os"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func LoggerInit(logPath string, serviceName string) {
    // 配置编码器（JSON格式）
    encoderConfig := zapcore.EncoderConfig{
        TimeKey:        "timestamp",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        MessageKey:     "message",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.LowercaseLevelEncoder,
        EncodeTime:     zapcore.ISO8601TimeEncoder,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    }

    // 创建文件输出
    file, _ := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    fileWriter := zapcore.AddSync(file)

    // 创建核心
    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(encoderConfig),  // JSON编码器
        fileWriter,
        zapcore.DebugLevel,
    )

    // 创建logger
    logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
    
    // 添加全局字段（服务名）
    Logger = logger.Sugar().With("service", serviceName)
}
```

**修改各服务的初始化：**

```go
// app/user/cmd/main.go
logger.LoggerInit("../../../logs/user.log", "user-service")

// app/task/cmd/main.go
logger.LoggerInit("../../../logs/task.log", "task-service")

// app/gateway/cmd/main.go
logger.LoggerInit("../../../logs/gateway.log", "gateway-service")
```

**日志输出示例：**
```json
{
  "timestamp": "2024-01-15T10:30:45.123Z",
  "level": "info",
  "service": "user-service",
  "caller": "service/login.go:25",
  "message": "用户登录成功",
  "user_id": 123,
  "username": "alice"
}
```

### 第2步：添加TraceID（可选但推荐）

**为什么需要TraceID？**
- 追踪跨服务的请求链路
- 快速定位问题
- 关联所有相关日志

**创建：`global/utils/trace.go`**

```go
package utils

import (
    "context"
    "github.com/google/uuid"
)

type contextKey string

const TraceIDKey contextKey = "trace_id"

// GenerateTraceID 生成TraceID
func GenerateTraceID() string {
    return uuid.New().String()
}

// WithTraceID 将TraceID添加到context
func WithTraceID(ctx context.Context, traceID string) context.Context {
    return context.WithValue(ctx, TraceIDKey, traceID)
}

// GetTraceID 从context获取TraceID
func GetTraceID(ctx context.Context) string {
    if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
        return traceID
    }
    return ""
}
```

**在服务中使用：**

```go
// 在gRPC拦截器中添加TraceID
func (u *UserSrv) Login(ctx context.Context, req *pb.LoginRequest, res *pb.LoginResponse) error {
    // 生成或获取TraceID
    traceID := utils.GetTraceID(ctx)
    if traceID == "" {
        traceID = utils.GenerateTraceID()
        ctx = utils.WithTraceID(ctx, traceID)
    }
    
    // 记录日志时带上TraceID
    logger.Logger.With("trace_id", traceID).Infof("用户登录请求: username=%s", req.Username)
    
    // ... 业务逻辑
}
```

### 第3步：部署Loki + Grafana

**创建：`docker-compose-logging.yml`**

```yaml
version: "3"

services:
  # Loki - 日志存储
  loki:
    image: grafana/loki:2.9.0
    ports:
      - "3100:3100"
    command: -config.file=/etc/loki/local-config.yaml
    volumes:
      - ./loki-config.yaml:/etc/loki/local-config.yaml
      - loki-data:/loki
    networks:
      - logging

  # Promtail - 日志采集器（User服务）
  promtail-user:
    image: grafana/promtail:2.9.0
    volumes:
      - ./logs:/logs
      - ./promtail-config.yaml:/etc/promtail/config.yaml
    command: -config.file=/etc/promtail/config.yaml
    environment:
      - SERVICE_NAME=user-service
    networks:
      - logging
    depends_on:
      - loki

  # Promtail - 日志采集器（Task服务）
  promtail-task:
    image: grafana/promtail:2.9.0
    volumes:
      - ./logs:/logs
      - ./promtail-config.yaml:/etc/promtail/config.yaml
    command: -config.file=/etc/promtail/config.yaml
    environment:
      - SERVICE_NAME=task-service
    networks:
      - logging
    depends_on:
      - loki

  # Grafana - 日志查询UI
  grafana:
    image: grafana/grafana:10.0.0
    ports:
      - "3000:3000"
    environment:
      - GF_AUTH_ANONYMOUS_ENABLED=true
      - GF_AUTH_ANONYMOUS_ORG_ROLE=Admin
    volumes:
      - grafana-data:/var/lib/grafana
    networks:
      - logging
    depends_on:
      - loki

volumes:
  loki-data:
  grafana-data:

networks:
  logging:
    driver: bridge
```

**创建：`loki-config.yaml`**

```yaml
auth_enabled: false

server:
  http_listen_port: 3100

ingester:
  lifecycler:
    address: 127.0.0.1
    ring:
      kvstore:
        store: inmemory
      replication_factor: 1
  chunk_idle_period: 5m
  chunk_retain_period: 30s

schema_config:
  configs:
    - from: 2020-10-24
      store: boltdb-shipper
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h

storage_config:
  boltdb_shipper:
    active_index_directory: /loki/index
    cache_location: /loki/cache
    shared_store: filesystem
  filesystem:
    directory: /loki/chunks

limits_config:
  enforce_metric_name: false
  reject_old_samples: true
  reject_old_samples_max_age: 168h

chunk_store_config:
  max_look_back_period: 0s

table_manager:
  retention_deletes_enabled: false
  retention_period: 0s
```

**创建：`promtail-config.yaml`**

```yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://loki:3100/loki/api/v1/push

scrape_configs:
  # User服务日志
  - job_name: user-service
    static_configs:
      - targets:
          - localhost
        labels:
          job: user-service
          __path__: /logs/user.log

  # Task服务日志
  - job_name: task-service
    static_configs:
      - targets:
          - localhost
        labels:
          job: task-service
          __path__: /logs/task.log

  # Gateway日志
  - job_name: gateway-service
    static_configs:
      - targets:
          - localhost
        labels:
          job: gateway-service
          __path__: /logs/gateway.log
```

### 第4步：启动日志系统

```bash
# 启动Loki + Grafana + Promtail
docker-compose -f docker-compose-logging.yml up -d

# 查看状态
docker-compose -f docker-compose-logging.yml ps

# 查看日志
docker-compose -f docker-compose-logging.yml logs -f
```

### 第5步：在Grafana中查询日志

**1. 访问Grafana**
```
http://localhost:3000
```

**2. 添加Loki数据源**
- 点击 Configuration → Data Sources
- 添加 Loki
- URL: `http://loki:3100`
- 保存

**3. 查询日志**

**查询所有日志：**
```logql
{job=~".+"}
```

**查询特定服务：**
```logql
{job="user-service"}
```

**查询特定级别：**
```logql
{job="user-service"} |= "error"
```

**查询特定TraceID：**
```logql
{job=~".+"} | json | trace_id="550e8400-e29b-41d4-a716-446655440000"
```

**查询特定用户的操作：**
```logql
{job=~".+"} | json | user_id="123"
```

**统计错误数量：**
```logql
sum(count_over_time({job=~".+"} |= "error" [5m]))
```

## 📊 Grafana Dashboard示例

### 创建日志监控面板

**1. 错误日志统计**
```logql
sum by (service) (count_over_time({job=~".+"} |= "error" [1m]))
```

**2. 请求量统计**
```logql
sum by (service) (count_over_time({job=~".+"} |= "请求" [1m]))
```

**3. 响应时间分布**
```logql
{job=~".+"} | json | __error__="" | unwrap duration | quantile_over_time(0.95, [5m])
```

## 🔍 日志查询实战

### 场景1：追踪一次完整的用户请求

```logql
# 通过TraceID查询所有相关日志
{job=~".+"} | json | trace_id="abc-123-def"
```

**结果示例：**
```
[Gateway] 收到登录请求 trace_id=abc-123-def
[User服务] 验证用户名密码 trace_id=abc-123-def user_id=123
[User服务] 登录成功 trace_id=abc-123-def user_id=123
[Gateway] 返回登录结果 trace_id=abc-123-def
```

### 场景2：查找某个用户的所有操作

```logql
{job=~".+"} | json | user_id="123"
```

### 场景3：查找最近的错误

```logql
{job=~".+"} |= "error" | json
```

### 场景4：性能分析

```logql
# 查找响应时间超过1秒的请求
{job=~".+"} | json | duration > 1000
```

## 🎯 简化方案（小型项目）

如果觉得Loki太复杂，可以用更简单的方案：

### 方案1：rsync定时同步

**原理：**
- 各服务输出日志到本地文件
- 使用rsync定时同步到中心服务器
- 在中心服务器上查看所有日志

**实现：**

```bash
# 在中心服务器上创建脚本：sync-logs.sh
#!/bin/bash

# 同步User服务日志
rsync -avz user@server1:/app/logs/user.log /logs/user.log

# 同步Task服务日志
rsync -avz task@server2:/app/logs/task.log /logs/task.log

# 同步Gateway日志
rsync -avz gateway@server3:/app/logs/gateway.log /logs/gateway.log
```

**添加定时任务：**
```bash
# 每分钟同步一次
* * * * * /path/to/sync-logs.sh
```

**查看日志：**
```bash
# 查看所有日志
tail -f /logs/*.log

# 搜索特定内容
grep "error" /logs/*.log

# 按时间排序
cat /logs/*.log | sort
```

### 方案2：云服务日志（推荐云上部署）

**阿里云日志服务（SLS）：**
- 自动采集日志
- 强大的查询功能
- 按量付费

**AWS CloudWatch Logs：**
- 与AWS服务集成
- 实时监控告警

**腾讯云日志服务（CLS）：**
- 类似阿里云SLS

## 📈 最佳实践

### 1. 日志级别规范

```go
// Debug - 调试信息（开发环境）
logger.Logger.Debugf("查询参数: %+v", params)

// Info - 正常业务流程
logger.Logger.Infof("用户登录成功: user_id=%d", userID)

// Warn - 警告（不影响业务）
logger.Logger.Warnf("Redis连接失败，使用降级方案")

// Error - 错误（影响业务）
logger.Logger.Errorf("创建任务失败: %v", err)

// Panic - 严重错误（服务崩溃）
logger.Logger.Panicf("数据库连接失败: %v", err)
```

### 2. 结构化日志

```go
// ❌ 不好的日志
logger.Logger.Infof("用户123登录成功")

// ✅ 好的日志
logger.Logger.With(
    "user_id", 123,
    "username", "alice",
    "ip", "192.168.1.1",
    "trace_id", traceID,
).Info("用户登录成功")
```

### 3. 敏感信息脱敏

```go
// ❌ 不要记录密码
logger.Logger.Infof("登录: username=%s, password=%s", username, password)

// ✅ 脱敏处理
logger.Logger.Infof("登录: username=%s, password=***", username)
```

### 4. 日志轮转

```go
// 使用lumberjack进行日志轮转
import "gopkg.in/natefinch/lumberjack.v2"

func LoggerInit(logPath string, serviceName string) {
    writer := &lumberjack.Logger{
        Filename:   logPath,
        MaxSize:    100,  // MB
        MaxBackups: 3,    // 保留3个备份
        MaxAge:     28,   // 天
        Compress:   true, // 压缩
    }
    
    // ... 使用writer
}
```

## 🆚 方案对比总结

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|---------|
| **Loki + Grafana** | 轻量、易用、强大 | 需要学习LogQL | 中大型项目（推荐） |
| **ELK Stack** | 功能最强大 | 重量级、资源消耗大 | 大型企业 |
| **rsync同步** | 最简单 | 不实时、功能弱 | 小型项目 |
| **云服务日志** | 免运维 | 按量付费 | 云上部署 |

## ✅ 推荐方案

### 小型项目（<5个服务）
**rsync + grep** - 简单够用

### 中型项目（5-20个服务）
**Loki + Grafana** - 最佳选择（推荐）

### 大型项目（>20个服务）
**ELK Stack** 或 **云服务日志**

## 📚 参考资料

- [Grafana Loki官方文档](https://grafana.com/docs/loki/latest/)
- [LogQL查询语言](https://grafana.com/docs/loki/latest/logql/)
- [Promtail配置](https://grafana.com/docs/loki/latest/clients/promtail/)
- [Zap日志库](https://github.com/uber-go/zap)

## 🎓 总结

**微服务日志管理的核心原则：**
1. ✅ 每个服务独立输出日志到本地文件
2. ✅ 使用日志采集器统一收集
3. ✅ 集中存储到日志系统
4. ✅ 通过统一界面查询和分析
5. ✅ 使用TraceID追踪请求链路

**不要：**
- ❌ 把所有日志写到一个文件
- ❌ 通过NFS共享日志文件
- ❌ 手动登录各服务器查看日志

对于你的ToDoList项目，**强烈推荐使用Loki + Grafana方案**！
