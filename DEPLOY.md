# Alertmanager Webhook 部署指南

## 一、简介

该服务用于将 Prometheus Alertmanager 的告警转发到企业微信、飞书、钉钉机器人，支持多平台，以 Markdown 格式发送告警通知。

## 二、镜像信息

```
registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1
```

## 三、Docker Compose 部署

```bash
docker compose up -d
```

## 四、配置参数说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| wecom.key | 企业微信机器人 Webhook Key | - |
| feishu.key | 飞书机器人 Webhook Key | - |
| dingtalk.key | 钉钉机器人 Webhook Key | - |
| addr | 服务监听地址 | :8999 |
| log.dir | 日志目录 | /var/log/wechat-webhook |
| log.file | 日志文件名 | webhook.log |

## 五、告警平台配置

### 5.1 微信

```yaml
receivers:
  - name: 'wechatbot'
    webhook_configs:
      - url: 'http://目标机器IP:8999/webhook/wecom'
        send_resolved: true
```

### 5.2 飞书

```yaml
receivers:
  - name: 'feishubot'
    webhook_configs:
      - url: 'http://目标机器IP:8999/webhook/feishu'
        send_resolved: true
```

### 5.3 钉钉

```yaml
receivers:
  - name: 'dingtalkbot'
    webhook_configs:
      - url: 'http://目标机器IP:8999/webhook/dingtalk'
        send_resolved: true
```

## 六、日志查看

```bash
# Docker Compose 日志
docker compose logs -f

# 查看日志文件
cat logs/webhook.log
```

## 七、健康检查

```bash
curl http://localhost:8999/health
```

## 八、消息格式

告警消息将以 Markdown 格式发送，包含：
- 告警名称和状态（触发中 / 已恢复）
- 严重级别（🔴 严重 / 🟡 警告 / 🔵 通知）
- 告警描述
- 触发时间（自动转换为 CST 时区）

## 九、本地构建

```bash
# 构建镜像（支持多平台）
docker build --platform linux/amd64 -t registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1 .

# 推送镜像（需先登录阿里云容器镜像服务）
docker push registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1
```