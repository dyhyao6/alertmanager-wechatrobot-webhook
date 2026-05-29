# Alertmanager 企业微信机器人 Webhook 部署指南

## 一、简介

该服务用于将 Prometheus Alertmanager 的告警转发到企业微信机器人，以 Markdown 格式发送告警通知。

## 二、镜像信息

```
registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1
```

## 三、Docker 部署

```bash
# 拉取镜像（如已推送）
docker pull registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1

# 运行容器
docker run -d \
  --name alertmanager-wechatbot-webhook \
  -p 8999:8999 \
  registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1 \
  /data/wechat-webhook/start.sh "你的RobotKey" ":8999" "grafanaUrl" "alertDomain"
```

## 四、配置参数说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| RobotKey | 企业微信机器人 Webhook Key | - |
| addr | 服务监听地址 | :8999 |
| grafanaUrl | Grafana 面板链接 | grafana.vnnox.com/d/PwMJtdvnr/k8s-chu-neng-cnanduat |
| alertDomain | Alertmanager 告警规则域名 | emscn-prometheus.ampaura.tech |

## 五、Alertmanager 配置

在 Alertmanager 的 `receivers` 中添加 webhook 配置：

```yaml
receivers:
  - name: 'wechatbot'
    webhook_configs:
      - url: 'http://目标机器IP:8999/webhook'
        send_resolved: true
```

### 5.1 多机器人配置

支持三种方式指定机器人：

1. **全局默认**：通过 `-RobotKey` 参数指定
2. **请求级**：通过 URL 参数 `?key=xxx`
3. **告警级**：在 Prometheus 告警规则中添加注解 `wechatRobot: xxx`

```yaml
annotations:
  wechatRobot: your-robot-key-here
```

## 六、日志查看

```bash
docker logs -f alertmanager-wechatbot-webhook
```

## 七、健康检查

```bash
curl http://localhost:8999/webhook
```

## 八、消息格式

告警消息将以 Markdown 格式发送，包含：
- 告警名称和状态（ firing / resolved ）
- 严重级别（ critical / warning / info ）
- 告警描述
- 触发时间（自动转换为 CST 时区）
- Grafana 面板和告警规则链接

## 九、本地构建

```bash
# 构建镜像
docker build -t registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1 .

# 推送镜像（需先登录阿里云容器镜像服务）
docker push registry.cn-hangzhou.aliyuncs.com/novacloud/wechat-webhook-new:ems1.1
```