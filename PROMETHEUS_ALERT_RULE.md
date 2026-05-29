# Prometheus 告警规则标准模板

## 一、告警规则示例

```yaml
groups:
- name: example-alerts
  rules:
  # 示例：实例宕机告警
  - alert: InstanceDown
    expr: up == 0
    for: 1m
    labels:
      severity: critical
      team: ops
      service: infrastructure
    annotations:
      summary: "实例 {{ $labels.instance }} 已宕机"
      description: "{{ $labels.instance }} 的 node_exporter 指标采集失败超过 1 分钟"
```

## 二、Labels 定义标准

| Label 名称 | 必填 | 说明 | 示例值 |
|-----------|------|------|--------|
| `severity` | ✅ | 告警级别 | critical / warning / info |
| `team` | ✅ | 负责团队 | ops / backend / frontend |
| `service` | ✅ | 服务/业务名称 | api / db / nginx |
| `instance` | ✅ | 告警实例（自动获取） | 172.16.11.236:9100 |

### severity 级别说明

| 级别 | 说明 | emoji |
|------|------|-------|
| critical | 严重，影响服务可用性 | 🔴 |
| warning | 警告，需要关注 | 🟡 |
| info | 通知，信息类 | 🔵 |

## 三、Annotations 定义标准

| Annotation 名称 | 必填 | 说明 |
|----------------|------|------|
| summary | ✅ | 告警简述（简短，一句话） |
| description | ✅ | 告警详情（详细描述问题） |

## 四、告警消息格式

最终微信消息格式如下：

```
━━━━━━━━━━━━━━━━━━
🚨 触发中 | 【服务名】告警简述
━━━━━━━━━━━━━━━━━━

📋 告警信息
├─ 🔴 级别：🔴 严重
├─ 📛 类型：告警规则名
├─ 🖥 主机：实例地址
├─ 📝 详情：告警详细描述
└─ ⏰ 时间：触发时间
```

## 五、建议的告警规则配置

```yaml
groups:
- name: infrastructure
  rules:
  # 服务宕机
  - alert: ServiceDown
    expr: up == 0
    for: 1m
    labels:
      severity: critical
      team: ops
      service: your-service-name
    annotations:
      summary: "服务 {{ $labels.instance }} 已宕机"
      description: "{{ $labels.instance }} 指标采集失败，服务可能已停止"

  # CPU 告警
  - alert: HighCPU
    expr: 100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
    for: 5m
    labels:
      severity: warning
      team: ops
      service: your-service-name
    annotations:
      summary: "【{{ $labels.service }}】CPU 使用率超过 80%"
      description: "主机 {{ $labels.instance }} CPU 使用率持续过高，当前值: {{ $value }}%"

  # 内存告警
  - alert: HighMemory
    expr: (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100 > 85
    for: 5m
    labels:
      severity: warning
      team: ops
      service: your-service-name
    annotations:
      summary: "【{{ $labels.service }}】内存使用率超过 85%"
      description: "主机 {{ $labels.instance }} 内存使用率过高，当前值: {{ $value }}%"

  # 磁盘告警
  - alert: DiskFull
    expr: (node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}) * 100 < 15
    for: 5m
    labels:
      severity: critical
      team: ops
      service: your-service-name
    annotations:
      summary: "【{{ $labels.service }}】磁盘空间不足"
      description: "主机 {{ $labels.instance }} 磁盘空间不足，当前剩余: {{ $value }}%"
```

## 六、消息示例

符合上述规则配置的实际告警消息：

```
━━━━━━━━━━━━━━━━━━
🚨 触发中 | 【api】服务 172.16.11.236:9100 已宕机
━━━━━━━━━━━━━━━━━━

📋 告警信息
├─ 🔴 级别：🔴 严重
├─ 📛 类型：ServiceDown
├─ 🖥 主机：172.16.11.236:9100
├─ 📝 详情：172.16.11.236:9100 指标采集失败，服务可能已停止
└─ ⏰ 时间：2026-05-29 20:00:00
```

## 六、字段与消息对应关系

| Prometheus Labels/Annotations | 消息显示 |
|------------------------------|---------|
| annotations.summary | 标题 |
| labels.severity | 级别（带 emoji） |
| labels.service | 服务 |
| alertname (自动) | 类型 |
| labels.instance | 主机 |
| annotations.description | 详情 |
| StartsAt (自动) | 触发时间 |

## 七、注意事项

1. **instance 必须包含端口**，用于标识具体主机
2. **service 名称建议与业务相关**，便于快速定位问题
3. **team 用于通知分流**，可按团队发送到不同群
4. **for 表示持续时间**，避免瞬时抖动产生告警