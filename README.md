# Alertmanager Webhook 多平台推送服务

支持将 Prometheus Alertmanager 的告警消息转发到**企业微信、飞书、钉钉**三大平台。

## 功能特性

- 支持多平台：企业微信、飞书、钉钉
- 统一的告警消息格式
- 支持按告警级别显示不同颜色和 emoji
- 自动将 UTC 时间转换为 CST 时区
- 消息超长自动分片

## 快速开始

### 启动服务

```bash
./bin/wechat-webhook \
  -wecom.key=your-wechat-key \
  -dingtalk.key=your-dingtalk-key \
  -feishu.key=your-feishu-key \
  -addr :8999
```

### Docker 部署

```bash
docker run -d \
  --name alertmanager-webhook \
  --network host \
  registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1 \
  /usr/bin/wechat-webhook \
  -wecom.key=your-wechat-key \
  -feishu.key=your-feishu-key \
  -addr :8999
```

## API 接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/webhook/wecom` | POST | 企业微信 |
| `/webhook/dingtalk` | POST | 钉钉 |
| `/webhook/feishu` | POST | 飞书 |
| `/webhook` | POST | 企业微信（兼容旧接口） |
| `/health` | GET | 健康检查 |

## Alertmanager 配置

根据不同平台配置多个 receiver：

```yaml
receivers:
  # 企业微信
  - name: 'wecom'
    webhook_configs:
      - url: 'http://服务地址:8999/webhook/wecom'
        send_resolved: true

  # 钉钉
  - name: 'dingtalk'
    webhook_configs:
      - url: 'http://服务地址:8999/webhook/dingtalk'
        send_resolved: true

  # 飞书
  - name: 'feishu'
    webhook_configs:
      - url: 'http://服务地址:8999/webhook/feishu'
        send_resolved: true
```

## 告警消息格式

```
━━━━━━━━━━━━━━━━━━
🚨 触发中 | 【服务名】告警简述
━━━━━━━━━━━━━━━━━━

📋 告警信息
├─ 🔴 级别：🔴 严重
├─ 📛 类型：ServiceDown
├─ 🖥 主机：192.168.1.100:9100
├─ 📝 详情：服务实例指标采集失败超过 1 分钟
└─ ⏰ 时间：2026-05-29 20:00:00
```

## Prometheus 告警规则标准

### Labels 定义

| Label | 说明 | 示例 |
|-------|------|------|
| severity | 告警级别 | critical / warning / info |
| service | 服务名称 | api / db / nginx |
| instance | 实例地址 | 192.168.1.100:9100 |
| alertname | 告警规则名 | ServiceDown |

### Annotations 定义

| Annotation | 说明 |
|------------|------|
| summary | 告警简述 |
| description | 告警详情 |

### 告警规则示例

```yaml
groups:
- name: example-alerts
  rules:
  - alert: ServiceDown
    expr: up == 0
    for: 1m
    labels:
      severity: critical
      service: your-service-name
    annotations:
      summary: "服务 {{ $labels.instance }} 已宕机"
      description: "{{ $labels.instance }} 指标采集失败超过 1 分钟"
```

## 测试接口

```bash
# 测试企业微信
curl -X POST http://localhost:8999/webhook/wecom \
  -H "Content-Type: application/json" \
  -d '{
    "version":"4","groupKey":"{}:test","status":"firing","receiver":"webhook",
    "groupLabels":{},"commonLabels":{},"commonAnnotations":{"summary":"测试告警：服务实例宕机"},
    "alerts":[{
      "status":"firing",
      "labels":{"alertname":"ServiceDown","severity":"critical","service":"api","instance":"192.168.1.100:9100"},
      "annotations":{"description":"服务实例 192.168.1.100:9100 指标采集失败超过 1 分钟"},
      "startsAt":"2026-05-29T12:00:00Z","endsAt":"0001-01-01T00:00:00Z"
    }]
  }'

# 测试飞书
curl -X POST http://localhost:8999/webhook/feishu \
  -H "Content-Type: application/json" \
  -d '{same payload as above}'

# 测试钉钉
curl -X POST http://localhost:8999/webhook/dingtalk \
  -H "Content-Type: application/json" \
  -d '{same payload as above}'
```

## 构建镜像

```bash
docker build --platform linux/amd64 -t registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1 .
```

## 项目结构

```
├── main.go              # 入口
├── config/             # 配置管理
├── model/              # 消息结构体
├── notifier/          # 各平台通知器
├── transformer/       # 消息转换
└── router/            # 路由处理
```